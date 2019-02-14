package math

import (
	"errors"
	"math"

	"gonum.org/v1/gonum/mat"
)

//Eye - create identity matrix
func Eye(n int) *mat.Dense {
	m := mat.NewDense(n, n, nil)
	for i := 0; i < n; i++ {
		m.Set(i, i, 1)
	}
	return m
}

// CheckInverse - checks if two matrix's product is identity matrix
func CheckInverse(m1, m2 mat.Matrix) bool {
	r1, c1 := m1.Dims()
	r2, c2 := m2.Dims()
	if r1 != c1 || r2 != c2 || r1 != r2 {
		return false
	}

	D := mat.NewDense(r1, r1, nil)
	D.Product(m1, m2)
	// println("Product")
	// io.MatPrint(D)

	eps := 1e-10
	r, c := D.Dims()
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			row := D.RawRowView(i)
			if (i == j && row[j] < 1-eps && row[i] > 1+eps) || (i != j && math.Abs(row[j]) > eps) {
				return false
			}
		}
	}
	return true
}

// ImprovedInverse - inverse matrix using optimised algorithm
func ImprovedInverse(A, Ai mat.Matrix, vec *mat.VecDense, index int) (*mat.Dense, error) {
	res := CheckInverse(A, Ai)
	if !res {
		return nil, errors.New("Matrices are incorrect")
	}

	r, _ := A.Dims()
	if vec.Len() != r {
		return nil, errors.New("Incorrect vector size")
	}

	l := mat.NewVecDense(r, nil)
	l.MulVec(Ai, vec)
	li := l.AtVec(index)
	if li == 0 {
		return nil, errors.New("l[i] equals to 0")
	}

	l.SetVec(index, -1)
	l.ScaleVec(-1.0/li, l)

	Q := Eye(r)
	Q.SetCol(index, l.RawVector().Data)
	ans := mat.NewDense(r, r, nil)
	ans.Product(Q, Ai)
	return ans, nil
}
