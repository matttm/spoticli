# spoticli-backend

This backend provides an API for downloading and streaming music.

For streaming, a naive approach is to segment an MP3 file randomly. This won't work though because mp3 decoders normally require a mp3 header for decoding. And even adding a fake header to randomly segmented partitions will not work because the frame header contains information that may be exclusive to a given frame.

The approach I ended up devising was to partition each frame at the frame boundaries. Then when data is sent, a cluster of frames are merged into a byte array sent.

Also a feature of this backend is when doing a range request, only the start position is respected.

## Getting started

First thing is environment variables. Source the following to your environment before running the backend.
```bash
```
