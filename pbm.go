package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
}

func ReadPBM(filename string) (*PBM, error) {
	var dimension string
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Lecture de la première ligne pour obtenir le magic number
	scanner.Scan()
	line := scanner.Text()
	line = strings.TrimSpace(line)
	if line != "P1" && line != "P4" {
		return nil, fmt.Errorf("Not a Portable Bitmap file: bad magic number %s", line)
	}
	magicNumber := line

	// Lecture des dimensions
	for scanner.Scan() {
		if scanner.Text()[0] == '#' {
			continue
		}
		break

	}

	dimension = scanner.Text()
	res := strings.Split(dimension, " ")
	height, _ := strconv.Atoi(res[1])
	width, _ := strconv.Atoi(res[0])

	// Lecture des données binaires
	data := make([][]bool, height)
	for i := range data {
		data[i] = make([]bool, width)
	}

	for i := 0; i < height; i++ {
		scanner.Scan()
		line := scanner.Text()
		hori := strings.Fields(line)
		for j := 0; j < width; j++ {
			verti, _ := strconv.Atoi(hori[j])
			if verti == 1 {
				data[i][j] = true
			}
		}
	}

	pbm := &PBM{
		data:        data,
		width:       width,
		height:      height,
		magicNumber: magicNumber,
	}

	fmt.Printf("%+v\n", PBM{data, width, height, magicNumber})

	return pbm, nil
}

func (pbm *PBM) Size() (int, int) {
	return pbm.width, pbm.height
}

func (pbm *PBM) At(x, y int) bool {
	return pbm.data[y][x]
}

func (pbm *PBM) Set(x, y int, value bool) {
	pbm.data[y][x] = value
}

func (pbm *PBM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "magicNumber: %s\n", pbm.magicNumber)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(file, "Width: %d\n", pbm.width)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(file, "Height: %d\n", pbm.height)
	if err != nil {
		return err
	}
	return nil
}

func (pbm *PBM) Invert() {
	for i := 0; i < pbm.height; i++ {
		for j := 0; j < pbm.width; j++ {
			if pbm.data[i][j] {
				pbm.data[i][j] = false
			} else if !pbm.data[i][j] {
				pbm.data[i][j] = true
			}
		}
	}
}

func (pbm *PBM) Flip() {
	for _, height := range pbm.data {
		for i, j := 0, len(height)-1; i < j; i, j = i+1, j-1 {
			height[i], height[j] = height[j], height[i]
		}
	}
}

func (pbm *PBM) Flop() {
	for i, j := 0, len(pbm.data)-1; i < j; i, j = i+1, j-1 {
		pbm.data[i], pbm.data[j] = pbm.data[j], pbm.data[i]
	}
}

func (pbm *PBM) SetMagicNumber(magicNumber string) {
	fmt.Println(magicNumber)
}

func main() {
	pbm, _ := ReadPBM("Dos_pbm/test.pbm")
	// (*PBM).Size(&PBM{})
	pbm.Save("Dos_pbm/save.pbm")
	fmt.Println("\n")

	pbm.SetMagicNumber("P4")
	pbm.Flip()
	fmt.Println("Flip:", pbm.data)
	fmt.Println("\n")

	pbm.Flop()
	fmt.Println("Flop:", pbm.data)
}
