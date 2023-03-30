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
inst, _ := resvg.NewResvg()
defer inst.Close()

// render the SVG as a PNG!
png, _ := inst.DefaultResvgRenderToPNG(svg)
```
