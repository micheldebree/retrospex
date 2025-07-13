# Retrospex Specification

## Project Overview
Retrospex is a command-line application written in Go that converts any image into an image adhering to the limitations of Commodore 64 graphic modes. The primary goal is to apply these limitations so that the resulting image can be converted and displayed on a real Commodore 64 computer.

## Key Features
- **Graphic Modes**: Supports different modes like bitmap, character mode, and sprites.
- **Multicolor and Hires Variants**: All modes support multicolor and high-resolution variants.
- **Palette Options**: Users can choose from different palettes based on the Commodore 64 standard palette with varying shades of colors.
- **Input/Output Formats**: Supports popular image formats for input and outputs a PNG image.

## Scope
- **Pre-processing**: Pre-processing tasks like scaling or cropping are not in scope. These should be done using other tools.
- **Binary Creation**: Creating Commodore 64 binary formats is not in scope. Tools like `png2prg` can be used for this purpose after Retrospex processing.

## Technical Details
- **Language**: Go (Golang)
- **Input Formats**: Popular image formats (e.g., PNG, JPEG)
- **Output Format**: PNG

## User Workflow
1. **Pre-processing**: Users pre-process the input image to ensure it fits within Commodore 64 screen dimensions if necessary.
2. **Conversion**: Use Retrospex to convert the image to adhere to Commodore 64 graphic mode limitations.
3. **Post-processing**: Convert the resulting PNG image into a Commodore 64 binary format using tools like `png2prg`.

## Example Usage
