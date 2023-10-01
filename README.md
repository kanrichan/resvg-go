<div align="center">
  <img src=".github/hua_nobg_512.gif" alt="椛" width = "400">
  <br>

  <h1>resvg-go</h1>
  <img src="https://counter.seku.su/cmoe?name=resvgo&theme=r34" /><br>
  A SVG renderer written in Go & WASM depended on resvg without CGO<br><br>
  
</div>


## Usage

### Render by default
```go
// initialize and don't forget to close!
// worker are not goroutine-safe!
worker, _ := NewDefaultWorker(context.Background())
defer worker.Close()

// render the SVG as a PNG!
png, _ := worker.Render(svg)
```

### Render with options
```go
// initialize and don't forget to close!
// worker are not goroutine-safe!
worker, _ := NewDefaultWorker(context.Background())
defer worker.Close()

// render the SVG as a PNG!
fontdb, _ := worker.NewFontDBDefault()
defer fontdb.Close()
fontdb.LoadFontData(ttf)

pixmap, _ := worker.NewPixmap(512, 512)
defer pixmap.Close()

tree, _ := worker.NewTreeFromData(svg, &Options{})
defer tree.Close()
tree.ConvertText(fontdb)
tree.Render(TransformIdentity(), pixmap)

png, _ := pixmap.EncodePNG()
```


## Thanks
- [resvg](https://github.com/RazrFalcon/resvg) - an SVG rendering library written in Rust
- [wazero](https://github.com/tetratelabs/wazero) - the zero dependency WebAssembly runtime for Go developers
