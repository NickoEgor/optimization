package math

import (
	"errors"
	"fmt"
	"math"

	"gonum.org/v1/gonum/mat"
	// io "labs/moiu"
)

const epsilon = 1e-10

// SimplexFirstPhase - first phase of simplex method
func SimplexFirstPhase(Aptr **mat.Dense, bptr **mat.VecDense) (mat.VecDense, mat.VecDense, error) {
	A := *Aptr
	b := *bptr

	for i := 0; i < b.Len(); i++ {
		if b.AtVec(i) < 0 {
			b.SetVec(i, -1*b.AtVec(i))
			r := A.RawRowView(i)
			for j := 0; j < len(r); j++ {
				r[j] *= -1
			}
		}
	}

	rows, cols := A.Dims()
	xF := mat.NewVecDense(cols+rows,
		append(make([]float64, cols), b.RawVector().Data...))

	indeces := make([]float64, rows)
	for i := range indeces {
		indeces[i] = float64(cols + i)
	}

	JF := mat.NewVecDense(rows, indeces)

	cF := mat.NewVecDense(cols+rows, nil)
	for i := 0; i < cols; i++ {
		cF.SetVec(i, 0)
	}
	for i := cols; i < cols+rows; i++ {
		cF.SetVec(i, -1)
	}

	AF := mat.NewDense(rows, cols+rows, nil)
	for i := 0; i < cols; i++ {
		AF.SetCol(i, mat.Col(nil, i, A))
	}

	E := Eye(rows)
	for i := cols; i < cols+rows; i++ {
		AF.SetCol(i, mat.Col(nil, i-cols, E))
	}

	// println("AF:")
	// io.MatPrint(AF)
	// println("b:")
	// io.MatPrint(b)
	// println("cF:")
	// io.MatPrint(cF)
	// println("xF:")
	// io.MatPrint(xF)
	// println("JF:")
	// io.MatPrint(JF)

	xM, JM, err := SimplexMainPhase(AF, b, cF, xF, JF)
	if err != nil {
		return *xF, *JF, err
	}

	// println("xF:")
	// io.MatPrint(xM)
	// println("JF:")
	// io.MatPrint(JM)

	xF = xM
	JF = JM

	for i := cols; i < cols+rows; i++ {
		if math.Abs(xF.AtVec(i)) > epsilon {
			println()
			return *xF, *JF, fmt.Errorf("%f is not 0. Initial task is incompatible", xF.AtVec(i))
		}
	}

	for {
		rows, cols = A.Dims()
		E := Eye(rows)

		artificial := -1

		for i := 0; i < rows; i++ {
			if int(JF.AtVec(i)) >= cols {
				artificial = i
				break
			}
		}

		if artificial == -1 {
			return *mat.VecDenseCopyOf(xF.SliceVec(0, cols)), *JF, nil
		}

		Ainv := mat.NewDense(rows, rows, nil)
		for i := 0; i < rows; i++ {
			if int(JF.AtVec(i)) < cols {
				Ainv.SetCol(i, mat.Col(nil, i, A))
			} else {
				Ainv.SetCol(i, mat.Col(nil, int(JF.AtVec(i))-cols, E))
			}
		}
		Ainv.Inverse(Ainv)

		_, cm := AF.Dims()
		diff := cm - cols

		hasFound := false
		hasChanged := false

		// unbasis unartificial
		for i := 0; i < cols; i++ {
			isAppropriate := false
			for j := 0; j < JF.Len(); j++ {
				if i == int(JF.AtVec(j)) {
					isAppropriate = true
					break
				}
			}

			if isAppropriate {
				hasFound = true

				ek := mat.NewDense(1, diff, nil)
				ek.Set(0, i, 1)
				ek.Product(ek, Ainv)

				Aj := A.ColView(i)

				alphaM := mat.NewDense(1, 1, nil)
				alphaM.Product(ek, Aj)
				alpha := alphaM.At(0, 0)

				if math.Abs(alpha) > epsilon { // alpha isn't zero
					JF.SetVec(artificial, float64(i))
					hasChanged = true
					break
				} else {
					continue
				}
			} else {
				continue
			}
		}

		if !hasFound {
			return *xF, *JF, errors.New("Why am I here?")
		}

		// all alpha got zero
		if !hasChanged {
			ind := int(JF.AtVec(artificial)) - cols

			newA := mat.NewDense(rows-1, cols, nil)
			newAF := mat.NewDense(rows-1, cols+rows-1, nil)
			for r := 0; r < ind; r++ {
				newA.SetRow(r, mat.Row(nil, r, A))

				frow := mat.Row(nil, r, AF)
				if ind != cols+rows-1 {
					frow = append(frow[:cols+ind], frow[cols+ind+1:]...)
				} else {
					frow = frow[:cols+ind]
				}
				newAF.SetRow(r, frow)
			}
			for r := ind + 1; r < rows; r++ {
				newA.SetRow(r-1, mat.Row(nil, r, A))
				frow := mat.Row(nil, r, AF)
				if ind != cols+rows-1 {
					frow = append(frow[:cols+ind], frow[cols+ind+1:]...)
				} else {
					frow = frow[:cols+ind]
				}
				newAF.SetRow(r-1, frow)
			}
			*Aptr = newA
			A = *Aptr
			AF = newAF

			newJF := JF.RawVector().Data
			if ind != JF.Len()-1 {
				newJF = append(newJF[:ind], newJF[ind+1:]...)
			} else {
				newJF = newJF[:ind]
			}
			for j := 0; j < JF.Len(); j++ {
				if int(JF.AtVec(j)) >= ind {
					JF.SetVec(j, JF.AtVec(j)-1)
				}
			}
			JF = mat.NewVecDense(JF.Len()-1, newJF)

			newxF := xF.RawVector().Data
			if ind != xF.Len()-1 {
				newxF = append(newxF[:cols+ind], newxF[cols+ind+1:]...)
			} else {
				newxF = newxF[:cols+ind]
			}
			xF = mat.NewVecDense(xF.Len()-1, newxF)

			newb := b.RawVector().Data
			if ind != b.Len()-1 {
				newb = append(newb[:ind], newb[ind+1:]...)
			} else {
				newb = newb[:ind]
			}
			*bptr = mat.NewVecDense(b.Len()-1, newb)
			b = *bptr

			continue
		}
	}
}
