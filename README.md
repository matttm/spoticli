# spoticli

Development: Complete

## Description

A program to stream music from the command line. This program can stream or download mp3 files from aws s3 to a command-line user. Downloading can be done directly (presigned url) or via the backend serving the content. Streaming  cannot be done using presigned urls as some backend processing is required.

Streaming is made possible due to the Decoder service. It decodes whats needed to comput the sizes of frames and tags and metadata.

See the README in [spoticli-backend](https://github.com/matttm/spoticli/tree/main/spoticli-backend#spoticli-backend) for more on the algorithm and backend architecture.

## Getting Started

The main parts are `spoticli-cli` and `spoticli-backend`, and a README.md describing their setup will be in each of these subprojects.

## Example

Running the stream or download command, will generate a prompt of songs to choose from. The songs are what is being stored in the database which is running in a docker container.
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

## Authors

-   Matt Maloney : matttm

## Contribute

If you want to contribute, just send me a message.
