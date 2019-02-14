package main

import (
	"fmt"

	"gonum.org/v1/gonum/mat"

	io "labs/moiu"
	cm "labs/moiu/math"
)

func main() {
	A, b, c, x, J := io.EnterConditions("conditions.txt")
	x, J, err := cm.SimplexMainPhase(A, b, c, x, J)
	if err != nil {
		fmt.Println(err)
		return
	}
	println("Optimized basis (J):")
	io.MatPrint(J)
	println("Optimized plan (x):")
	io.MatPrint(x)
	println("Loss function:")
	loss := mat.NewVecDense(1, nil)
	loss.MulVec(c.T(), x)
	io.MatPrint(loss)
}
