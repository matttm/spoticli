# spoticli

Development: Active

## Description

A program to stream music from the command line. This program can stream or download mp3 files from aws s3 to a command-line user. Downloading can be done directly (presigned url) or via the backend serving the content. Streaming  cannot be done using presigned urls as some backend processing is required.

This backend processing includes removing the ID3v2 header from the beginning of the mp3 file, and then partitioning it into clusters of mp3 frames, to ensure every payload sent to the frontend contains the mp3 header.

See the README in Spoticli-backend for more on the algorithm.

## Getting Started

The main parts are `spoticli-cli` and `spoticli-backend`, and a README.md describing their setup will be in each of these subprojects.
## Authors

-   Matt Maloney : matttm

## Contribute

If you want to contribute, just send me a message.
