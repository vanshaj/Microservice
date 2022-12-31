package main

import (
	"fmt"

	"github.com/google/go-containerregistry/pkg/crane"
)

func main() {
	img, err := crane.Pull("")
	if err != nil {
		panic(err)
	}
	layers, err := img.Layers()
	if err != nil {
		panic(err)
	}
	config, _ := img.ConfigFile()
	hist := config.History
	for _, val := range hist {
		fmt.Println(val)
	}
	fmt.Println("--------------------")
	for _, layer := range layers {
		hash, _ := layer.DiffID()
		fmt.Println(hash.Hex)
	}
}
