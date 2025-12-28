# MockGen Usage Guide

This document explains how to use `mockgen` to generate mock implementations for testing in this project.

## Overview

This project uses [gomock](https://github.com/uber-go/mock) for creating mock implementations of interfaces in unit tests. The `mockgen` tool automatically generates these mocks from interface definitions.

## Installation

Install mockgen:

```bash
go install go.uber.org/mock/mockgen@latest
```

Verify installation:

```bash
mockgen -version
```

## Project Structure

Our testing pattern follows these conventions:

```
internal/services/
├── audio.go                    # Service implementation
├── audio_api.go                # Interface definition
├── audio_test.go               # Unit tests
└── mocks/
    ├── audio_mock.go           # Generated mock (DO NOT EDIT)
    ├── config_mock.go
    ├── file-meta-info_mock.go
    └── storage_mock.go
```

**Key Pattern:**
- `*_api.go` files contain interface definitions
- `mocks/*_mock.go` files contain generated mocks
- Mock files should never be manually edited (they're regenerated)

## Generating Mocks

### For a Single Interface

Generate a mock from an interface file:

```bash
mockgen -source=internal/services/<name>_api.go \
        -destination=internal/services/mocks/<name>_mock.go
```

### Examples from This Project

**Audio Service:**
```bash
mockgen -source=internal/services/audio_api.go \
        -destination=internal/services/mocks/audio_mock.go
```

**Config Service:**
```bash
mockgen -source=internal/services/config_api.go \
        -destination=internal/services/mocks/config_mock.go
```

**Storage Service:**
```bash
mockgen -source=internal/services/storage_api.go \
        -destination=internal/services/mocks/storage_mock.go
```

**File Meta Info Service:**
```bash
mockgen -source=internal/services/file-meta-info_api.go \
        -destination=internal/services/mocks/file-meta-info_mock.go
```

### Generate All Mocks

To regenerate all mocks at once:

```bash
# From project root
mockgen -source=internal/services/audio_api.go -destination=internal/services/mocks/audio_mock.go
mockgen -source=internal/services/config_api.go -destination=internal/services/mocks/config_mock.go
mockgen -source=internal/services/file-meta-info_api.go -destination=internal/services/mocks/file-meta-info_mock.go
mockgen -source=internal/services/storage_api.go -destination=internal/services/mocks/storage_mock.go
```

## Using Mocks in Tests

### Basic Example

```go
package services

import (
    "testing"
    
    mock_services "github.com/matttm/spoticli/spoticli-backend/internal/services/mocks"
    "github.com/stretchr/testify/assert"
    "go.uber.org/mock/gomock"
)

func TestExample(t *testing.T) {
    // 1. Create a mock controller
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    // 2. Create a mock instance
    mockAudio := mock_services.NewMockApiAudioService(ctrl)
    
    // 3. Set expectations
    mockAudio.EXPECT().
        GetPresignedUrl(1).
        Return("https://s3.amazonaws.com/...", nil).
        Times(1)
    
    // 4. Use the mock in your test
    url, err := mockAudio.GetPresignedUrl(1)
    
    // 5. Assert results
    assert.NoError(t, err)
    assert.Equal(t, "https://s3.amazonaws.com/...", url)
}
```

### Real Example from This Project

From `storage_test.go`:

```go
func TestStorageService_GetPresignedUrl_Success(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockPresignClient := mock_services.NewMockS3PresignClientApi(ctrl)
    storageService := &StorageService{
        psClient: mockPresignClient,
    }

    expectedURL := "https://s3.amazonaws.com/spoticli-tracks/test-key?presigned=true"
    key := "test-key"

    mockPresignClient.EXPECT().
        PresignGetObject(
            gomock.Any(),
            gomock.Any(),
        ).
        Return(&s3.PresignedRequest{URL: expectedURL}, nil).
        Times(1)

    url, err := storageService.GetPresignedUrl(key)

    assert.NoError(t, err)
    assert.Equal(t, expectedURL, url)
}
```

## Common Mockgen Options

| Option | Description | Example |
|--------|-------------|---------|
| `-source` | Source file containing interfaces | `-source=audio_api.go` |
| `-destination` | Output file for generated mocks | `-destination=mocks/audio_mock.go` |
| `-package` | Package name for generated code | `-package=mock_services` |
| `-mock_names` | Custom mock type names | `-mock_names=Reader=MockReader` |

## Gomock Matchers

Use these matchers in `EXPECT()` calls:

- `gomock.Any()` - Matches any value
- `gomock.Eq(x)` - Matches values equal to x
- `gomock.Nil()` - Matches nil
- `gomock.Not(x)` - Matches values not equal to x
- `gomock.AssignableToTypeOf(x)` - Matches values assignable to x's type

Example:
```go
mockClient.EXPECT().
    GetObject(
        gomock.Any(),                    // Any context
        gomock.Eq(&s3.GetObjectInput{    // Exact match for input
            Bucket: aws.String("my-bucket"),
            Key:    aws.String("my-key"),
        }),
    ).
    Return(mockOutput, nil).
    Times(1)
```

## When to Regenerate Mocks

Regenerate mocks whenever:

1. You add a new method to an interface
2. You change a method signature in an interface
3. You create a new `*_api.go` file
4. You encounter "undefined method" errors in tests

## Creating New Interfaces

When adding a new service with mocks:

1. **Create the interface file** (`*_api.go`):
   ```go
   package services
   
   type MyServiceApi interface {
       DoSomething(id int) (string, error)
   }
   ```

2. **Update your service to use the interface**:
   ```go
   type MyService struct {
       dependency MyServiceApi
   }
   ```

3. **Generate the mock**:
   ```bash
   mockgen -source=internal/services/my_service_api.go \
           -destination=internal/services/mocks/my_service_mock.go
   ```

4. **Create tests** using the generated mock:
   ```go
   mockDep := mock_services.NewMockMyServiceApi(ctrl)
   ```

## Troubleshooting

### "undefined: MockXXX"

**Solution:** Regenerate the mock file. The interface likely changed.

### "expected call" errors

**Solution:** Check that your `EXPECT()` matchers match the actual arguments passed. Use `gomock.Any()` for flexible matching.

### Import cycle errors

**Solution:** Ensure interfaces are in separate `*_api.go` files, not in the same file as implementations.

## Best Practices

1. ✅ **DO** define interfaces in separate `*_api.go` files
2. ✅ **DO** add mock generation commands to project documentation
3. ✅ **DO** use `defer ctrl.Finish()` in every test
4. ✅ **DO** use descriptive test names: `TestService_Method_Condition`
5. ❌ **DON'T** manually edit generated mock files
6. ❌ **DON'T** commit merge conflicts in mock files (regenerate instead)
7. ❌ **DON'T** use mocks for simple functions (use real implementations)

## Additional Resources

- [gomock GitHub](https://github.com/uber-go/mock)
- [gomock Documentation](https://pkg.go.dev/go.uber.org/mock/gomock)
- [Testing in Go](https://go.dev/doc/tutorial/add-a-test)
