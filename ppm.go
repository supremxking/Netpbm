package main

import (
	"errors"
	"fmt"
	"math"
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

// ReadPPM reads a PPM image from a file and returns a struct that represents the image.
func ReadPPM(filename string) (*PPM, error) {
	var err error
	var magicNumber string = ""
	var width int
	var height int
	var maxval int
	var counter int
	var headersize int
	var splitfile []string
	file, err := os.ReadFile(filename)
	if err != nil {
	}
	if strings.Contains(string(file), "\r") { // if string of file contains return chariot
		splitfile = strings.SplitN(string(file), "\r\n", -1) // put in splitfile every data in string of file splitted if it contain \r\n
	} else {
		splitfile = strings.SplitN(string(file), "\n", -1)
	}
	for i, _ := range splitfile {
		if strings.Contains(splitfile[i], "P3") { //loop that check and confirm the magic number
			magicNumber = "P3"
		} else if strings.Contains(splitfile[i], "P6") {
			magicNumber = "P6"
		}
		if strings.HasPrefix(splitfile[i], "#") && maxval != 0 {
			headersize = counter
		}
		splitl := strings.SplitN(splitfile[i], " ", -1) // Split when there is space
		if width == 0 && height == 0 && len(splitl) >= 2 {
			width, err = strconv.Atoi(splitl[0])  // take width data / convert split string in int
			height, err = strconv.Atoi(splitl[1]) // take height data
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
			for j := 0; j < width*3; j += 3 { //loop with 3 incr to take data 3 by 3 R G B + condition to not go out of range
				r, _ := strconv.Atoi(splitdata[j])
				if r > maxval {
					r = maxval
				}
				g, _ := strconv.Atoi(splitdata[j+1])
				if g > maxval {
					g = maxval
				}
				b, _ := strconv.Atoi(splitdata[j+2])
				if b > maxval {
					b = maxval
				}
				data[i][j/3] = Pixel{R: uint8(r), G: uint8(g), B: uint8(b)}
			}
		}
	}
	return &PPM{data: data, width: width, height: height, magicNumber: magicNumber, max: uint8(maxval)}, err
}

// Size returns the width and height of the image.
func (ppm *PPM) Size() (int, int) {
	return ppm.width, ppm.height
}

// At returns the value of the pixel at (x, y).
func (ppm *PPM) At(x, y int) Pixel {
	return ppm.data[y][x]
}

// Set sets the value of the pixel at (x, y).
func (ppm *PPM) Set(x, y int, value Pixel) {
	if x >= 0 && x < ppm.width && y >= 0 && y < ppm.height {
		ppm.data[y][x] = value
	}
}

// Save saves the PPM image to a file and returns an error if there was a problem.
func (ppm *PPM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		file.Close()
		return err
	}
	_, err = fmt.Fprintln(file, ppm.magicNumber)
	if err != nil {
		file.Close()
		return err
	}
	_, err = fmt.Fprintln(file, ppm.width, ppm.height)
	if err != nil {
		file.Close()
		return err
	}
	_, err = fmt.Fprintln(file, ppm.max)
	if err != nil {
		file.Close()
		return err
	}

	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			if ppm.data[y][x].R > ppm.max || ppm.data[y][x].G > ppm.max || ppm.data[y][x].B > ppm.max {
				errors.New("data value is too high")
			} else {
				fmt.Fprint(file, ppm.data[y][x].R, ppm.data[y][x].G, ppm.data[y][x].B, " ")
			}
		}
		fmt.Fprintln(file)
	}
	return err
}

// Invert the color of the data for R G and B
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
	numRows := len(ppm.data)
	if numRows == 0 {
		return
	}
	for i := 0; i < numRows/2; i++ {
		ppm.data[i], ppm.data[numRows-i-1] = ppm.data[numRows-i-1], ppm.data[i]
	}
}

// SetMagicNumber sets the magic number of the PPM image.
func (ppm *PPM) SetMagicNumber(magicNumber string) {
	ppm.magicNumber = magicNumber
}

// SetMaxValue sets the max value of the PPM image.
func (ppm *PPM) SetMaxValue(maxValue uint8) {
	oldmax := ppm.max
	ppm.max = maxValue
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			// Convert each color component individually
			ppm.data[i][j].R = uint8(float64(ppm.data[i][j].R) * float64(ppm.max) / float64(oldmax))
			ppm.data[i][j].G = uint8(float64(ppm.data[i][j].G) * float64(ppm.max) / float64(oldmax))
			ppm.data[i][j].B = uint8(float64(ppm.data[i][j].B) * float64(ppm.max) / float64(oldmax))
		}
	}
}

