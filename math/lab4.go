package math

import (
	"errors"
	"math"
	// "time"

	"gonum.org/v1/gonum/mat"
	// io "labs/moiu"
)

// DualSimplex - dual simplex method
func DualSimplex(A *mat.Dense, b, c, y, J *mat.VecDense) (*mat.VecDense, *mat.VecDense, error) {
	rows, cols := A.Dims()

	Ab := mat.NewDense(rows, rows, nil)
	cB := mat.NewVecDense(rows, nil)
	for i := 0; i < rows; i++ {
		j := int(J.AtVec(i))
		cB.SetVec(i, float64(c.AtVec(j)))
		Ab.SetCol(i, mat.Col(nil, j, A))
	}

	Abi := mat.NewDense(rows, rows, nil)
	Abi.Inverse(Ab)

	something := mat.NewVecDense(rows, nil)
	something.MulVec(Abi.T(), cB)

	delta := mat.NewVecDense(cols, nil)
	delta.MulVec(A.T(), something)
	delta.SubVec(delta, c)

	cnt := 0
	for {
		JNon := mat.NewVecDense(cols-rows, nil)
		k := 0
		for i := 0; i < cols; i++ {
			isBasis := false
			for j := 0; j < rows; j++ {
				if int(J.AtVec(j)) == i {
					isBasis = true
					break
				}
			}
			if !isBasis {
				JNon.SetVec(k, float64(i))
				k++
			}
		}

		An := mat.NewDense(rows, cols-rows, nil)
		for i := 0; i < cols-rows; i++ {
			j := int(JNon.AtVec(i))
			An.SetCol(i, mat.Col(nil, j, A))
		}

		deltaN := mat.NewVecDense(cols-rows, nil)
		for i := 0; i < cols-rows; i++ {
			j := int(JNon.AtVec(i))
			deltaN.SetVec(i, delta.AtVec(j))
		}

		nu := mat.NewVecDense(rows, nil)
		nu.MulVec(Abi, b)

		isAllGreaterThenZero := true
		for i := 0; i < rows; i++ {
			if nu.AtVec(i) < 0 {
				isAllGreaterThenZero = false
			}
		}

		if isAllGreaterThenZero {
			result := mat.NewVecDense(cols, nil)
			for i := 0; i < rows; i++ {
				j := int(J.AtVec(i))
				result.SetVec(j, nu.AtVec(i))
			}
			return result, J, nil
		}

		minVal := mat.Min(nu)
		minInd := -1
		for i := 0; i < rows; i++ {
			if nu.AtVec(i) == minVal {
				minInd = i
				break
			}
		}

		jk := int(J.AtVec(minInd))

		mu := mat.NewVecDense(cols-rows, nil)
		mu.MulVec(An.T(), Abi.RowView(minInd))

		sigma := mat.NewVecDense(cols-rows, nil)
		for i := 0; i < cols-rows; i++ {
			if mu.AtVec(i) < 0 {
				sigma.SetVec(i, -deltaN.AtVec(i)/mu.AtVec(i))
			} else {
				sigma.SetVec(i, math.Inf(1))
			}
		}

		sigma0 := mat.Min(sigma)

		if sigma0 == math.Inf(1) {
			return nil, nil, errors.New("No solutions")
		}

		minSigmaInd := -1
		for i := 0; i < cols; i++ {
			if sigma.AtVec(i) == sigma0 {
				minSigmaInd = i
				break
			}
		}

		j0 := int(JNon.AtVec(minSigmaInd))
		J.SetVec(minInd, float64(j0))

		for i := 0; i < cols-rows; i++ {
			j := int(JNon.AtVec(i))
			delta.SetVec(j, delta.AtVec(j)+sigma0*mu.AtVec(i))
		}

		delta.SetVec(jk, sigma0)

		newCol := mat.NewVecDense(rows, mat.Col(nil, j0, A))
		newAbi, err := ImprovedInverse(Ab, Abi, newCol, minInd)
		if err != nil {
			return nil, nil, err
		}
		Ab.SetCol(minInd, mat.Col(nil, j0, A))
		Abi = newAbi

		cnt++
		// time.Sleep(2 * time.Second)
	}
}
