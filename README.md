# Self Got

A Discord selfbot with a graphical user interface built in Go.

![help](https://files.catbox.moe/unuqw0.png)
## Features

- Graphical user interface using Fyne toolkit
- CLI mode available with `--no-ui` flag
- Configuration via GUI or config.txt file
- Any user can interact with the selfbot


### Available Commands

- **info** - Displays the current memory usage
- **bounce** - Creates a bouncing GIF animation from an image
- **remind** - Sets a reminder after specified duration
- **ocr** - Performs optical character recognition on images
- **delete** - Delete messages in the given channel
- **icon** - Extracts the avatar of a mentioned\replied user

## Installation

1. Install the required dependencies:
   - Go 1.17+
   - Tesseract OCR (for OCR command)
   - libvips (for image processing)

2. Clone the repository:
```sh
git clone https://github.com/awangelo/self-got.git
cd self-got
```

3. Build and run:
```sh
go build
./self-got
```

## Configuration

On first run, you'll be prompted to enter:
- Discord token
- Command prefix (default: \\)

These settings are saved to `config.txt`

## Usage

### The commands are used directly on discord
   - `\info` - Show memory stats
   - `\bounce [image/url]` - Create bouncing GIF
   - `\remind [duration] [message]` - Set a reminder
   - `\ocr [image]` - Extract text from image
   - `\delete [n/all]` - Delete `n` or `all` messages in a channel
   - `\icon @user` - Extract the user avatar

### Usage Notes for Image Commands

Commands that work with images support two ways to provide the image:

1. Direct attachment with the command
2. Replying to a message that contains:
   - An image attachment
   - A valid image URL

For example:
- `\bounce` while attaching an image
- Reply to a message containing an image with `\bounce`
- Reply to a message containing an image URL with `\bounce`
