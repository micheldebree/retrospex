# Retrospex Development Guidelines

## Commands

- Build: `make` or `go build -v -trimpath -ldflags="-w -s -X main.Version=$(VERSION).$(REV)"`
- Install: `make install` (installs to ~/bin/retrospex)
- Test: `make test` or `go test ./...`
- Clean: `make clean` (removes build artifacts and generated files)
- Examples: `make examples` (builds example Commodore 64 files)
- Cross-compilation: Use specific targets like `make retrospex_macos_arm64`

## Testing

- Uses standard Go testing package (not testify)
- Test files are named \*\_test.go
- Run individual test package: `go test ./internal/c64io`
- Run single test function: `go test -run TestGetBitmapData ./internal/c64io`

## Project Structure

- main.go: Application entry point
- internal/: Library packages (c64io, conversion, dithering, imageio,
  indexedimage, pixels, prepost)
- c64/: Commodore 64 specific files and makefile includes
- examples/: Example usage files

## Important Notes

- Binary output format varies by graphics mode (see README.md for size details)
- Default graphics mode is "koala"
- Default output format is "png" (use -f bin for Commodore 64 binary data)
- The tool requires exactly one input image argument
- Color palettes include: colodore, pepto, vice, vice_old_lum,
  vice_371_internal, pantaloon, archmage, electric, ste, perplex_1, perplex_2,
  lenyn

