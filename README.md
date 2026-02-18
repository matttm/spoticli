# spoticli

## Description

A CLI tool for streaming and downloading music from Amazon S3. This program streams or downloads MP3 files from AWS S3 to a command-line user. Downloading can be done directly (via presigned URLs) or through the backend serving the content. Streaming cannot be done using presigned URLs as backend processing is required to ensure proper MP3 frame alignment.

Streaming is made possible by the Decoder service, which decodes what's needed to ensure the streamed music chunks are frame-aligned for seamless playback.

See the [spoticli-backend/README.md](./spoticli-backend/README.md) for more on the algorithm and backend architecture.

## Backend Features

The spoticli-backend provides a comprehensive audio streaming API with the following capabilities:

- **Intelligent MP3 Streaming**: Frame-based chunking algorithm that decodes MP3 frames at boundaries for seamless playback
- **Flexible Access Methods**: Audio files accessible via presigned S3 URLs, proxy downloads, or streaming with range request support
- **File Management**: RESTful API for uploading/downloading audio files and querying file metadata by content type
- **Storage Integration**: S3-compatible storage with LocalStack support for local development
- **Performance Optimization**: Caching layer for decoded frame boundaries to improve streaming performance
- **Database-Backed**: MySQL database for file metadata with automatic schema initialization
- **Docker Deployment**: Complete Docker Compose setup with app, database, and LocalStack services
- **Testing Suite**: Integration tests with helper scripts and test asset generation
- **API Documentation**: Full OpenAPI 3.0.3 specification with Postman collections

## Prerequisites

- **Docker and Docker Compose** - For running the backend services locally
- **Go 1.21+** - For building and running the CLI
- **AWS Account** (optional) - Or use LocalStack for local S3 emulation

## Getting Started

The project consists of two main components: `spoticli-cli` and `spoticli-backend`. Each includes a README.md with detailed setup instructions.

### Quick Start

1. **Start the backend services:**
   ```bash
   cd spoticli-backend
   docker-compose up -d
   ```

2. **Build the CLI:**
   ```bash
   cd spoticli-cli
   go build
   ```

3. **Upload and play music:**
   ```bash
   # Upload songs from a directory
   ./spoticli-cli song upload assets/your-music-folder
   
   # Play a song
   ./spoticli-cli song play
   ```

For detailed configuration and advanced usage, see the component READMEs:
- [spoticli-backend/README.md](./spoticli-backend/README.md) - Backend setup and API documentation
- [spoticli-cli/README.md](./spoticli-cli/README.md) - CLI commands and usage

## Example Usage

After starting the backend and uploading songs, you can stream or download music using the CLI. The `play` command presents an interactive prompt showing all songs stored in the database:

```
❯ ./spoticli-cli song play
Use the arrow keys to navigate: ↓ ↑ → ←
? Select Song:
    blinded_in_chains.mp3
    the_wicked_end.mp3
  ▸ bat_country.mp3
    sidewinder.mp3
↓   blinded_in_chains.mp3

```

Use the arrow keys to navigate and press Enter to select a song for streaming.

## Authors

-   Matt Maloney : matttm

## Contribute

If you want to contribute, just send me a message.



