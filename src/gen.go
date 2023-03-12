package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
)

const (
	targetwasm = "target/wasm32-unknown-unknown/release/resvg_go.wasm"
)

func main() {
	fmt.Println("Run cargo build --release --target wasm32-unknown-unknown ...")
	carcmd := exec.Command("cargo", "build", "--release", "--target", "wasm32-unknown-unknown")
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
	fmt.Println("Pack", targetwasm, "to zip ...")
	f, err := os.Create(targetwasm + ".zip")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	zw := zip.NewWriter(f)
	defer zw.Close()
	w, err := zw.Create("resvg_go.wasm")
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(w, tgt)
	if err != nil {
		panic(err)
	}
	fmt.Println("Success.")
}
