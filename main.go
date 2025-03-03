package main

import (
	"fmt"

	"example.com/price-calculator-go/file"
	"example.com/price-calculator-go/prices"
)

func main() {
	taxRates := []float64{0, 0.07, 0.1, 0.15}
	doneChannels := make([]chan bool, len(taxRates))
	errorChannels := make([]chan error, len(taxRates))

	for index, taxRate := range taxRates {
		doneChannels[index] = make(chan bool)
		errorChannels[index] = make(chan error)

		fm := file.New("prices.txt", fmt.Sprintf("result_%.0f.json", taxRate*100))
		// cmdm := cmdmanager.New()

		priceJob := prices.NewTaxIncludedPriceJob(fm, taxRate)
		go priceJob.Process(doneChannels[index], errorChannels[index])
	}

	for i := range taxRates {
		select {
		case err := <-errorChannels[i]:
			if err != nil {
				fmt.Println(err)
			}
		case <-doneChannels[i]:
			fmt.Println("Done!")
		}
	}
}
