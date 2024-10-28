# spoticli

Development: Active

## Description

A program to stream music from the command line. This program can stream or download mp3 files from aws s3 to a command-line user. Downloading can be done directly (presigned url) or via the backend serving the content. Streaming  cannot be done using presigned urls as some backend processing is required.

This backend processing includes removing the ID3v2 header from the beginning of the mp3 file, and then partitioning it into clusters of mp3 frames, to ensure every payload sent to the frontend contains the mp3 header.

The anatomy of an mp3 is as shown below,
```
+---------------------+
|   MP3 File Header   |  --> Metadata (e.g., ID3v2 tags)
+---------------------+
|    Audio Frame 1    |  --> Contains header + data
+---------------------+
|    Audio Frame 2    |
+---------------------+
|    Audio Frame 3    |
+---------------------+
|    Audio Frame N    |  --> Last audio frame
+---------------------+
|  Optional Metadata  |  --> Footer (e.g., ID3v1 tags)
+---------------------+
```
In my algorithm, I strip the initial ID3v2 tag, then break the file apart by frames, as shown below,
```
+-------------------+   +-------------------+   +-------------------+   +-------------------+
| Frame 1 Header    |   | Frame 2 Header    |   | Frame 3 Header    |   | Frame N Header    |
+-------------------+   +-------------------+   +-------------------+   +-------------------+
| Frame 1 Data      |   | Frame 2 Data      |   | Frame 3 Data      |   | Frame N Data      |
+-------------------+   +-------------------+   +-------------------+   +-------------------+
    |                     |                     |                      |
(1152 samples)      (1152 samples)      (1152 samples)      (1152 samples)  
```
The frame slices are then grouped together, such that there is x frames per cluster.

## Getting Started

The main parts are `spoticli-cli` and `spoticli-backend`, and a README.md describing their setup will be in each of these subprojects.
## Authors

-   Matt Maloney : matttm

## Contribute

If you want to contribute, just send me a message.
