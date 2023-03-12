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
	inst, err := DefaultResvg()
	if err != nil {
		t.Fatal(err)
	}
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
	inst, err := DefaultResvg()
	if err != nil {
		t.Fatal(err)
	}
	rb, err := inst.NewRustBytesPointer()
	if err != nil {
		t.Fatal(err)
	}
	if err := rb.Free(); err != nil {
		t.Fatal(err)
	}
}

func TestResvgRender(t *testing.T) {
	inst, err := DefaultResvg()
	if err != nil {
		t.Fatal(err)
	}
	out, err := inst.DefaultResvgRenderToPNG(svg)
	if err != nil {
		t.Fatal(err)
	}
	fi, err := os.OpenFile("out.png", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer fi.Close()
	fi.Write(out)
}
