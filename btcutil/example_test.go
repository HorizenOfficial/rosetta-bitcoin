package btcutil_test

import (
	"fmt"
	"math"

	"github.com/HorizenOfficial/rosetta-zen/btcutil"
)

func ExampleAmount() {

	a := btcutil.Amount(0)
	fmt.Println("Zero Zentoshi:", a)

	a = btcutil.Amount(1e8)
	fmt.Println("100,000,000 Zentoshi:", a)

	a = btcutil.Amount(1e5)
	fmt.Println("100,000 Zentoshi:", a)
	// Output:
	// Zero Zentoshi: 0 ZEN
	// 100,000,000 Zentoshi: 1 ZEN
	// 100,000 Zentoshi: 0.001 ZEN
}

func ExampleNewAmount() {
	amountOne, err := btcutil.NewAmount(1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountOne) //Output 1

	amountFraction, err := btcutil.NewAmount(0.01234567)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountFraction) //Output 2

	amountZero, err := btcutil.NewAmount(0)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountZero) //Output 3

	amountNaN, err := btcutil.NewAmount(math.NaN())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountNaN) //Output 4

	// Output: 1 ZEN
	// 0.01234567 ZEN
	// 0 ZEN
	// invalid zen amount
}

func ExampleAmount_unitConversions() {
	amount := btcutil.Amount(44433322211100)

	fmt.Println("Zentoshi to kZEN:", amount.Format(btcutil.AmountKiloZEN))
	fmt.Println("Zentoshi to ZEN:", amount)
	fmt.Println("Zentoshi to MilliZEN:", amount.Format(btcutil.AmountMilliZEN))
	fmt.Println("Zentoshi to MicroZEN:", amount.Format(btcutil.AmountMicroZEN))
	fmt.Println("Zentoshi to Zentoshi:", amount.Format(btcutil.AmountZentoshi))

	// Output:
	// Zentoshi to kZEN: 444.333222111 kZEN
	// Zentoshi to ZEN: 444333.222111 ZEN
	// Zentoshi to MilliZEN: 444333222.111 mZEN
	// Zentoshi to MicroZEN: 444333222111 μZEN
	// Zentoshi to Zentoshi: 44433322211100 Zentoshi
}
