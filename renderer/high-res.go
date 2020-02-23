package renderer

import (
	"fmt"
	"image"
	"math"

	cv "gocv.io/x/gocv"
)

type Character struct {
	Mask [4 * 8]int
	C    string
}

var ChSet = []Character{
	Character{Mask: [4 * 8]int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, C: "█"},
	Character{Mask: [4 * 8]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, C: " "},
	Character{Mask: [4 * 8]int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, C: "▀"},
	Character{Mask: [4 * 8]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1}, C: "▁"},
	Character{Mask: [4 * 8]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1}, C: "▂"},
	Character{Mask: [4 * 8]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, C: "▃"},
	Character{Mask: [4 * 8]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, C: "▄"},
	Character{Mask: [4 * 8]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, C: "▅"},
	Character{Mask: [4 * 8]int{0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, C: "▆"},
	Character{Mask: [4 * 8]int{0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, C: "▇"},
	Character{Mask: [4 * 8]int{1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1, 0}, C: "▉"},
	Character{Mask: [4 * 8]int{1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0}, C: "▋"},
	Character{Mask: [4 * 8]int{1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0}, C: "▍"},
	Character{Mask: [4 * 8]int{0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1}, C: "▐"},
	Character{Mask: [4 * 8]int{1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, C: "▔"},
	Character{Mask: [4 * 8]int{0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1}, C: "▕"},
	Character{Mask: [4 * 8]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0}, C: "▖"},
	Character{Mask: [4 * 8]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1}, C: "▗"},
	Character{Mask: [4 * 8]int{0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, C: "▘"},
	Character{Mask: [4 * 8]int{1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, C: "▙"},
	Character{Mask: [4 * 8]int{1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1}, C: "▚"},
	Character{Mask: [4 * 8]int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0}, C: "▛"},
	Character{Mask: [4 * 8]int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1}, C: "▜"},
	Character{Mask: [4 * 8]int{0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, C: "▝"},
	Character{Mask: [4 * 8]int{0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0}, C: "▞"},
	Character{Mask: [4 * 8]int{0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, C: "▟"},
	Character{Mask: [4 * 8]int{1, 1, 1, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1}, C: "@"},
	Character{Mask: [4 * 8]int{0, 0, 0, 1, 0, 0, 1, 1, 0, 0, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 0, 0, 1, 1, 0, 0, 1, 0, 0, 0}, C: "/"},
	Character{Mask: [4 * 8]int{1, 0, 0, 0, 1, 1, 0, 0, 0, 1, 0, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 0, 1, 0, 0, 0, 1, 1, 0, 0, 0, 1}, C: "\\"},
}

func init() {
	RendererMap["high-res"] = HighPix
}

func CompMask(m1, m2 [4 * 8]int) float64 {
	score := float64(0)
	for i := 0; i < 4*8; i++ {
		if m2[i] == m1[i] {
			score++
		}
	}
	return score
}

func HighPix(img cv.Mat, size image.Point) error {
	size.X *= 4
	size.Y *= 8
	size.Y -= 16
	cv.Resize(img, &img, size, 0, 0, 1)
	cv.GaussianBlur(img, &img, image.Point{3, 3}, 2, 2, 0)

	imgPtr := img.DataPtrUint8()

	fmt.Print("\033[0;0H")

	for y := 0; y < img.Rows()/2; y += 8 {
		for x := 0; x <= img.Cols()*3; x += 3 * 4 {
			histogram := [64]int{}
			histogramSum := [64][3]int{}
			for SubY := 0; SubY < 8; SubY++ {
				for SubX := 0; SubX < 4*3; SubX += 3 {
					r := int(imgPtr[y*img.Cols()*3+(y+SubY)*img.Cols()*3+SubX+x])
					g := int(imgPtr[y*img.Cols()*3+(y+SubY)*img.Cols()*3+SubX+x+1])
					b := int(imgPtr[y*img.Cols()*3+(y+SubY)*img.Cols()*3+SubX+x+2])
					index := r>>6*16 + g>>6*4 + b>>6
					histogram[index]++
					histogramSum[index][0] += r
					histogramSum[index][1] += g
					histogramSum[index][2] += b
				}
			}
			top_result := [2]int{}
			max := 0
			max_index := 0
			for i := 0; i < 64; i++ {
				if histogram[i] > max {
					max = histogram[i]
					max_index = i
				}
			}
			top_result[0] = max_index
			max = 0
			max_index = 0
			for i := 0; i < 64; i++ {
				if histogram[i] > max && i != top_result[0] {
					max = histogram[i]
					max_index = i
				}
			}
			top_result[1] = max_index
			for i := 0; i < 64; i++ {
				if histogram[i] == 0 {
					continue
				}
				histogramSum[i][0] /= histogram[i]
				histogramSum[i][1] /= histogram[i]
				histogramSum[i][2] /= histogram[i]
			}

			mask := [4 * 8]int{}
			for SubY := 0; SubY < 8; SubY++ {
				for SubX := 0; SubX < 4*3; SubX += 3 {
					r := int(imgPtr[y*img.Cols()*3+(y+SubY)*img.Cols()*3+SubX+x])
					g := int(imgPtr[y*img.Cols()*3+(y+SubY)*img.Cols()*3+SubX+x+1])
					b := int(imgPtr[y*img.Cols()*3+(y+SubY)*img.Cols()*3+SubX+x+2])
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

			fmt.Printf("\033[48;2;%d;%d;%dm\033[38;2;%d;%d;%dm"+ChSet[minIndex].C+"\033[49m\033[39m", b1, g1, r1, b2, g2, r2)
		}
		fmt.Printf("\n")
	}
	return nil
}
