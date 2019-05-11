package main

import (
	"fmt"
	"os"

	"gonum.org/v1/gonum/mat"

	io "labs/moiu"
	cm "labs/moiu/math"
)

func main() {
	var filename string
	if len(os.Args) < 2 {
		filename = "conditions.txt"
	} else {
		filename = os.Args[1]
	}

	A, b, c := io.EnterCommonConditions(filename)
	xF, JF, err := cm.SimplexFirstPhase(&A, &b)
	if err != nil {
		fmt.Println(err)
		return
	}

	x, J := &xF, &JF
	println("A:")
	io.MatPrint(A)
	println("b:")
	io.MatPrint(b)
	println("c:")
	io.MatPrint(c)
	println("x:")
	io.MatPrint(x)
	println("J:")
	io.MatPrint(J)

	x, J, err = cm.SimplexMainPhase(A, b, c, x, J)
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
