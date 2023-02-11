package main

import (
	"os"

	"github.com/jessevdk/go-flags"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var tax = 0.017

type Options struct {
	InterestPerMonth   float64 `long:"interest" default:"0.20" description:"Interest per month"`
	Months             int     `long:"duration" default:"60" description:"Simulation's duration in months"`
	Money              float64 `long:"initial-deposit" default:"20000" description:"Initial investment to begin with"`
	MaxMoneyWithBroker float64 `long:"max-money-with-broker" default:"4000000" description:"Maximum amount of money to keep in the brokerage account"`
}

func main() {
	var ops Options

	if _, err := flags.Parse(&ops); err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
			os.Exit(1)
		default:
			os.Exit(1)
		}
	}

	p := message.NewPrinter(language.English)

	bank := 0.0
	withdrawalPrinted := false

	for i := 0; i < ops.Months; i++ {
		interest := ops.Money * ops.InterestPerMonth
		ops.Money += interest

		// Can we start taking money?
		if ops.Money > ops.MaxMoneyWithBroker {
			moneyToTake := ops.Money - ops.MaxMoneyWithBroker
			netMoneyToTake := subtractTax(moneyToTake)
			if !withdrawalPrinted {
				p.Printf("You can take at least %.0f EUR per month after %d months (~ %d years)\n", netMoneyToTake, i, int(i/12))
				withdrawalPrinted = true
			}

			bank += netMoneyToTake
			ops.Money -= moneyToTake
		}
	}

	p.Printf("\n")
	p.Printf("Balance after %d years:\n", int(ops.Months/12))
	p.Printf("    Bank   = %.0f   (yields %.0f EUR passive income p.m)\n", round(bank), calcPassiveIncomePerMonth(bank))
	p.Printf("    Broker = %.0f\n", round(ops.Money))
	p.Printf("\n")
	p.Printf("Total possible passive income: %.0f EUR\n", calcPassiveIncomePerMonth(bank+subtractTax(ops.Money)))
}

func round(m float64) float64 {
	return float64(100.0 * int(m/100))
}

func calcPassiveIncomePerMonth(money float64) float64 {
	passiveAnnualInterest := 0.04
	return round(money * passiveAnnualInterest / 12)
}

func subtractTax(money float64) float64 {
	return (1 - tax) * money
}
