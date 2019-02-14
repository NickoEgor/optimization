package iohandler

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"gonum.org/v1/gonum/mat"
)

// MatPrint - prints beautiful matrix/vector
func MatPrint(X mat.Matrix) {
	fa := mat.Formatted(X, mat.Prefix(""), mat.Squeeze())
	// fmt.Printf("%.3f\n", fa)
	fmt.Printf("%v\n", fa)
}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func checkScanner(s *bufio.Scanner) {
	if err := s.Err(); err != nil {
		panic(err)
	}
}

// EnterConditions - reads matrix from file
func EnterConditions(filename string) (*mat.Dense, *mat.VecDense, *mat.VecDense, *mat.VecDense, *mat.VecDense) {
	f, err := os.Open(filename)
	checkError(err)
	defer f.Close()

	scanner := bufio.NewScanner(f)

	// ignore 'A' line
	scanner.Scan()
	checkScanner(scanner)

	scanner.Scan()
	checkScanner(scanner)
	sizeStr := strings.Split(scanner.Text(), " ")
	rows, err := strconv.Atoi(sizeStr[0])
	checkError(err)
	cols, err := strconv.Atoi(sizeStr[1])
	checkError(err)

	A := mat.NewDense(rows, cols, nil)

	for i := 0; i < rows && scanner.Scan(); i++ {
		line := scanner.Text()
		values := strings.Split(line, " ")
		for j := 0; j < cols; j++ {
			v, err := strconv.ParseFloat(values[j], 64)
			checkError(err)
			A.Set(i, j, v)
		}
	}
	checkScanner(scanner)

	// ignore 'b' line
	scanner.Scan()
	checkScanner(scanner)

	scanner.Scan()
	checkScanner(scanner)
	b := mat.NewVecDense(rows, nil)
	line := scanner.Text()
	values := strings.Split(line, " ")
	for i := 0; i < rows; i++ {
		v, err := strconv.ParseFloat(values[i], 64)
		checkError(err)
		b.SetVec(i, v)
	}

	// ignore 'c' line
	scanner.Scan()
	checkScanner(scanner)

	scanner.Scan()
	checkScanner(scanner)
	c := mat.NewVecDense(cols, nil)
	line = scanner.Text()
	values = strings.Split(line, " ")
	for i := 0; i < cols; i++ {
		v, err := strconv.ParseFloat(values[i], 64)
		checkError(err)
		c.SetVec(i, v)
	}

	// ignore 'x' line
	scanner.Scan()
	checkScanner(scanner)

	scanner.Scan()
	checkScanner(scanner)
	x := mat.NewVecDense(cols, nil)
	line = scanner.Text()
	values = strings.Split(line, " ")
	for i := 0; i < cols; i++ {
		v, err := strconv.ParseFloat(values[i], 64)
		checkError(err)
		x.SetVec(i, v)
	}

	// ignore 'J' line
	scanner.Scan()
	checkScanner(scanner)

	scanner.Scan()
	checkScanner(scanner)
	J := mat.NewVecDense(rows, nil)
	line = scanner.Text()
	values = strings.Split(line, " ")
	for i := 0; i < rows; i++ {
		v, err := strconv.Atoi(values[i])
		checkError(err)
		J.SetVec(i, float64(v))
	}

	return A, b, c, x, J
}
