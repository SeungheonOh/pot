package renderer

import (
	"bytes"
	"fmt"
	"image"
	"math"
)

type Character struct {
	Mask [4 * 8]int
	C    string
}

var ChSet = []Character{

	Character{Mask: [32]int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, C: "█"},
	Character{Mask: [32]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, C: " "},
	Character{Mask: [32]int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, C: "▀"},
	Character{Mask: [32]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1}, C: "▁"},
	Character{Mask: [32]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1}, C: "▂"},
	Character{Mask: [32]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, C: "▃"},
	Character{Mask: [32]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, C: "▄"},
	Character{Mask: [32]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, C: "▅"},
	Character{Mask: [32]int{0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, C: "▆"},
	Character{Mask: [32]int{0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, C: "▇"},
	Character{Mask: [32]int{1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1, 0}, C: "▉"},
	Character{Mask: [32]int{1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0}, C: "▋"},
	Character{Mask: [32]int{1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0}, C: "▍"},
	Character{Mask: [32]int{0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1}, C: "▐"},
	Character{Mask: [32]int{1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, C: "▔"},
	Character{Mask: [32]int{0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1}, C: "▕"},
	Character{Mask: [32]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0}, C: "▖"},
	Character{Mask: [32]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1}, C: "▗"},
	Character{Mask: [32]int{0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, C: "▘"},
	Character{Mask: [32]int{1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, C: "▙"},
	Character{Mask: [32]int{1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1}, C: "▚"},
	Character{Mask: [32]int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0}, C: "▛"},
	Character{Mask: [32]int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1}, C: "▜"},
	Character{Mask: [32]int{0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, C: "▝"},
	Character{Mask: [32]int{0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0}, C: "▞"},
	Character{Mask: [32]int{0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, C: "▟"},
	Character{Mask: [32]int{0, 0, 0, 1, 0, 0, 1, 1, 0, 0, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 0, 0, 1, 1, 0, 0, 1, 0, 0, 0}, C: "/"},
	Character{Mask: [32]int{1, 0, 0, 0, 1, 1, 0, 0, 0, 1, 0, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 0, 1, 0, 0, 0, 1, 1, 0, 0, 0, 1}, C: "\\"},
	Character{Mask: [32]int{0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0}, C: "|"},
	Character{Mask: [32]int{0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0}, C: "|"},
	Character{Mask: [32]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, C: "-"},
	Character{Mask: [32]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, C: "─"},
	Character{Mask: [32]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0}, C: "┏"},
	Character{Mask: [32]int{0, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 1, 1, 0}, C: "┇"},

	// ASCII
	/* It kindda looks bad
	Character{Mask: [32]int{1, 1, 1, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1}, C: "@"},
	Character{Mask: [32]int{1, 1, 1, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 1, 1, 0}, C: "["},
	Character{Mask: [32]int{0, 1, 1, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1, 1, 1}, C: "]"},
	Character{Mask: [32]int{0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 1, 1, 0}, C: "!"},
	Character{Mask: [32]int{0, 1, 1, 0, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1}, C: "A"},
	Character{Mask: [32]int{1, 1, 1, 0, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 1, 0, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 1, 0}, C: "B"},
	Character{Mask: [32]int{0, 1, 1, 0, 1, 0, 0, 1, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 1, 0, 1, 1, 0}, C: "C"},
	Character{Mask: [32]int{1, 1, 1, 0, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 1, 0}, C: "D"},
	Character{Mask: [32]int{1, 1, 1, 1, 1, 0, 0, 0, 1, 0, 0, 0, 1, 1, 1, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 1, 1, 1}, C: "E"},
	Character{Mask: [32]int{1, 1, 1, 1, 1, 0, 0, 0, 1, 0, 0, 0, 1, 1, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0}, C: "F"},
	Character{Mask: [32]int{0, 1, 1, 0, 1, 0, 0, 1, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1, 1, 1, 0, 0, 1, 1, 0, 0, 1, 0, 1, 1, 0}, C: "G"},
	Character{Mask: [32]int{1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1}, C: "H"},
	Character{Mask: [32]int{1, 1, 1, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 1, 1, 1, 0}, C: "I"},
	Character{Mask: [32]int{0, 0, 1, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 1, 0, 0, 1, 0, 1, 1, 0}, C: "J"},
	Character{Mask: [32]int{1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 1, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 0, 1, 0, 1, 0, 0, 1, 1, 0, 0, 1}, C: "K"},
	Character{Mask: [32]int{1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 1, 1, 1}, C: "L"},
	Character{Mask: [32]int{1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1}, C: "M"},
	Character{Mask: [32]int{1, 0, 0, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 0, 0, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1, 0, 0, 1, 1, 0, 0, 1}, C: "N"},
	Character{Mask: [32]int{0, 1, 1, 0, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 0, 1, 1, 0}, C: "O"},
	Character{Mask: [32]int{1, 1, 1, 0, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 1, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0}, C: "P"},
	Character{Mask: [32]int{0, 1, 1, 0, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 0, 1, 1, 0, 0, 0, 1, 1}, C: "Q"},
	Character{Mask: [32]int{1, 1, 1, 0, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 1, 0, 1, 1, 0, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 0, 1}, C: "R"},
	Character{Mask: [32]int{0, 1, 1, 0, 1, 0, 0, 1, 1, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 1, 0, 0, 0, 1, 1, 0, 0, 1, 0, 1, 1, 0}, C: "S"},
	Character{Mask: [32]int{1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 0, 1, 1, 1}, C: "U"},
	Character{Mask: [32]int{1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 0, 1, 1, 0, 0, 1, 1, 0}, C: "V"},
	Character{Mask: [32]int{1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 0, 1, 1, 0, 0, 1, 1, 0, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1}, C: "X"},
	Character{Mask: [32]int{1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0}, C: "Y"},
	Character{Mask: [32]int{1, 1, 1, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 1, 1, 1}, C: "Z"},
	*/
	/* Templete
	Character{
		Mask: [32]int{
			0, 0, 0, 0,
			0, 0, 0, 0,
			0, 0, 0, 0,
			0, 0, 0, 0,
			0, 0, 0, 0,
			0, 0, 0, 0,
			0, 0, 0, 0,
			0, 0, 0, 0,
		},
		C: " ",
	},
	*/
}

func init() {
	RendererMap["4x8"] = NewMask4x8
}

type Mask4x8 struct {
}

func NewMask4x8(options ...string) RenderEngine {
	return &Mask4x8{}
}

func (r *Mask4x8) Size(size image.Point) image.Point {
	return image.Point{size.X * 4, size.Y * 8}
}

func (r *Mask4x8) Render(img image.Image) (string, error) {
	var ret bytes.Buffer
	for i := 0; i < img.Bounds().Max.Y-8; i += 8 {
		for j := 0; j < img.Bounds().Max.X; j += 4 {
			histogram := [64]int{}
			histogramSum := [64][3]int{}
			top_result := [2]int{-1, -1}

			for SubY := 0; SubY < 8; SubY++ {
				for SubX := 0; SubX < 4; SubX++ {
					r, g, b, _ := img.At(j+SubX, i+SubY).RGBA()
					r /= 257
					b /= 257
					g /= 257
					index := r>>6*16 + g>>6*4 + b>>6
					histogram[index]++
					histogramSum[index][0] += int(r)
					histogramSum[index][1] += int(g)
					histogramSum[index][2] += int(b)
				}
			}
			for i, v := range histogram {
				if top_result[0] == -1 || v > histogram[top_result[0]] {
					top_result[1] = top_result[0]
					top_result[0] = i
				} else if top_result[1] == -1 || v > histogram[top_result[1]] {
					top_result[1] = i
				}
			}
			index := top_result[0]
			if histogram[index] != 0 {
				histogramSum[index][0] /= histogram[index]
				histogramSum[index][1] /= histogram[index]
				histogramSum[index][2] /= histogram[index]
			}
			index = top_result[1]
			if histogram[index] != 0 {
				histogramSum[index][0] /= histogram[index]
				histogramSum[index][1] /= histogram[index]
				histogramSum[index][2] /= histogram[index]
			}

			mask := [4 * 8]int{}
			for SubY := 0; SubY < 8; SubY++ {
				for SubX := 0; SubX < 4; SubX++ {
					r, g, b, _ := img.At(j+SubX, i+SubY).RGBA()
					r /= 257
					b /= 257
					g /= 257
					rdiff1 := math.Abs(float64(histogramSum[top_result[0]][0] - int(r)))
					gdiff1 := math.Abs(float64(histogramSum[top_result[0]][1] - int(g)))
					bdiff1 := math.Abs(float64(histogramSum[top_result[0]][2] - int(b)))
					rdiff2 := math.Abs(float64(histogramSum[top_result[1]][0] - int(r)))
					gdiff2 := math.Abs(float64(histogramSum[top_result[1]][1] - int(g)))
					bdiff2 := math.Abs(float64(histogramSum[top_result[1]][2] - int(b)))

					if rdiff1+gdiff1+bdiff1 > rdiff2+gdiff2+bdiff2 {
						mask[SubY*4+SubX] = 0
					} else {
						mask[SubY*4+SubX] = 1
					}
				}
			}

			min := float64(math.MaxFloat64)
			minIndex := 0
			_ = minIndex
			for i, c := range ChSet {
				if CompMask(mask, c.Mask) < min {
					minIndex = i
					min = CompMask(mask, c.Mask)
				}
			}
			r1 := histogramSum[top_result[0]][0]
			g1 := histogramSum[top_result[0]][1]
			b1 := histogramSum[top_result[0]][2]
			r2 := histogramSum[top_result[1]][0]
			g2 := histogramSum[top_result[1]][1]
			b2 := histogramSum[top_result[1]][2]

			fmt.Fprintf(&ret, "\033[48;2;%d;%d;%dm\033[38;2;%d;%d;%dm"+ChSet[minIndex].C+"\033[49m\033[39m",
				r1,
				g1,
				b1,
				r2,
				g2,
				b2)
		}
		fmt.Fprintf(&ret, "\n")
	}
	return ret.String(), nil
}

func CompMask(m1, m2 [32]int) float64 {
	score := float64(0)
	for i := 0; i < 4*8; i++ {
		if m2[i] == m1[i] {
			score++
		}
	}
	return score
}
