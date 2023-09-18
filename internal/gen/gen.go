package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
)

const (
	targetwasm = "target/wasm32-wasi/release/resvg.wasm"
)

func main() {
	fmt.Println("Run cargo build --release --target wasm32-wasi ...")
	carcmd := exec.Command("cargo", "build", "--release", "--target", "wasm32-wasi")
	carcmd.Stdout = os.Stdout
	carcmd.Stderr = os.Stderr
	err := carcmd.Run()
	if err != nil {
		panic(err)
	}
	tgt, err := os.Open(targetwasm)
	if err != nil {
		panic(err)
	}
	defer tgt.Close()
	fmt.Println("Pack", targetwasm, "to gzip ...")
	target := path.Join("internal", path.Base(targetwasm)) + ".gz"
	f, err := os.Create(target)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	zw, _ := gzip.NewWriterLevel(f, gzip.BestCompression)
	defer zw.Close()
	zw.Header.Name = "resvg.wasm"
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(zw, tgt)
	if err != nil {
		panic(err)
	}
	fmt.Println("Success.")
}
