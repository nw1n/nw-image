# NW-IMAGE CLI tool

Helper tool for handy simple image transformations.
Supports jpg and png.

Build Command:

`go build -o ./bin/nw-image`

Dev command (grayscale example):

`go run . grayscale ./tmp/example.jpg`

## Specific commands (from binary)

`nw-image grayscale ./tmp/example.jpg`

`nw-image resize ./tmp/example.jpg 240`
