# W1 IMAGE CLI tool

Helper tool for handy simple image transformations.
Supports jpg and png.

Build Command:

`govvv build -o ./bin/w1-image`

Dev command (grayscale example):

`go run . grayscale ./tmp/example.jpg`

## Specific commands (from binary)

`w1-image grayscale ./tmp/example.jpg`

`w1-image resize ./tmp/example.jpg 240`
