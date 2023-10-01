package resvg

import (
	"context"
	"os"
	"testing"
)

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
	tree, err := worker.NewTreeFromData(svg, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer tree.Close()
	width, height, err := tree.GetSize()
	if err != nil {
		t.Fatal(err)
	}
	pixmap, err := worker.NewPixmap(uint32(width), uint32(height))
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
	err = os.WriteFile("out.png", data, 0644)
	if err != nil {
		t.Fatal(err)
	}
}