// Rotate90CW rotates the PPM image 90Â° clockwise.
func (ppm *PPM) Rotate90CW() {
	NumRows := ppm.width
	NumColumns := ppm.height
	var newData [][]Pixel
	for i := 0; i < NumRows; i++ {
		newData = append(newData, make([]Pixel, NumColumns))
	}

	for i := 0; i < NumRows; i++ {
		for j := 0; j < NumColumns; j++ {
			newData[i][j] = ppm.data[NumColumns-j-1][i]
		}
	}
	ppm.data = newData
}

// ToPGM converts the PPM image to PGM.

func (ppm *PPM) ToPGM() *PGM {
	var newNumber string
	if ppm.magicNumber == "P3" { // Check magic number and change it
		newNumber = "P2"
	} else if ppm.magicNumber == "P6" {
		newNumber = "P5"
	}

	NumRows := ppm.width
	NumColumns := ppm.height
	var newData = make([][]uint8, NumColumns) // New Matrix to put the new data and allocate NumColums size

	for i := 0; i < NumColumns; i++ {
		newData[i] = make([]uint8, NumRows)
		for j := 0; j < NumRows; j++ {
			newData[i][j] = uint8((int(ppm.data[i][j].R) + int(ppm.data[i][j].G) + int(ppm.data[i][j].B)) / 3) //Convert the data
		}
	}

	return &PGM{data: newData, width: NumRows, height: NumColumns, magicNumber: newNumber, max: ppm.max}
}

// ToPBM converts the PPM image to PBM.

func (ppm *PPM) ToPBM() *PBM {
	var newNumber string
	if ppm.magicNumber == "P3" {
		newNumber = "P1"
	} else if ppm.magicNumber == "P6" {
		newNumber = "P4"
	}
	NumRows := ppm.width
	NumColumns := ppm.height
	var newData = make([][]bool, NumColumns)
	for i := 0; i < NumColumns; i++ {
		newData[i] = make([]bool, NumRows)
		for j := 0; j < NumRows; j++ {
			newData[i][j] = (uint8((int(ppm.data[i][j].R)+int(ppm.data[i][j].G)+int(ppm.data[i][j].B))/3) < ppm.max/2)
		}
	}
	return &PBM{data: newData, width: NumRows, height: NumColumns, magicNumber: newNumber}
}

type Point struct {
	X, Y int
}

// DrawLine draws a line between two points.
func (ppm *PPM) DrawLine(p1, p2 Point, color Pixel) {
	dx := float64(p2.X - p1.X) //Draw a Line with the bresenham algorithm
	dy := float64(p2.Y - p1.Y)
	steps := int(math.Max(math.Abs(dx), math.Abs(dy))) // new variable steps = int of the max absolute value f dx and dy

	xIncrement := dx / float64(steps)
	yIncrement := dy / float64(steps)

	x, y := float64(p1.X), float64(p1.Y)

	for i := 0; i <= steps; i++ {
		ppm.Set(int(x), int(y), color)
		x += xIncrement
		y += yIncrement
	}
}

// Draw an rectangle with drawline function
func (ppm *PPM) DrawRectangle(p1 Point, width, height int, color Pixel) {
	p2 := Point{p1.X + width, p1.Y} //Initalize the second point of the rectangle
	// Draw the rectangle width
	ppm.DrawLine(p1, p2, color)

	p3 := Point{p2.X, p2.Y + height}
	// Draw  the rectangle height
	ppm.DrawLine(p2, p3, color)

	p4 := Point{p1.X, p1.Y + height}
	ppm.DrawLine(p3, p4, color)

	ppm.DrawLine(p4, p1, color)
}

// DrawFilledRectangle draws a filled rectangle.
func (ppm *PPM) DrawFilledRectangle(p1 Point, width, height int, color Pixel) {
	ppm.DrawRectangle(p1, width, height, color)
	for j := p1.Y + 1; j < p1.Y+height; j++ {
		for i := p1.X + 1; i < p1.X+width; i++ {
			ppm.Set(i, j, color)
		}
	}
}

// DrawTriangle draws a triangle with Drawline func
func (ppm *PPM) DrawTriangle(p1, p2, p3 Point, color Pixel) {
	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p2, p3, color)
	ppm.DrawLine(p3, p1, color)
}
