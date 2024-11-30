# spoticli

Development: Active

## Description

A program to stream music from the command line. This program can stream or download mp3 files from aws s3 to a command-line user. Downloading can be done directly (presigned url) or via the backend serving the content. Streaming  cannot be done using presigned urls as some backend processing is required.

Streaming is made possible due to the Decoder service. It decodes whats needed to comput the sizes of frames and tags and metadata.

See the README in [spoticli-backend](https://github.com/matttm/spoticli/tree/main/spoticli-backend#spoticli-backend) for more on the algorithm and backend architecture.

## Getting Started

The main parts are `spoticli-cli` and `spoticli-backend`, and a README.md describing their setup will be in each of these subprojects.
## Authors

-   Matt Maloney : matttm

## Contribute

If you want to contribute, just send me a message.
