package main

import (
	"flag"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var tax = 0.017

func main() {
	var maxMoneyWithBroker float64
	var money float64
	var interestPerMonth float64
	var months int

	flag.Float64Var(&interestPerMonth, "interest", 0.20, "Interest per month")
	flag.IntVar(&months, "duration", 5*12, "Simulation's duration in months")
	flag.Float64Var(&money, "initial-deposit", 20_000.0, "Initial investment to begin with")
	flag.Float64Var(&maxMoneyWithBroker, "max-money-with-broker", 4_000_000.0, "Maximum amount of money to keep in the brokerage account")
	flag.Parse()

	p := message.NewPrinter(language.English)

	bank := 0.0
	withdrawalPrinted := false

	for i := 0; i < months; i++ {
		interest := money * interestPerMonth
		money += interest

		// Can we start taking money?
		if money > maxMoneyWithBroker {
			moneyToTake := money - maxMoneyWithBroker
			netMoneyToTake := subtractTax(moneyToTake)
			if !withdrawalPrinted {
				p.Printf("You can take at least %.0f EUR per month after %d months (~ %d years)\n", netMoneyToTake, i, int(i/12))
				withdrawalPrinted = true
			}

			bank += netMoneyToTake
			money -= moneyToTake
		}
	}

	p.Printf("\n")
	p.Printf("Balance after %d years:\n", int(months/12))
	p.Printf("    Bank   = %.0f   (yields %.0f EUR passive income p.m)\n", round(bank), calcPassiveIncomePerMonth(bank))
	p.Printf("    Broker = %.0f\n", round(money))
	p.Printf("\n")
	p.Printf("Total possible passive income: %.0f EUR\n", calcPassiveIncomePerMonth(bank+subtractTax(money)))
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
