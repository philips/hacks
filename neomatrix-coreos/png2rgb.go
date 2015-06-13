// This example demonstrates decoding a JPEG image and examining its pixels.
package main

import (
	"fmt"
	"image"
	"log"
	"os"

	_ "image/png"

	"github.com/lucasb-eyer/go-colorful"
)

func main() {
	// Decode the JPEG data. If reading from file, create a reader with
	//
	reader, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	m, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	bounds := m.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			ra, ga, ba, _ := m.At(x, y).RGBA()
			c := colorful.Color{float64(ra) / 65536.0, float64(ga) / 65536.0, float64(ba) / 65536.0}
			r, g, b := c.RGB255()
			fmt.Printf("0x%x, 0x%x, 0x%x\n", r, g, b)
		}
	}

}
