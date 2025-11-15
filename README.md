# Snapchat Memory Downloader

- [x] Download snapchat memories
- [x] Add date EXIF information
- [ ] Add geo tagging EXIF data

## Why?

Snapchat's memory export did not have date information on images. Exporting lots of memories was time-consuming and error-prone.

## Usage

```bash
snapchat-memory-downloader -f json/memories_history.json
```

This will download everything to the `./downloads` directory.


## Install

### Downloading

Download a binary from the [releases](https://github.com/defektive/snapchat-memory-downloader/releases) section.

### Building

```bash
go install github.com/defektive/snapchat-memory-downloader@latest
```

