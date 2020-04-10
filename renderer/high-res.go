package renderer

import (
  "bytes"
	"fmt"
	"image"
	"math"
	"sync"

	cv "gocv.io/x/gocv"
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
	RendererMap["high-res"] = HighRes
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

func HighRes(img cv.Mat, size image.Point) (string, error) {
	var wg sync.WaitGroup
	buffer := make([]string, (size.X)*(size.Y))
	size.X *= 4
	size.Y *= 4
	cv.Resize(img, &img, size, 0, 0, 1)
	cv.GaussianBlur(img, &img, image.Point{3, 3}, 0, 0, 1)

	imgPtr := img.DataPtrUint8()

	fmt.Print("\033[0;0H")

	for y := 0; y < img.Rows(); y += 8 {
		for x := 0; x < img.Cols()*3; x += 3 * 4 {
			if y+8 >= img.Rows() {
				continue
			}
			// Rendering with gorutines, saving cells to buffer
			wg.Add(1)
			go func(x, y int) {
				histogram := [64]int{}
				histogramSum := [64][3]int{}
				top_result := [2]int{-1, -1}

				for SubY := 0; SubY < 8; SubY++ {
					for SubX := 0; SubX < 4*3; SubX += 3 {
						index := (y+SubY)*img.Cols()*3 + SubX + x
						r := int(imgPtr[index])
						g := int(imgPtr[index+1])
						b := int(imgPtr[index+2])
						index = r>>6*16 + g>>6*4 + b>>6
						histogram[index]++
						histogramSum[index][0] += r
						histogramSum[index][1] += g
						histogramSum[index][2] += b
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
					for SubX := 0; SubX < 4*3; SubX += 3 {
						index := (y+SubY)*img.Cols()*3 + SubX + x
						r := int(imgPtr[index])
						g := int(imgPtr[index+1])
						b := int(imgPtr[index+2])
						rdiff1 := math.Abs(float64(histogramSum[top_result[0]][0] - r))
						gdiff1 := math.Abs(float64(histogramSum[top_result[0]][1] - g))
						bdiff1 := math.Abs(float64(histogramSum[top_result[0]][2] - b))
						rdiff2 := math.Abs(float64(histogramSum[top_result[1]][0] - r))
						gdiff2 := math.Abs(float64(histogramSum[top_result[1]][1] - g))
						bdiff2 := math.Abs(float64(histogramSum[top_result[1]][2] - b))

						if rdiff1+gdiff1+bdiff1 > rdiff2+gdiff2+bdiff2 {
							mask[SubY*4+SubX/3] = 0
						} else {
							mask[SubY*4+SubX/3] = 1
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

				buffer[y/8*(img.Cols()/4)+x/12] = fmt.Sprintf("\033[48;2;%d;%d;%dm\033[38;2;%d;%d;%dm"+ChSet[minIndex].C+"\033[49m\033[39m", b1, g1, r1, b2, g2, r2)
				wg.Done()
			}(x, y)
		}
	}
	wg.Wait()

  var retbuffer bytes.Buffer

	// printing buffer
	for y := 0; y < img.Rows()/8; y++ {
		for x := 0; x < img.Cols()/4; x++ {
			fmt.Fprintf(&retbuffer, buffer[y*(img.Cols()/4)+x])
		}
		if y != img.Rows()/8-1 {
			fmt.Fprintf(&retbuffer, "\n")
		}
	}
	return retbuffer.String(), nil
}
