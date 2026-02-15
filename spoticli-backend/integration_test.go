package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestIntegration_UploadDownloadStream uploads a small blob via the presigned
// flow, then verifies download and streaming proxy endpoints. The test
// assumes the integration services are already running (the `scripts/run-
// integration.sh` helper manages docker-compose lifecycle).
func TestIntegration_UploadDownloadStream(t *testing.T) {
	fmt.Println("=== running TestIntegration_UploadDownloadStream ===")
	if testing.Short() {
		t.Skip("skipping integration test in -short mode")
	}
	// Require explicit opt-in via environment variable so integration tests
	// don't run during ordinary `go test` runs.
	if os.Getenv("INTEGRATION") != "1" {
		t.Skip("skipping integration test; set INTEGRATION=1 or use ./scripts/run-integration.sh to run")
	}
	// repo root (test now lives at repo root)
	repoRoot := "."

	client := &http.Client{Timeout: 5 * time.Second}
	healthUrl := "http://localhost:4200/"
	ready := false
	deadline := time.Now().Add(3 * time.Minute)
	for time.Now().Before(deadline) {
		resp, err := client.Get(healthUrl)
		if err == nil {
			t.Logf("health check status: %d", resp.StatusCode)
			if resp.StatusCode == http.StatusOK {
				ready = true
				resp.Body.Close()
				break
			}
			resp.Body.Close()
		} else {
			t.Logf("health check error: %v", err)
		}
		time.Sleep(2 * time.Second)
	}
	if !ready {
		t.Fatal("backend did not become ready in time")
	}

	track := "integration/test-integration.mp3"
	audioPath := filepath.Join(repoRoot, "assets", "test_tone.mp3")
	b64Path := filepath.Join(repoRoot, "assets", "test_tone.mp3.b64")
	var fileBytes []byte
	if b, err := os.ReadFile(b64Path); err == nil {
		if decoded, err := base64.StdEncoding.DecodeString(strings.TrimSpace(string(b))); err == nil {
			fileBytes = decoded
		} else {
			t.Logf("failed to decode base64 fixture: %v", err)
		}
	}
	if fileBytes == nil {
		if b, err := os.ReadFile(audioPath); err == nil {
			fileBytes = b
		} else {
			fileBytes = []byte("THIS IS A TEST MP3 BLOB")
		}
	}

	payload := map[string]interface{}{"Key_name": track, "File_size": len(fileBytes)}
	pb, _ := json.Marshal(payload)
	resp, err := client.Post("http://localhost:4200/audio", "application/json", bytes.NewReader(pb))
	if err != nil {
		t.Fatalf("failed to call upload endpoint: %v", err)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	t.Logf("upload endpoint status: %d body: %s", resp.StatusCode, string(body))
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("upload endpoint returned %d: %s", resp.StatusCode, string(body))
	}
	presigned := string(body)
	if presigned == "" {
		t.Fatal("empty presigned url returned")
	}
	t.Logf("presigned URL: %s", presigned)
	if strings.Contains(presigned, "localstack") {
		if u, err := url.Parse(presigned); err == nil {
			host := u.Host
			if !strings.Contains(host, ":") {
				host = "localhost:4566"
			} else {
				parts := strings.Split(host, ":")
				host = "localhost:" + parts[len(parts)-1]
			}
			u.Host = host
			presigned = u.String()
			t.Logf("rewritten presigned URL for host access: %s", presigned)
		}
	}

	putReq, _ := http.NewRequest("PUT", presigned, bytes.NewReader(fileBytes))
	putReq.Header.Set("Content-Type", "application/octet-stream")
	putResp, err := client.Do(putReq)
	if err != nil {
		t.Fatalf("failed to PUT to presigned url: %v", err)
	}
	putBody, _ := io.ReadAll(putResp.Body)
	putResp.Body.Close()
	t.Logf("presigned PUT status: %d body: %s", putResp.StatusCode, string(putBody))
	if putResp.StatusCode < 200 || putResp.StatusCode >= 300 {
		t.Fatalf("presigned PUT returned %d: %s", putResp.StatusCode, string(putBody))
	}

	filesResp, err := client.Get("http://localhost:4200/files/1")
	if err != nil {
		t.Fatalf("failed to query files list: %v", err)
	}
	filesBody, _ := io.ReadAll(filesResp.Body)
	filesResp.Body.Close()
	t.Logf("files list body: %s", string(filesBody))
	var files []struct {
		Id       int    `json:"Id"`
		Key_name string `json:"Key_name"`
	}
	if err := json.Unmarshal(filesBody, &files); err != nil {
		t.Fatalf("failed to parse files list: %v -- body: %s", err, string(filesBody))
	}
	var id int
	for _, f := range files {
		if f.Key_name == track {
			id = f.Id
			break
		}
	}
	if id == 0 {
		t.Fatalf("uploaded file not found in files list: %s", string(filesBody))
	}

	dlUrl := fmt.Sprintf("http://localhost:4200/audio/proxy/%d", id)
	dlResp, err := client.Get(dlUrl)
	if err != nil {
		t.Fatalf("failed download proxy: %v", err)
	}
	dlBytes, _ := io.ReadAll(dlResp.Body)
	dlResp.Body.Close()
	if dlResp.StatusCode != http.StatusOK {
		t.Fatalf("download proxy returned %d", dlResp.StatusCode)
	}
	if len(dlBytes) != len(fileBytes) {
		t.Fatalf("downloaded length mismatch: expected %d got %d", len(fileBytes), len(dlBytes))
	}

	streamUrl := fmt.Sprintf("http://localhost:4200/audio/proxy/stream/%d", id)
	streamReq, _ := http.NewRequest("GET", streamUrl, nil)
	streamReq.Header.Set("Range", "bytes=0-9")
	streamResp, err := client.Do(streamReq)
	if err != nil {
		t.Fatalf("failed stream proxy: %v", err)
	}
	streamData, _ := io.ReadAll(streamResp.Body)
	streamResp.Body.Close()
	if streamResp.StatusCode != http.StatusPartialContent {
		t.Fatalf("expected 206 Partial Content, got %d", streamResp.StatusCode)
	}
	if len(streamData) != 10 {
		t.Fatalf("streamed bytes length expected 10 got %d", len(streamData))
	}
}
