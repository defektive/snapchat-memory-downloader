FROM alpine:latest

ADD dist/snapchat-memory-downloader_linux_amd64_v1/snapchat-memory-downloader /bin/
WORKDIR /workspace
ENTRYPOINT ["/bin/snapchat-memory-downloader"]