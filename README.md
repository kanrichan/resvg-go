<div align="center">
  <img src=".github/hua_nobg_512.gif" alt="椛" width = "400">
  <br>

  <h1>resvg-go</h1>
  <img src="https://counter.seku.su/cmoe?name=resvgo&theme=r34" /><br>
  A SVG renderer written in Go & WASM depended on resvg without CGO<br><br>
  
  [![Go Reference](https://pkg.go.dev/badge/github.com/kanrichan/resvg-go.svg)](https://pkg.go.dev/github.com/kanrichan/resvg-go)
  [![Go Report Card](https://goreportcard.com/badge/github.com/kanrichan/resvg-go)](https://goreportcard.com/report/github.com/kanrichan/resvg-go)
</div>

## Table of Contents

- [Overview](#overview)
- [Key Features](#key-features)
- [Installation](#installation)
- [Basic Usage](#basic-usage)
  - [Simple Rendering](#simple-rendering)
  - [Rendering with Options](#rendering-with-options)
- [Advanced Features](#advanced-features)
  - [Transformation Functions](#transformation-functions)
  - [Font Management](#font-management)
  - [Customizing Rendering Options](#customizing-rendering-options)
  - [Image Processing](#image-processing)
- [Examples](#examples)
  - [Converting SVG to PNG](#converting-svg-to-png)
  - [Custom Font Rendering](#custom-font-rendering)
  - [Batch Processing](#batch-processing)
- [Performance Tips](#performance-tips)
- [Troubleshooting](#troubleshooting)
- [Architecture](#architecture)
- [Contributing](#contributing)
- [Important Notes](#important-notes)
- [Acknowledgements](#acknowledgements)
- [License](#license)

## Overview

`resvg-go` is an SVG renderer implemented in pure Go and WebAssembly without CGO dependencies. This library uses the [resvg](https://github.com/RazrFalcon/resvg) Rust library compiled to WebAssembly and invoked through Go's [wazero](https://github.com/tetratelabs/wazero) runtime.

The library provides SVG rendering features in Go applications, offering high-quality rendering with minimal dependencies. By leveraging WebAssembly, `resvg-go` achieves the performance and correctness of the Rust implementation while providing a pure Go API without CGO.

## Key Features

- **No CGO Dependencies**: Implemented in pure Go and WebAssembly
- **Various Transformations**: Support for rotation, scaling, translation, skewing, and more
- **Font Management**: Custom font loading and configuration
- **Image Processing**: PNG encoding/decoding support
- **Rendering Options**: Detailed rendering option customization
- **High Performance**: Leverages the speed of the Rust resvg implementation
- **Minimal Dependencies**: Only requires wazero as a dependency

## Installation

### Installing the Package

```bash
go get github.com/kanrichan/resvg-go
```

### Verifying Installation

After installation, you can verify that everything is working correctly by running the test suite:

```bash
go test github.com/kanrichan/resvg-go
```

## Basic Usage

### Simple Rendering

The simplest way to render an SVG to PNG:

```go
// Initialize a worker (not goroutine-safe!)
worker, _ := resvg.NewDefaultWorker(context.Background())
defer worker.Close()

// Render SVG to PNG
png, _ := worker.Render(svg)
```

This code snippet creates a new worker with default settings, renders the SVG, and returns the PNG data.

### Rendering with Options

For more fine-grained control:

```go
// Initialize a worker
worker, _ := resvg.NewDefaultWorker(context.Background())
defer worker.Close()

// Create and configure a font database
fontdb, _ := worker.NewFontDBDefault()
defer fontdb.Close()
fontdb.LoadFontData(ttf)  // Load custom font data

// Create a pixmap (specify output image size)
pixmap, _ := worker.NewPixmap(512, 512)
defer pixmap.Close()

// Create and render an SVG tree
tree, _ := worker.NewTreeFromData(svg, &resvg.Options{
    Dpi: 96.0,
    FontFamily: "Arial",
    FontSize: 14.0,
})
defer tree.Close()
tree.ConvertText(fontdb)
tree.Render(resvg.TransformIdentity(), pixmap)

// Encode to PNG
png, _ := pixmap.EncodePNG()
```

This example demonstrates a more complete rendering pipeline with custom font configuration and rendering options.

## Advanced Features

### Transformation Functions

You can apply various transformations when rendering SVGs:

```go
// Identity transform (no transformation)
transform := resvg.TransformIdentity()

// Translation transform (move by 100 on x-axis and 50 on y-axis)
transform = resvg.TransformFromTranslate(100.0, 50.0)

// Scaling transform (scale by 2.0 horizontally and 1.5 vertically)
transform = resvg.TransformFromScale(2.0, 1.5)

// Skewing transform
transform = resvg.TransformFromSkew(0.5, 0.3)

// Rotation transform (rotate by 45 degrees)
transform = resvg.TransformFromRotate(45.0)

// Rotation around a specific point (rotate by 45 degrees around point x:100, y:100)
transform = resvg.TransformFromRotateAt(45.0, 100.0, 100.0)

// Custom transformation matrix
transform = resvg.TransformFromRow(sx, ky, kx, sy, tx, ty)
```

These transformations can be combined by applying them in sequence to achieve complex visual effects.

### Font Management

The library provides various font configuration features:

```go
// Create a font database
fontdb, _ := worker.NewFontDBDefault()
defer fontdb.Close()

// Load a font file
fontdb.LoadFontFile("/path/to/font.ttf")

// Load fonts from a directory
fontdb.LoadFontsDir("/path/to/fonts")

// Load font data directly
fontdb.LoadFontData(ttfBytes)

// Set font family preferences
fontdb.SetSerifFamily("Times New Roman")
fontdb.SetSansSerifFamily("Arial")
fontdb.SetCursiveFamily("Comic Sans MS")
fontdb.SetFantasyFamily("Impact")
fontdb.SetMonospaceFamily("Courier New")

// Check the size of the font database
count, _ := fontdb.Len()
```

Properly configuring fonts is essential for maintaining text fidelity in SVG rendering.

### Customizing Rendering Options

You can adjust various rendering options to customize the output:

```go
options := &resvg.Options{
    ResourcesDir: "/path/to/resources",
    Dpi: 300.0,  // High-resolution output
    FontFamily: "Helvetica",
    FontSize: 16.0,
    Languages: []string{"en", "en-US"},
    
    // Set rendering modes
    ShapeRenderingMode: resvg.ShapeRenderingModeGeometricPrecision,
    TextRenderingMode: resvg.TextRenderingModeOptimizeLegibility,
    ImageRenderingMode: resvg.ImageRenderingModeOptimizeQuality,
    
    // Set default sizes
    DefaultSizeWidth: 800.0,
    DefaultSizeHeight: 600.0,
}

tree, _ := worker.NewTreeFromData(svg, options)
```

The available rendering options are:

| Option | Description | Default Value |
|--------|-------------|---------------|
| ResourcesDir | Directory for resolving relative paths | None |
| Dpi | Target DPI for units conversion | 96.0 |
| FontFamily | Default font family | Times New Roman |
| FontSize | Default font size | 12.0 |
| Languages | Languages for systemLanguage resolution | ["en"] |
| ShapeRenderingMode | Default shape rendering method | GeometricPrecision |
| TextRenderingMode | Default text rendering method | OptimizeLegibility |
| ImageRenderingMode | Default image rendering method | OptimizeQuality |
| DefaultSizeWidth | Default viewport width | 100.0 |
| DefaultSizeHeight | Default viewport height | 100.0 |

### Image Processing

The library provides PNG image processing capabilities:

```go
// Create an empty pixmap
pixmap, _ := worker.NewPixmap(width, height)
defer pixmap.Close()

// Create a pixmap from PNG data
pngData := loadPNGFromSomewhere()
pixmap, _ := worker.NewPixmapDecodePNG(pngData)
defer pixmap.Close()

// Encode pixmap to PNG
pngBytes, _ := pixmap.EncodePNG()
```

## Examples

### Converting SVG to PNG

This example shows how to convert an SVG file to a PNG file:

```go
package main

import (
    "context"
    "io/ioutil"
    "log"
    "os"

    "github.com/kanrichan/resvg-go"
)

func main() {
    // Read SVG file
    svgData, err := ioutil.ReadFile("input.svg")
    if err != nil {
        log.Fatalf("Failed to read SVG file: %v", err)
    }

    // Initialize worker
    worker, err := resvg.NewDefaultWorker(context.Background())
    if err != nil {
        log.Fatalf("Failed to create worker: %v", err)
    }
    defer worker.Close()

    // Render SVG to PNG
    pngData, err := worker.Render(svgData)
    if err != nil {
        log.Fatalf("Failed to render SVG: %v", err)
    }

    // Write PNG file
    err = ioutil.WriteFile("output.png", pngData, 0644)
    if err != nil {
        log.Fatalf("Failed to write PNG file: %v", err)
    }

    log.Println("Successfully converted SVG to PNG")
}
```

### Custom Font Rendering

This example demonstrates how to use custom fonts for SVG rendering:

```go
package main

import (
    "context"
    "io/ioutil"
    "log"

    "github.com/kanrichan/resvg-go"
)

func main() {
    // Read SVG file
    svgData, err := ioutil.ReadFile("text_heavy.svg")
    if err != nil {
        log.Fatalf("Failed to read SVG file: %v", err)
    }

    // Read font data
    fontData, err := ioutil.ReadFile("custom_font.ttf")
    if err != nil {
        log.Fatalf("Failed to read font file: %v", err)
    }

    // Initialize worker
    worker, err := resvg.NewDefaultWorker(context.Background())
    if err != nil {
        log.Fatalf("Failed to create worker: %v", err)
    }
    defer worker.Close()

    // Create font database
    fontdb, err := worker.NewFontDBDefault()
    if err != nil {
        log.Fatalf("Failed to create font database: %v", err)
    }
    defer fontdb.Close()

    // Load custom font
    err = fontdb.LoadFontData(fontData)
    if err != nil {
        log.Fatalf("Failed to load font data: %v", err)
    }

    // Set up pixmap
    pixmap, err := worker.NewPixmap(800, 600)
    if err != nil {
        log.Fatalf("Failed to create pixmap: %v", err)
    }
    defer pixmap.Close()

    // Create SVG tree with custom options
    options := &resvg.Options{
        FontFamily: "CustomFont",
        FontSize: 16.0,
    }
    tree, err := worker.NewTreeFromData(svgData, options)
    if err != nil {
        log.Fatalf("Failed to create tree: %v", err)
    }
    defer tree.Close()

    // Apply fonts and render
    tree.ConvertText(fontdb)
    err = tree.Render(resvg.TransformIdentity(), pixmap)
    if err != nil {
        log.Fatalf("Failed to render tree: %v", err)
    }

    // Save result
    pngData, err := pixmap.EncodePNG()
    if err != nil {
        log.Fatalf("Failed to encode PNG: %v", err)
    }
    err = ioutil.WriteFile("output_with_custom_font.png", pngData, 0644)
    if err != nil {
        log.Fatalf("Failed to write PNG file: %v", err)
    }

    log.Println("Successfully rendered SVG with custom font")
}
```

### Batch Processing

This example shows how to process multiple SVG files efficiently:

```go
package main

import (
    "context"
    "io/ioutil"
    "log"
    "os"
    "path/filepath"
    "strings"
    "sync"

    "github.com/kanrichan/resvg-go"
)

func main() {
    // Create a worker pool
    workerCount := 4
    svgChan := make(chan string, 100)
    var wg sync.WaitGroup

    // Start worker goroutines
    for i := 0; i < workerCount; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            
            // Each worker gets its own resvg worker
            worker, err := resvg.NewDefaultWorker(context.Background())
            if err != nil {
                log.Printf("Failed to create worker: %v", err)
                return
            }
            defer worker.Close()
            
            // Process files from the channel
            for svgPath := range svgChan {
                // Read SVG
                svgData, err := ioutil.ReadFile(svgPath)
                if err != nil {
                    log.Printf("Failed to read SVG file %s: %v", svgPath, err)
                    continue
                }
                
                // Render
                pngData, err := worker.Render(svgData)
                if err != nil {
                    log.Printf("Failed to render SVG %s: %v", svgPath, err)
                    continue
                }
                
                // Write PNG
                pngPath := strings.TrimSuffix(svgPath, filepath.Ext(svgPath)) + ".png"
                err = ioutil.WriteFile(pngPath, pngData, 0644)
                if err != nil {
                    log.Printf("Failed to write PNG file %s: %v", pngPath, err)
                    continue
                }
                
                log.Printf("Converted %s to %s", svgPath, pngPath)
            }
        }()
    }
    
    // Walk the directory and find SVG files
    inputDir := "svg_files"
    err := filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() && strings.ToLower(filepath.Ext(path)) == ".svg" {
            svgChan <- path
        }
        return nil
    })
    
    if err != nil {
        log.Fatalf("Failed to walk directory: %v", err)
    }
    
    // Close channel and wait for workers to finish
    close(svgChan)
    wg.Wait()
    
    log.Println("Batch processing complete")
}
```

## Performance Tips

To get the best performance from resvg-go:

1. **Reuse Workers**: Create workers once and reuse them rather than creating new ones for each rendering operation.
2. **Worker per Goroutine**: Since workers are not goroutine-safe, use one worker per goroutine.
3. **Close Resources**: Always close resources (Workers, FontDB, Pixmap, Tree) to prevent memory leaks.
4. **Batch Processing**: When processing multiple SVGs, use a worker pool pattern as shown in the batch processing example.
5. **Image Size Considerations**: Set the pixmap to the exact size needed rather than creating oversized pixmaps.

## Troubleshooting

### Common Issues

1. **Error: worker is being used**
   - This occurs when trying to use a worker that's already in use.
   - Solution: Either wait for the worker to finish or create a new worker.

2. **Error: pointer is nil**
   - This happens when trying to use a resource after it's been closed.
   - Solution: Make sure you're not closing resources before using them and not using resources after they're closed.

3. **Memory Usage Growing**
   - If memory usage keeps growing, you might not be closing resources properly.
   - Solution: Ensure all resources have a matching `Close()` call, ideally using `defer`.

4. **SVG Rendering Issues**
   - If text is not rendering correctly, check font configuration.
   - For images not loading, check the ResourcesDir setting.

### Debug Techniques

1. **Check Error Returns**: Always check error returns from function calls.
2. **Memory Profiling**: Use Go's built-in memory profiling to detect leaks.
3. **Validate SVG**: Ensure your SVG is valid and follows the SVG specification.

## Architecture

The architecture of resvg-go consists of several key components:

1. **Worker**: The main entry point that manages the WebAssembly runtime.
2. **FontDB**: Manages font loading and configuration.
3. **Tree**: Represents the parsed SVG document.
4. **Pixmap**: The rendering target for the SVG.
5. **Transform**: Applies geometric transformations to the SVG during rendering.

These components work together to provide a complete SVG rendering pipeline:

```
SVG Data → Tree → Apply Fonts → Apply Transform → Render to Pixmap → Encode to PNG
```

The WebAssembly module contains the compiled Rust code from the resvg library, which handles the actual rendering logic.

## Contributing

Contributions are welcome! Here's how you can contribute:

1. **Report Issues**: If you find a bug or have a feature request, please open an issue.
2. **Submit Pull Requests**: Feel free to submit PRs for bug fixes or new features.
3. **Improve Documentation**: Help improve the documentation or add examples.

Please follow these guidelines when contributing:

- Follow Go's code style and conventions.
- Add tests for new features.
- Update documentation for any changes.
- Ensure all tests pass before submitting PRs.

## Important Notes

- Workers are not goroutine-safe. Do not use the same worker concurrently in multiple goroutines.
- All resources (Worker, FontDB, Pixmap, Tree, etc.) must be cleaned up with Close() after use.
- Consider memory usage when processing large SVG files.
- The library follows the SVG 1.1 specification with some additional features from SVG 2.0.

## Acknowledgements

- [resvg](https://github.com/RazrFalcon/resvg) - an SVG rendering library written in Rust
- [wazero](https://github.com/tetratelabs/wazero) - the zero dependency WebAssembly runtime for Go developers

## License

For license information, please refer to the LICENSE file.
