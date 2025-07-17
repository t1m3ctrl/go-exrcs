package main

import (
	"fmt"
	"math/big"
)

func main() {
	a := big.NewInt(0)
	b := big.NewInt(0)

	a.Exp(big.NewInt(2), big.NewInt(50), nil)
	b.Exp(big.NewInt(2), big.NewInt(49), nil)

	fmt.Println("a =", a)
	fmt.Println("b =", b)

	sum := new(big.Int).Add(a, b)
	fmt.Println("a + b =", sum)

	diff := new(big.Int).Sub(a, b)
	fmt.Println("a - b =", diff)

	product := new(big.Int).Mul(a, b)
	fmt.Println("a * b =", product)

	quotient := new(big.Int).Div(a, b)
	fmt.Println("a / b =", quotient)
}
