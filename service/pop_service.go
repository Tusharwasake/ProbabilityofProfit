package service

import (
	"fmt"
	"math"
	"math/rand"
	"pop-calculator/firstock"
	"pop-calculator/model"
	"time"
)

// func CalculatePoPValue(spot float64, daysToExpiry float64, expiryDate string, symbol string, optionList []model.OptionLeg) float64 {

// 	const numSimulations = 10000

// 	T := daysToExpiry / 365
// 	mean := spot

// 	// Get real-time IV for each option leg using Firstock

// 	optionIVs := make(map[string]float64)

// 	for _, leg := range optionList {

// 		optionKey := fmt.Sprintf("%.0f_%s", leg.Strike, leg.OptionType)

//         if _, exists := optionIVs[optionKey];
// 		!exists {
//             optionSymbol := formatOptionSymbol(symbol, expiryDate, leg.Strike, leg.OptionType)
//             exchange := getExchangeForSymbol(symbol)
//             optionIVs[optionKey] = firstock.GetIV(exchange, optionSymbol)
//         }
//     }

// 	// Use average IV for market simulation
// 	totalIV := 0.0
// 	for _, iv := range optionIVs {
// 		totalIV += iv
// 	}

// 	avgIV := totalIV / float64(len(optionIVs))

// 	if avgIV == 0 {
// 		avgIV = 0.2 // fallback
// 	}

// 	stdDev := spot * avgIV * math.Sqrt(T)

// 	r := rand.New(rand.NewSource(time.Now().UnixNano()))

// 	profitableCount := 0

// 	for i:= 0; i < numSimulations; i++ {

// 		//  This assumes price can go negative
// 		// simulatedPrice := r.NormFloat64()*stdDev + mean

// 	// Stock prices should be log-normally distributed

// 	simulatedPrice := spot * math.Exp((r.NormFloat64()*avgIV*math.Sqrt(T)) - (0.5*avgIV*avgIV*T))

// 	pnl := calculatePnLWithIV(simulatedPrice, optionList, optionIVs, expiryDate, symbol)
// 		if pnl >= 0 {
// 			profitableCount++
// 		}
// 	}

// 	pop := float64(profitableCount) / float64(numSimulations)

// 	return math.Round(pop*100) / 100
// }

func CalculatePoPValue(spot float64, daysToExpiry float64, expiryDate string, symbol string, optionList []model.OptionLeg) float64 {
    const numSimulations = 10000
    T := daysToExpiry / 365

    // mean := spot

    // Get real-time IV for each unique option
    optionIVs := make(map[string]float64)
    
    for _, leg := range optionList {
        optionKey := fmt.Sprintf("%.0f_%s", leg.Strike, leg.OptionType)
        if _, exists := optionIVs[optionKey]; !exists {
            optionSymbol := formatOptionSymbol(symbol, expiryDate, leg.Strike, leg.OptionType)
            exchange := getExchangeForSymbol(symbol)
            optionIVs[optionKey] = firstock.GetIV(exchange, optionSymbol)
        }
    }
    
    // Calculate average IV with safety check
    totalIV := 0.0
    for _, iv := range optionIVs {
        totalIV += iv
    }
    
    var avgIV float64
    if len(optionIVs) == 0 {
        avgIV = 0.2 // fallback
    } else {
        avgIV = totalIV / float64(len(optionIVs))
        if avgIV == 0 {
            avgIV = 0.2 // fallback
        }
    }
    
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    profitableCount := 0
    
    for i := 0; i < numSimulations; i++ {

		// This assumes price can go negative
		// simulatedPrice := r.NormFloat64()*stdDev + mean

        // Correct log-normal price simulation
        randomFactor := r.NormFloat64()
        simulatedPrice := spot * math.Exp((randomFactor*avgIV*math.Sqrt(T)) - (0.5*avgIV*avgIV*T))
        
        pnl := calculatePnLWithIV(simulatedPrice, optionList, optionIVs, expiryDate, symbol)
        if pnl >= 0 {
            profitableCount++
        }

    }
    
    pop := float64(profitableCount) / float64(numSimulations)
    return math.Round(pop*100) / 100
}

func calculatePnLWithIV(price float64, optionList []model.OptionLeg, optionIVs map[string]float64, expiryDate string, symbol string) float64 {
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

		switch leg.TransactionType {
		case "B":
			totalPnL += (payoff - leg.LTP) * float64(leg.Quantity)
		case "S":
			totalPnL += (leg.LTP - payoff) * float64(leg.Quantity)
		}
	}

	return totalPnL
}

// formatOptionSymbol formats an option symbol for Firstock API
func formatOptionSymbol(symbol, expiryDate string, strike float64, optionType string) string {
	strikeStr := fmt.Sprintf("%.0f", strike)
	return fmt.Sprintf("%s%s%s%s", symbol, expiryDate, strikeStr, optionType)
}

// getExchangeForSymbol returns the appropriate exchange for a given symbol
func getExchangeForSymbol(symbol string) string {
	switch symbol {
	case "NIFTY", "BANKNIFTY", "FINNIFTY":
		return "NFO"
	case "SENSEX", "BANKEX":
		return "BFO"
	default:
		return "NFO"
	}
}
