package resvg

import (
	"os"
	"testing"
)

var svg = []byte(
	`<svg width="100" height="100" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
    <rect id="rect1" x="10" y="10" width="80" height="80" fill="black"/>
</svg>`)

func TestRustBytes(t *testing.T) {
	inst, err := NewResvg()
	if err != nil {
		t.Fatal(err)
	}
	defer inst.Close()

	rb, err := inst.NewRustBytes(12)
	if err != nil {
		t.Fatal(err)
	}
	err = rb.WriteString("hello world!")
	if err != nil {
		t.Fatal(err)
	}
	if err := rb.Free(); err != nil {
		t.Fatal(err)
	}
}

func TestRustBytesPointer(t *testing.T) {
	inst, err := NewResvg()
	if err != nil {
		t.Fatal(err)
	}
	defer inst.Close()

	rb, err := inst.NewRustBytesPointer()
	if err != nil {
		t.Fatal(err)
	}
	if err := rb.Free(); err != nil {
		t.Fatal(err)
	}
}

func TestResvgRender(t *testing.T) {
	inst, err := NewResvg()
	if err != nil {
		t.Fatal(err)
	}
	defer inst.Close()

	out, err := inst.DefaultResvgRenderToPNG(svg)
	if err != nil {
		t.Fatal(err)
	}

	if err = os.WriteFile("out.png", out, 0644); err != nil {
		t.Fatal(err)
	}
}
