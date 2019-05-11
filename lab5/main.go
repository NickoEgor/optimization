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
		filename = "conditions6.txt"
	} else {
		filename = os.Args[1]
	}

	Cost, Need, Stock := io.EnterTransportConditions(filename)

	flow, err := cm.TransportPotentials(Cost, Need, Stock)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Flow matrix:")
	io.MatPrint(flow)

	rows, cols := flow.Dims()

	fmt.Println("Consumption:")
	consumption := mat.NewVecDense(rows, nil)
	for i := 0; i < rows; i++ {
		consumption.SetVec(i, mat.Sum(flow.RowView(i)))
	}
	io.MatPrint(consumption.T())

	fmt.Println("Satisfied needs:")
	satisfaction := mat.NewVecDense(cols, nil)
	for i := 0; i < cols; i++ {
		satisfaction.SetVec(i, mat.Sum(flow.ColView(i)))
	}
	io.MatPrint(satisfaction.T())

	fmt.Println("Minimum cost:")
	tempSum := mat.NewDense(rows, cols, nil)
	tempSum.MulElem(Cost, flow)
	cost := mat.Sum(tempSum)
	fmt.Println(cost)
}
