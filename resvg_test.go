package resvg

import (
	"context"
	"os"
	"testing"

	"github.com/kanrichan/resvg-go/internal"
)

func TestMemory(t *testing.T) {
	ctx := context.Background()
	worker, err := NewDefaultWorker(ctx)
	if err != nil {
		t.Fatal(err)
	}
	p, err := internal.MemoryMalloc(worker.ctx, worker.mod, 8)
	if err != nil {
		t.Fatal(err)
	}
	err = internal.MemoryFree(worker.ctx, worker.mod, p, 8)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFontDB(t *testing.T) {
	ctx := context.Background()
	worker, err := NewDefaultWorker(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer worker.Close()
	fontdb, err := worker.NewFontDBDefault()
	if err != nil {
		t.Fatal(err)
	}
	ttf, err := os.ReadFile("./testdata/arial.ttf")
	if err != nil {
		t.Fatal(err)
	}
	err = fontdb.LoadFontData(ttf)
	if err != nil {
		t.Fatal(err)
	}
	err = fontdb.LoadFontFile("./testdata/arial.ttf")
	if err != nil {
		t.Fatal(err)
	}
	err = fontdb.LoadFontsDir("./testdata")
	if err != nil {
		t.Fatal(err)
	}
	num, err := fontdb.Len()
	if err != nil {
		t.Fatal(err)
	}
	if num != 3 {
		t.Fatal("fontdb len must be 3")
	}
	err = fontdb.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestPixmap(t *testing.T) {
	ctx := context.Background()
	worker, err := NewDefaultWorker(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer worker.Close()
	pixmap, err := worker.NewPixmap(100, 100)
	if err != nil {
		t.Fatal(err)
	}
	data, err := pixmap.EncodePNG()
	if err != nil {
		t.Fatal(err)
	}
	if data[1] != 80 || data[2] != 78 || data[3] != 71 {
		t.Fatal("illegal PNG")
	}
	err = pixmap.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestTree(t *testing.T) {
	ctx := context.Background()
	worker, err := NewDefaultWorker(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer worker.Close()
	svg, err := os.ReadFile("./testdata/beach.svg")
	if err != nil {
		t.Fatal(err)
	}
	tree, err := worker.NewTreeFromData(svg, &Options{
		ResourcesDir:       "",
		Dpi:                96.0,
		FontFamily:         "Times New Roman",
		FontSize:           12.0,
		Languages:          []string{"en"},
		ShapeRenderingMode: ShapeRenderingModeGeometricPrecision,
		TextRenderingMode:  TextRenderingModeOptimizeLegibility,
		ImageRenderingMode: ImageRenderingModeOptimizeQuality,
		DefaultSizeWidth:   100.0,
		DefaultSizeHeight:  100.0,
	})
	if err != nil {
		t.Fatal(err)
	}
	width, height, err := tree.GetSize()
	if err != nil {
		t.Fatal(err)
	}
	if width != 512.0 || height != 512.0 {
		t.Fatal("width and height should be 512.0")
	}
	fontdb, err := worker.NewFontDBDefault()
	if err != nil {
		t.Fatal(err)
	}
	defer fontdb.Close()
	err = tree.ConvertText(fontdb)
	if err != nil {
		t.Fatal(err)
	}
	pixmap, err := worker.NewPixmap(512, 512)
	if err != nil {
		t.Fatal(err)
	}
	defer pixmap.Close()
	err = tree.Render(TransformIdentity(), pixmap)
	if err != nil {
		t.Fatal(err)
	}
	data, err := pixmap.EncodePNG()
	if err != nil {
		t.Fatal(err)
	}
	if data[1] != 80 || data[2] != 78 || data[3] != 71 {
		t.Fatal("illegal PNG")
	}
	err = tree.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestRender(t *testing.T) {
	var svg = []byte(
		`<svg width="100" height="100" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
			<rect id="rect1" x="10" y="10" width="80" height="80" fill="black"/>
		</svg>`)
	worker, err := NewDefaultWorker(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer worker.Close()
	data, err := worker.Render(svg)
	if err != nil {
		t.Fatal(err)
	}
	if data[1] != 80 || data[2] != 78 || data[3] != 71 {
		t.Fatal("illegal PNG")
	}
}
