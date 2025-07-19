package service

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"pop-calculator/firstock"
	"pop-calculator/model"
)

func CalculatePoPValue(spot float64, daysToExpiry float64, expiryDate string, symbol string, optionList []model.OptionLeg) float64 {
	const numSimulations = 500000
	T := daysToExpiry / 365

	profitableCount := 0
	r := rand.New(rand.NewSource(12345))

	optionIVs := getOptionIVs(optionList, spot, T)

	totalIV := 0.0
	validIVCount := 0
	for _, iv := range optionIVs {
		if iv > 0 {
			totalIV += iv
			validIVCount++
		}
	}

	var avgIV float64
	if validIVCount > 0 {
		avgIV = totalIV / float64(validIVCount)
		log.Printf("Using average IV: %.4f (%.2f%%) from %d options", avgIV, avgIV*100, validIVCount)
	} else {
		// If no IV could be calculated, we cannot proceed with accurate simulation
		log.Printf("Error: Cannot calculate PoP - no valid IV data available")
		log.Printf("Ensure all options have valid LTP (Last Traded Price) values")
		return 0.0
	}

	for i := 0; i < numSimulations; i++ {
		z := r.NormFloat64()
		simulatedPrice := spot * math.Exp((z*avgIV*math.Sqrt(T))-(0.5*avgIV*avgIV*T))

		pnl := calculatePnLWithIV(simulatedPrice, optionList)
		if pnl >= 0 {
			profitableCount++
		}
	}
	result := float64(profitableCount) / float64(numSimulations) * 100
	return math.Round(result*100) / 100
}

// getOptionIVs calculates IV for each option using Black-Scholes

func getOptionIVs(optionList []model.OptionLeg, spot float64, timeToExpiry float64) map[string]float64 {
	optionIVs := make(map[string]float64)
	riskFreeRate := 0.065 // consider 6.5%

	log.Printf("Calculating IV for %d options", len(optionList))

	for _, leg := range optionList {
		strikeKey := fmt.Sprintf("%.0f_%s", leg.Strike, leg.OptionType)

		if leg.LTP <= 0 {
			log.Printf("Skipping %s: invalid LTP %.2f", strikeKey, leg.LTP)
			continue
		}

		isCall := leg.OptionType == "CE"

		calculatedIV, err := firstock.CalculateImpliedVolatility(
			spot,
			leg.Strike,
			timeToExpiry,
			riskFreeRate,
			leg.LTP,
			isCall,
		)

		if err != nil {
			log.Printf("Failed to calculate IV for %s (LTP %.2f): %v", strikeKey, leg.LTP, err)
			continue
		}

		if calculatedIV <= 0 {
			log.Printf("Invalid IV result for %s: %.4f from LTP %.2f", strikeKey, calculatedIV, leg.LTP)
			continue
		}

		optionIVs[strikeKey] = calculatedIV
		log.Printf("Calculated IV for %s: %.4f (%.2f%%) from LTP %.2f", strikeKey, calculatedIV, calculatedIV*100, leg.LTP)
	}

	if len(optionIVs) == 0 {
		log.Printf("Warning: No valid IV calculations - all options lack proper LTP data")
	} else {
		log.Printf("Successfully calculated IV for %d out of %d options", len(optionIVs), len(optionList))
	}

	return optionIVs
}

func calculatePnLWithIV(price float64, optionList []model.OptionLeg) float64 {
	totalPnL := 0.0

	for _, leg := range optionList {
		var payoff float64

		switch leg.OptionType {
		case "CE":
			payoff = math.Max(price-leg.Strike, 0)
		case "PE":
			payoff = math.Max(leg.Strike-price, 0)
		default:
			continue
		}

		adjustedLTP := leg.LTP

		switch leg.TransactionType {
		case "B":
			totalPnL += (payoff - adjustedLTP) * float64(leg.Quantity)
		case "S":
			totalPnL += (adjustedLTP - payoff) * float64(leg.Quantity)
		}
	}

	return totalPnL
}

