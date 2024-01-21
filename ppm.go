package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PPM struct {
	data          [][]Pixel
	width, height int
	magicNumber   string
	max           uint8
}

type Pixel struct {
	R, G, B uint8
}

func ReadPPM(filename string) (*PPM, error) {
	var err error
	magicNumber := ""
	var width, height, maxval, counter, headersize int

	file, err := os.ReadFile(filename)
	if err != nil {
	}
	splitfile := strings.SplitN(string(file), "\n", -1)
	for i, _ := range splitfile {
		if strings.Contains(splitfile[i], "P3") {
			magicNumber = "P3"
		} else if strings.Contains(splitfile[i], "P6") {
			magicNumber = "P6"
		}
		if strings.HasPrefix(splitfile[i], "#") && maxval != 0 {
			headersize = counter
		}
		splitl := strings.SplitN(splitfile[i], " ", -1)
		if width == 0 && height == 0 && len(splitl) >= 2 {
			width, err = strconv.Atoi(splitl[0])
			height, err = strconv.Atoi(splitl[1])
			headersize = counter
		}
		if maxval == 0 && width != 0 {
			maxval, err = strconv.Atoi(splitfile[i])
			headersize = counter
		}
		counter++

	}

	data := make([][]Pixel, height)

	for j := 0; j < height; j++ {
		data[j] = make([]Pixel, width)
	}
	var splitdata []string

	if counter > headersize {
		for i := 0; i < height; i++ {
			splitdata = strings.SplitN(splitfile[headersize+1+i], " ", -1)
			for j := 0; j < width*3; j += 3 {
				r, _ := strconv.Atoi(splitdata[j])
				g, _ := strconv.Atoi(splitdata[j+1])
				b, _ := strconv.Atoi(splitdata[j+2])
				data[i][j/3] = Pixel{R: uint8(r), G: uint8(g), B: uint8(b)}
			}
		}
	}
	return &PPM{data: data, width: width, height: height, magicNumber: magicNumber, max: uint8(maxval)}, err
}

func display(data [][]Pixel) {
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data[0]); j++ {
			fmt.Print(data[i][j], " ")
		}
		fmt.Println()
	}
}

// Size returns the width and height of the image.
func (ppm *PPM) Size() (int, int) {
	return ppm.height, ppm.width
}

// At returns the value of the pixel at (x, y).
func (ppm *PPM) At(x, y int) Pixel {
	return ppm.data[y][x]
}

// Set sets the value of the pixel at (x, y).
func (ppm *PPM) Set(x, y int, value Pixel) {
	ppm.data[x][y] = value
}

// Invert inverts the colors of the PPM image.
func (ppm *PPM) Invert() {
	for i := 0; i < len(ppm.data); i++ {
		for j := 0; j < len(ppm.data[0]); j++ {
			ppm.data[i][j].R = ppm.max - ppm.data[i][j].R
			ppm.data[i][j].G = ppm.max - ppm.data[i][j].G
			ppm.data[i][j].B = ppm.max - ppm.data[i][j].B
		}
	}
}

// Flip flips the PPM image horizontally.
func (ppm *PPM) Flip() {
	// Height = Colums
	// Width = Rows
	NumRows := ppm.width
	NumColums := ppm.height
	for i := 0; i < NumRows; i++ {
		for j := 0; j < NumColums/2; j++ {
			ppm.data[i][j], ppm.data[i][NumColums-j-1] = ppm.data[i][NumColums-j-1], ppm.data[i][j]
		}
	}
}

// Flop flops the PPM image vertically.
func (ppm *PPM) Flop() {
	// Height = Colums
	// Width = Rows
	NumRows := ppm.width
	for i := 0; i < NumRows/2; i++ {
		ppm.data[i], ppm.data[NumRows-i-1] = ppm.data[NumRows-i-1], ppm.data[i]
	}
}

func (ppm *PPM) SetMagicNumber(magicNumber string) {
	ppm.magicNumber = magicNumber
}

// SetMaxValue sets the max value of the PPM image.
func (ppm *PPM) SetMaxValue(maxValue uint8) {
	oldMax := ppm.max
	ppm.max = maxValue
	for i := 0; i < len(ppm.data); i++ {
		for j := 0; j < len(ppm.data[0]); j++ {
			ppm.data[i][j].R = uint8(float64(ppm.data[i][j].R) * float64(ppm.max) / float64(oldMax))
			ppm.data[i][j].G = uint8(float64(ppm.data[i][j].G) * float64(ppm.max) / float64(oldMax))
			ppm.data[i][j].B = uint8(float64(ppm.data[i][j].B) * float64(ppm.max) / float64(oldMax))
		}
	}
}

// Rotate90CW rotates the PPM image 90° clockwise.
func (ppm *PPM) Rotate90CW() {
	// Height = Colums = Colonne vers le bas
	// Width = Rows = Ligne vers la droite
	NumRows := ppm.width
	NumColums := ppm.height
	for i := 0; i < NumColums; i++ {
		for j := i + 1; j < NumRows; j++ {
			vartemp := ppm.data[i][j]
			ppm.data[i][j] = ppm.data[j][i]
			ppm.data[j][i] = vartemp
		}
	}
	for i := 0; i < NumColums; i++ {
		for j := 0; j < NumRows/2; j++ {
			vartemp := ppm.data[i][j]
			ppm.data[i][j] = ppm.data[i][NumRows-j-1]
			ppm.data[i][NumRows-j-1] = vartemp
		}
	}
}

// ToPGM converts the PPM image to PGM.
func (ppm *PPM) ToPGM() *PGM {
	// Height = Colums = Colonne vers le bas
	// Width = Rows = Ligne vers la droite
	var newmagicnumber string

	if ppm.magicNumber == "P3" {
		newmagicnumber = "P2"
	} else if ppm.magicNumber == "P6" {
		newmagicnumber = "P5"
	}
	Numrows := ppm.width
	NumColumns := ppm.height
	var newdata = make([][]uint8, NumColumns)
	for i := 0; i < NumColumns; i++ {
		newdata[i] = make([]uint8, Numrows)
		for j := 0; j < Numrows; j++ {
			{
				newdata[i][j] = uint8((int(ppm.data[i][j].R) + int(ppm.data[i][j].G) + int(ppm.data[i][j].B)) / 3)
			}
		}
	}
	return &PGM{data: newdata, width: Numrows, height: NumColumns, max: ppm.max, magicNumber: newmagicnumber}
}

// ToPBM converts the PPM image to PBM.
func (ppm *PPM) ToPBM() *PBM {
	var newmagicnumber string

	if ppm.magicNumber == "P3" {
		newmagicnumber = "P1"
	} else if ppm.magicNumber == "P6" {
		newmagicnumber = "P4"
	}
	Numrows := ppm.width
	NumColumns := ppm.height
	var newdata = make([][]bool, NumColumns)
	for i := 0; i < NumColumns; i++ {
		newdata[i] = make([]bool, Numrows)
		for j := 0; j < Numrows; j++ {
			newdata[i][j] = uint8((int(ppm.data[i][j].R)+int(ppm.data[i][j].G)+int(ppm.data[i][j].B))/3) < ppm.max/2
		}
	}
	return &PBM{data: newdata, width: Numrows, height: NumColumns, magicNumber: newmagicnumber}
}

type Point struct {
	X, Y int
}

// DrawLine draws a line between two points.
func (ppm *PPM) DrawLine(p1, p2 Point, color Pixel) {

}
