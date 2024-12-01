# spoticli-cli

## Getting started

### Building

To build, just execute
```
go build
```
Then run with the following
```
./spoticli-cli
```
This will produce the following output
```
❯ ./spoticli-cli
NAME:
   spoticli-cli - A new cli application

USAGE:
   spoticli-cli [global options] command [command options]

COMMANDS:
   song     song <action>
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help

```
With more under the song command
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
There is also the `./spoticli-cli song upload <path>` which will upload a directory to s3 and create db entries for each song in the database. For now the path shold be from the user's directory. So if you want to upload `~/assets/a7x`. Run the following:
```
./spoticli-cli song upload assets/a7x
```

