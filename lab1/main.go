package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"gonum.org/v1/gonum/mat"
	c "labs/moiu/math"
)

func getKnownMatrices() (*mat.Dense, *mat.Dense) {
	v1 := []float64{
		7, 2, 1,
		0, 3, -1,
		-3, 4, -2,
	}
	v2 := []float64{
		-2, 8, -5,
		3, -11, 7,
		9, -34, 21,
	}
	return mat.NewDense(3, 3, v1), mat.NewDense(3, 3, v2)
}

func getReplacingVector() *mat.VecDense {
	return mat.NewVecDense(3, []float64{1, 2, 3})
}

func main() {
	A, Ai := getKnownMatrices()
	println("Source matrix:")
	c.MatPrint(A)
	println("Inversed matrix:")
	c.MatPrint(Ai)

	vec := getReplacingVector()
	println("Replacing vector:")
	c.MatPrint(vec)

	print("Enter column index: ")
	reader := bufio.NewReader(os.Stdin)
	indexStr, err := reader.ReadString('\n')
	indexStr = strings.TrimSuffix(indexStr, "\n")
	index, err := strconv.Atoi(indexStr)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ans, err := c.ImprovedInverse(A, Ai, vec, index)
	if err != nil {
		println(err)
		os.Exit(1)
	}
	c.MatPrint(ans)
}
