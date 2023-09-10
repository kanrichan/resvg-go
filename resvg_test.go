package resvg

import (
	"context"
	"os"
	"testing"
)

var svg = []byte(
	`<svg width="100" height="100" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
    <rect id="rect1" x="10" y="10" width="80" height="80" fill="black"/>
</svg>`)

func TestRender(t *testing.T) {
	ctx, err := NewContext(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer ctx.Close()
	renderer, err := ctx.NewRenderer()
	if err != nil {
		t.Fatal(err)
	}
	defer renderer.Close()
	b, err := renderer.Render(svg)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile("out.png", b, 0644)
	if err != nil {
		t.Fatal(err)
	}
}
