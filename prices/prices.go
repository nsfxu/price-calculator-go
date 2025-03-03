package prices

import (
	"fmt"

	"example.com/price-calculator-go/conversion"
	"example.com/price-calculator-go/iomanager"
)

type TaxIncludedPriceJob struct {
	IOManager         iomanager.IOManager `json:"-"`
	TaxRate           float64             `json:"taxRate"`
	Prices            []float64           `json:"prices"`
	TaxIncludedPrices map[string]string   `json:"taxIncludedPrices"`
}

func (job *TaxIncludedPriceJob) LoadData() error {
	lines, err := job.IOManager.ReadLines()

	if err != nil {
		return err
	}

	prices, err := conversion.StringsToFloats(lines)

	if err != nil {
		return err
	}

	if len(prices) > 0 {
		job.Prices = prices
	}

	return nil
}

func (job TaxIncludedPriceJob) Process(doneChannel chan bool, errorChannel chan error) {
	err := job.LoadData()

	// errorChannel <- errors.New("An error!")

	if err != nil {
		errorChannel <- err
		return
	}

	result := make(map[string]string)

	for _, price := range job.Prices {
		taxIncludedPrice := price * (1 + job.TaxRate)
		result[fmt.Sprintf("%.2f", price)] = fmt.Sprintf("%.2f", taxIncludedPrice)
	}

	job.TaxIncludedPrices = result
	err = job.IOManager.WriteResult(job)

	if err != nil {
		errorChannel <- err
		return
	}

	doneChannel <- true
}

func NewTaxIncludedPriceJob(iom iomanager.IOManager, taxRate float64) *TaxIncludedPriceJob {
	return &TaxIncludedPriceJob{
		IOManager: iom,
		Prices:    []float64{10, 20, 30},
		TaxRate:   taxRate,
	}
}
