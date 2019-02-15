package math

import (
	"errors"
	"math"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
)

func isBasisIndex(m *mat.VecDense, value int) (bool, int) {
	length := m.Len()
	for i := 0; i < length; i++ {
		if int(m.AtVec(i)) == value {
			return true, i
		}
	}
	return false, -1
}

// SimplexMainPhase - main phase of simplex method
func SimplexMainPhase(A *mat.Dense, b, c, x, J *mat.VecDense) (*mat.VecDense, *mat.VecDense, error) {
	rows, cols := A.Dims()

	for i := 0; i < rows; i++ {
		J.SetVec(i, J.AtVec(i)-1)
	}

	isFirstIteration := true
	var theta0idx, j0 int
	Ab := mat.NewDense(rows, rows, nil)
	Abi := mat.NewDense(rows, rows, nil)
	cb := mat.NewVecDense(rows, nil)

	for {
		if isFirstIteration {
			for i := 0; i < rows; i++ {
				j := int(J.AtVec(i))
				Ab.SetCol(i, mat.Col(nil, j, A))
				cb.SetVec(i, c.AtVec(j))
			}

			Abi.Inverse(Ab)
			isFirstIteration = false
		} else {
			replaced := mat.NewVecDense(rows, mat.Col(nil, j0, A))
			abi, err := ImprovedInverse(Ab, Abi, replaced, theta0idx)
			if err != nil {
				return nil, nil, err
			}

			Abi = abi
			Ab.SetCol(theta0idx, replaced.RawVector().Data)
			cb.SetVec(theta0idx, c.AtVec(j0))
		}

		u := mat.NewVecDense(rows, nil)
		u.MulVec(Abi.T(), cb)

		delta := mat.NewVecDense(cols, nil)
		delta.MulVec(A.T(), u)
		delta.SubVec(delta, c)

		deltas := make(map[int]int)
		isOptimized := true
		for i := 0; i < cols; i++ {
			res, _ := isBasisIndex(delta, i)
			if !res {
				value := delta.AtVec(i)
				if value < 0 {
					isOptimized = false
				}
				deltas[i] = int(delta.AtVec(i))
			}
		}

		if isOptimized {
			return x, J, nil
		}

		for k, v := range deltas {
			if v < 0 {
				j0 = k
				break
			}
		}

		z := mat.NewVecDense(rows, nil)
		z.MulVec(Abi, A.ColView(j0))

		theta := mat.NewVecDense(rows, nil)
		for i := 0; i < rows; i++ {
			if z.AtVec(i) > 0 {
				theta.SetVec(i, x.AtVec(int(J.AtVec(i)))/z.AtVec(i))
			} else {
				theta.SetVec(i, math.Inf(1))
			}
		}

		theta0 := mat.Min(theta)
		if theta0 == math.Inf(1) {
			return nil, nil, errors.New("Loss function is not limited from above")
		}

		theta0idx = int(floats.MinIdx(theta.RawVector().Data))
		J.SetVec(theta0idx, float64(j0))

		for i := 0; i < cols; i++ {
			if i == j0 {
				x.SetVec(i, theta0)
			} else {
				res, k := isBasisIndex(J, i)
				if res {
					x.SetVec(i, x.AtVec(i)-theta0*z.AtVec(k))
				} else {
					x.SetVec(i, 0)
				}
			}
		}
	}
}
