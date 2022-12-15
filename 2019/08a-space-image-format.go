package main

import (
	"fmt"
	"log"
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

func histogram(data []int) []int {
	h := make([]int, 10)

	for _, v := range data {
		h[v]++
	}

	return h
}

const w = 25
const h = 6

func main() {
	var data string
	fmt.Scanln(&data)

	img := makeImage(data, w, h)

	var bestHistogram []int

	for l := 0; l < img.L; l++ {
		h := histogram(img.getLayer(l))
		if bestHistogram == nil || h[0] < bestHistogram[0] {
			bestHistogram = h
		}
	}

	fmt.Println(bestHistogram[1] * bestHistogram[2])
}
