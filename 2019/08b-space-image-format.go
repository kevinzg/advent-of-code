package main

import (
	"fmt"
	"log"
	"strings"
)

type image struct {
	data    []int
	w, h, L int
}

func makeImage(data string, w, h int) image {
	n := len(data)
	pixelsPerLayer := w * h

	if n%pixelsPerLayer != 0 {
		log.Fatal("Wrong size for image.")
	}

	layers := n / pixelsPerLayer

	intData := make([]int, n)

	for i, s := range data {
		intData[i] = int(s - '0')
	}

	return image{data: intData, w: w, h: h, L: layers}
}

func (img *image) getIndex(layer, i, j int) int {
	pixelsPerLayer := img.w * img.h
	return pixelsPerLayer*layer + i*img.w + j
}

func (img *image) getLayer(layer int) []int {
	p := img.getIndex(layer, 0, 0)
	q := img.getIndex(layer+1, 0, 0)
	return img.data[p:q]
}

func (img *image) render() string {
	canvas := make([]int, img.h*img.w)

	for l := 0; l < img.L; l++ {
		layer := img.L - 1 - l
		for i := range canvas {
			pixel := img.getLayer(layer)[i]
			if pixel != 2 {
				canvas[i] = pixel
			}
		}
	}

	colors := []string{"█", " ", "░"}

	var stringCanvas strings.Builder

	for i := 0; i < img.h; i++ {
		if i != 0 {
			stringCanvas.WriteString("\n")
		}

		for j := 0; j < img.w; j++ {
			pixel := canvas[i*img.w+j]
			color := colors[pixel]
			stringCanvas.WriteString(color)
		}
	}

	return stringCanvas.String()
}

const w = 25
const h = 6

func main() {
	var data string
	fmt.Scanln(&data)

	img := makeImage(data, w, h)

	fmt.Println(img.render())
}
