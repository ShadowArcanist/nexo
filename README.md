# Nexo
A blazing-fast, lightweight single static binary to convert images from `JPG`, `PNG` formats to `Webp`

## Features
- **Lightning Fast** - Converts 30 `.jpg` (Dimensions: 5120 × 2880) images to `.webp` within 10s
- **Ultra Lightweight** - Single ~5MB binary, no external dependencies
- **Zero Config** - Just run and forget
- **Multiple Formats** - Supports JPG, JPEG, PNG
- **Auto-Cleanup** - Automatically deletes original images after conversion
- **Cross-Platform** - Works on macOS, Linux, and Windows
- 
More aren't planned but feel free to add them yourself.

## How It Works
1. **Scan** - Recursively scans the input directory for supported image formats
2. **Queue** - Creates a work queue with all found images
3. **Process** - Uses a worker pool (up to 8 concurrent workers) to convert images in parallel
4. **Encode** - Uses the high-performance `gen2brain/webp` library for WebP encoding
5. **Cleanup** - Removes original images after successful conversion
6. **Exit** - Automatically exits when all conversions are complete


## Real-World Test Results
**High-Resolution Batch Conversion:**
- **Images:** 30 JPEG files
- **Total Size:** 313MB
- **Per Image:** 10.8MB, 5120 × 2880 pixels (5K resolution)
- **Completion Time:** 9.34 seconds
- **Result:** 30 high-resolution images converted to WebP in under 10s


## Installation

### Prerequisites

- [Go 1.26+](https://go.dev/dl/) (for building from source)
- [Task](https://taskfile.dev/installation/) (optional, for using Taskfile)

### Quick Install

```bash
# Clone the repository
git clone <repository-url>
cd imgconverter

# Build the binary
task build

# Or with go directly
go build -o build/imgconverter -ldflags="-s -w" .
```

## Usage

### Basic Usage

1. Create a `converter/` folder in the same directory as the binary
2. Put your images inside the folder
3. Run the binary:

```bash
./build/imgconverter
```

4. Done! All images are converted to `.webp` in the same directory and originals are removed

### Custom Directory

You can specify a different input directory:

```bash
./build/imgconverter /path/to/your/images
```

### Example Output
**Success:**
```
Found 42 image(s) to convert
Conversion completed!
Successful: 42
Time taken: 850ms
```

**With Failures:**
```
Found 50 image(s) to convert
Conversion completed!
Successful: 45
Failure: 5
Failed Images: corrupt1.png, corrupt2.jpg, missing.png, bad.png, error.jpg
Time taken: 1.23s
```

## Building

### Build for Current Platform

```bash
task build
```

### Build for All Platforms

```bash
task build-all
```

This creates binaries for:
- macOS (Intel & Apple Silicon)
- Linux (AMD64 & ARM64)
- Windows (AMD64)


## Troubleshooting
### "converter directory not found"

The binary looks for a `converter/` folder by default. Create it:
```bash
mkdir converter
```

Or specify a custom path:
```bash
./build/imgconverter /path/to/images
```

### "No supported images found"
Ensure your images have supported extensions:
- `.jpg` or `.jpeg`
- `.png`


### Permission Denied
Make the binary executable:
```bash
chmod +x build/imgconverter
```


## Notes
1. This project was entirely created using AI, but the application has been thoroughly tested (only on macOS Sequoia).
2. This project was built primarily for my personal use, so I will not be merging pull requests or adding new features unless I need them myself. If you want to make changes or add features, feel free to fork this repository. It’s open-sourced so others can learn from it, use it as a base for their own projects, or even run the application as-is.