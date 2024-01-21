package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PGM struct {
	data          [][]uint8
	width, height int
	magicNumber   string
	max           uint8
}

func main() {

}

// Read PGM file P2 Only
func ReadPGM(filename string) (*PGM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var width, height, max int
	var data [][]uint8

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	magicNumber := scanner.Text()
	if magicNumber != "P2" && magicNumber != "P5" {
		return nil, errors.New("type de fichier non pris en charge")
	}

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "#") { //if the line don't contain # it check for width and height with the format asked
			_, err := fmt.Sscanf(line, "%d %d", &width, &height)
			if err == nil {
				break
			} else {
				fmt.Println("Largeur ou hauteur invalide :", err)
			}
		}
	}

	scanner.Scan()
	max, err = strconv.Atoi(scanner.Text()) //Check max value
	if err != nil {
		return nil, errors.New("valeur maximale de pixel invalide")
	}

	for scanner.Scan() {
		line := scanner.Text() // take the data and stock it
		if magicNumber == "P2" {
			row := make([]uint8, 0)
			for _, char := range strings.Fields(line) {
				pixel, err := strconv.Atoi(char)
				if err != nil {
					fmt.Println("Erreur de conversion en entier :", err)
				}
				if pixel >= 0 && pixel <= max {
					row = append(row, uint8(pixel))
				} else {
					fmt.Println("Valeur de pixel invalide :", pixel)
				}
			}
			data = append(data, row)
		}
	}

	return &PGM{ // return the new value of every thing in struct
		data:        data,
		width:       width,
		height:      height,
		magicNumber: magicNumber,
		max:         uint8(max),
	}, nil
}

func (pgm *PGM) Size() (int, int) {
	return pgm.width, pgm.height
}

// At returns the value of the pixel at (x, y).
func (pgm *PGM) At(x, y int) uint8 {
	return pgm.data[x][y]
}

// Set sets the value of the pixel at (x, y).
func (pgm *PGM) Set(x, y int, value uint8) {
	pgm.data[x][y] = value
}

// Invert inverts the colors of the PGM image.
func (pgm *PGM) Invert() {
	for i := 0; i < len(pgm.data); i++ {
		for j := 0; j < len(pgm.data[i]); j++ {
			pgm.data[i][j] = pgm.max - pgm.data[i][j]
		}
	}
}

// Flip flips the PGM image horizontally.
func (pgm *PGM) Flip() {
	NumRows := pgm.width
	Numcolums := pgm.height
	for i := 0; i < NumRows; i++ {
		for j := 0; j < Numcolums/2; j++ {
			pgm.data[i][j], pgm.data[i][Numcolums-j-1] = pgm.data[i][Numcolums-j-1], pgm.data[i][j]
		}
	}
}

// Flop flops the PGM image vertically.
func (pgm *PGM) Flop() {
	numRows := len(pgm.data)
	if numRows == 0 {
		return
	}
	for i := 0; i < numRows/2; i++ {
		pgm.data[i], pgm.data[numRows-i-1] = pgm.data[numRows-i-1], pgm.data[i]
	}
}

// Set new magic number
func (pgm *PGM) SetMagicNumber(magicNumber string) {
	pgm.magicNumber = magicNumber
}

// Set new max value
func (pgm *PGM) SetMaxValue(maxValue uint8) {
	oldmax := pgm.max
	pgm.max = maxValue
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {

			pgm.data[i][j] = pgm.data[i][j] * uint8(5) / oldmax
		}
	}

}

// Rotate90CW rotates the PGM image 90° clockwise.
func (pgm *PGM) Rotate90CW() {
	NumRows := pgm.width
	NumColums := pgm.height
	for i := 0; i < len(pgm.data); i++ { // First the transpose of the data
		var temp uint8
		for j := i + 1; j < len(pgm.data[0]); j++ {
			temp = pgm.data[i][j]
			pgm.data[i][j] = pgm.data[j][i]
			pgm.data[j][i] = temp
		}
	}
	for i := 0; i < NumColums; i++ {
		for j := 0; j < NumRows/2; j++ { // And then the reverse of the data after transpose
			temp := pgm.data[i][j]
			pgm.data[i][j] = pgm.data[i][NumRows-j-1]
			pgm.data[i][NumRows-j-1] = temp
		}
	}
}

// ToPBM converts the PGM image to PBM.
func (pgm *PGM) ToPBM() *PBM {
	if pgm.magicNumber == "P2" {
		pgm.magicNumber = "P1"
	} else if pgm.magicNumber == "P5" {
		pgm.magicNumber = "P4"
	}
	Numrows := pgm.width
	NumColumns := pgm.height
	var newdata = make([][]bool, NumColumns) // allocate size to newdata bool
	for i := 0; i < NumColumns; i++ {
		newdata[i] = make([]bool, Numrows)
		for j := 0; j < Numrows; j++ {
			if pgm.data[i][j] < pgm.max/2 { //In the newdata created for each uint8 in pgm.data lower than max/2 it stack true in it and else false
				newdata[i][j] = true
			} else {
				newdata[i][j] = false
			}
		}
	}
	return &PBM{data: newdata, magicNumber: pgm.magicNumber, width: Numrows, height: NumColumns}
}
