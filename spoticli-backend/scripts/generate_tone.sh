#!/usr/bin/env bash
set -euo pipefail
# Generates simple test audio files using ffmpeg.
# Usage: ./generate_tone.sh [out_dir]
# Requires ffmpeg installed on the PATH.

# Check dependencies
deps=(ffmpeg)
missing=()
for cmd in "${deps[@]}"; do
	if ! command -v "$cmd" >/dev/null 2>&1; then
		missing+=("$cmd")
	fi
done

if [ ${#missing[@]} -ne 0 ]; then
	echo "Error: missing required dependencies: ${missing[*]}" >&2
	echo "Please install the missing tools and ensure they're on your PATH." >&2
	echo "On macOS you can often run: brew install ffmpeg" >&2
	exit 1
fi

OUT_DIR="${1:-./assets}"
mkdir -p "$OUT_DIR"

# 5s 440Hz sine wave WAV
ffmpeg -y -f lavfi -i "sine=frequency=440:duration=5" -c:a pcm_s16le "$OUT_DIR/test_tone.wav"

# Same tone encoded to MP3 (VBR quality)
ffmpeg -y -f lavfi -i "sine=frequency=440:duration=5" -c:a libmp3lame -q:a 2 "$OUT_DIR/test_tone.mp3"

echo "Generated: $OUT_DIR/test_tone.wav and $OUT_DIR/test_tone.mp3"
