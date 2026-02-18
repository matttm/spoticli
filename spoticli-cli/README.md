# spoticli-cli

A CLI tool for managing and streaming music - upload songs to S3, play music interactively, and download your music library.

## Prerequisites

- **Go 1.21+** - For building the CLI
- **spoticli-backend running** - The backend must be running and accessible (see [spoticli-backend/README.md](../spoticli-backend/README.md))

## Getting started

### Building

To build, just execute:
```
go build
```
Then run with the following:
```
./spoticli-cli
```
This will produce the following output:
```
❯ ./spoticli-cli
NAME:
   spoticli-cli - A CLI tool for managing and streaming music

USAGE:
   spoticli-cli [global options] command [command options]

COMMANDS:
   song     song <action>
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help

```
With more under the song command:
```
❯ ./spoticli-cli song
NAME:
   spoticli-cli song - song <action>

USAGE:
   spoticli-cli song command [command options]

COMMANDS:
   upload    upload <path>
   play      play
   download  download <song-id>
   ls        ls
   help, h   Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help

```

## Usage

### Upload Songs

The `upload` command uploads a directory of MP3 files to S3 and creates database entries for each song. Use relative paths from your home directory:

```bash
# To upload ~/assets/a7x, run:
./spoticli-cli song upload assets/a7x
```

### Play Songs

The `play` command presents an interactive menu to select and stream a song:

```bash
./spoticli-cli song play
```

### List Songs

List all songs in the database:

```bash
./spoticli-cli song ls
```

### Download Songs

Download a specific song by ID:

```bash
./spoticli-cli song download <song-id>
```

