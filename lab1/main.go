package main

import (
	"bufio"
	"fmt"
	"gonum.org/v1/gonum/mat"
	// "math/rand"
	"os"
	"strconv"
	"strings"
	// "time"
)

func matPrint(X mat.Matrix) {
	fa := mat.Formatted(X, mat.Prefix(""), mat.Squeeze())
	// fmt.Printf("%.3f\n", fa)
	fmt.Printf("%v\n", fa)
}

func getKnownMatrices() (*mat.Dense, *mat.Dense) {
	// s := rand.NewSource(time.Now().Unix())
	// r := rand.New(s)
	// length := 3
	// vl := make([]float64, length*length)
	// m1 := mat.NewDense(length, length, vl)
	// m2 := mat.NewDense(length, length, nil)
	// is_acceptable := false

	// for !is_acceptable {
	// for i := 0; i < length*length; i++ {
	// vl[i] = float64((r.Int() % 7) - 4)
	// }

	// if mat.Det(m1) == 1 {
	// is_acceptable = true
	// }
	// }

	// m2.Inverse(m1)
	// matPrint(m1)
	// matPrint(m2)

	// return m1, m2

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

func eye(n int) *mat.Dense {
	m := mat.NewDense(n, n, nil)
	for i := 0; i < n; i++ {
		m.Set(i, i, 1)
	}
	return m
}

func checkInverse(m1, m2 mat.Matrix) bool {
	r1, c1 := m1.Dims()
	r2, c2 := m2.Dims()
	if r1 != c1 || r2 != c2 || r1 != r2 {
		return false
	}

	D := mat.NewDense(3, 3, nil)
	D.Product(m1, m2)
	// matPrint(D)
	r, c := D.Dims()
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			row := D.RawRowView(i)
			if (i == j && row[j] != 1) || (i != j && row[j] != 0) {
				return false
			}
		}
	}
	return true
}

func main() {
	A, A_i := getKnownMatrices()
	println("Source matrix:")
	matPrint(A)
	println("Inversed matrix:")
	matPrint(A_i)

	res := checkInverse(A, A_i)
	if !res {
		println("Matrices are incorrect")
		os.Exit(1)
	}

	r, _ := A.Dims()

	vec := getReplacingVector()
	println("Replacing vector:")
	matPrint(vec)

	if vec.Len() != r {
		println("Incorrect vector size")
		os.Exit(1)
	}

	print("Enter column index: ")
	reader := bufio.NewReader(os.Stdin)
	index_str, err := reader.ReadString('\n')
	index_str = strings.TrimSuffix(index_str, "\n")
	index, err := strconv.Atoi(index_str)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	A.SetCol(index, vec.RawVector().Data)
	println("Inversion needed for matrix:")
	matPrint(A)

	l := mat.NewVecDense(r, nil)
	l.MulVec(A_i, vec)

	li := l.AtVec(index)
	if li == 0 {
		println("l[i] equals to 0")
		os.Exit(0)
	}

	l.SetVec(index, -1)
	l.ScaleVec(-1.0/li, l)
	Q := eye(r)
	Q.SetCol(index, l.RawVector().Data)
	Ans := mat.NewDense(r, r, nil)
	Ans.Product(Q, A_i)
	println("New inversed matrix:")
	matPrint(Ans)

	println("Product:")
	A.Product(A, Ans)
	matPrint(A)
}
