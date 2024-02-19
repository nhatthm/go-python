package main

import (
	"fmt"

	"go.nhat.io/python3"
)

func main() { //nolint: govet
	math := python3.MustImportModule("math")

	pyResult := math.CallMethodArgs("sqrt", 4)
	defer pyResult.DecRef()

	result := python3.AsFloat64(pyResult)

	fmt.Printf("sqrt(4) = %.2f\n", result)

	// Output:
	// sqrt(4) = 2.00
}
