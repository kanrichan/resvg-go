<div align="center">
  <img src=".github/hua_nobg_512.gif" alt="椛" width = "400">
  <br>

  <h1>resvg-go</h1>
  <img src="https://counter.seku.su/cmoe?name=resvgo&theme=r34" /><br>
  A SVG render written in Go & WASM depended on resvg without CGO<br><br>
  
</div>


## Example

```go
// initialize and don't forget to close!
worker, _ := NewDefaultWorker(context.Background())
defer worker.Close()
renderer, _ := worker.NewRenderer()
defer renderer.Close()

// render the SVG as a PNG!
png, _ := renderer.Render(svg)
```
