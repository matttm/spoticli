# spoticli-backend

This backend provides an API for downloading and streaming music.

For streaming, a naive approach is to segment an MP3 file randomly. This won't work though because mp3 decoders normally require a mp3 header for decoding. And even adding a fake header to randomly segmented partitions will not work because the frame header contains information that may be exclusive to a given frame.

The approach I ended up devising, was to partition each frame at the frame boundaries. Then when data is sent, a cluster of frames are merged into a byte array sent.

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

Also a feature of this backend is when doing a range request, only the start position is respected.

## Getting started

First thing is environment variables. Source the following to your environment before running the backend.
```bash
# variables determining processing
export STREAM_SEGMENT_SIZE=1000000
export FRAME_CLUSTER_SIZE=30

# database
export DB_HOST="localhost"
export DB_PORT=""3306
export DB_USERNAME="ADMIN"
export DB_PASSWORD="ADMIN"

# aws variables
export AWS_ACCESS_KEY_ID="key"
export AWS_SECRET_ACCESS_KEY="secret"
export AWS_REGION=us-east-1


```

Once these are in your enironment, you can run the backend in the terminal as follows
```
$ go build
$ ./spoticli-backend
```
## Architecture
