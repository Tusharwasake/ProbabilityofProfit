package test

import (
	"testing"

	"pop-calculator/model"
	"pop-calculator/service"

	"github.com/stretchr/testify/assert"
)

// Test data constants
const (
	TestSpot         = 22913.15
	TestDaysToExpiry = 8.0
	TestExpiry       = "06-MAR-2025"
	TestSymbol       = "NIFTY"
)

// TestLongCallStrategy tests a bullish long call strategy
func TestLongCallStrategy(t *testing.T) {
	
	// Long Call setup
	optionList := []model.OptionLeg{
		{
			OptionType:      "CE",
			TransactionType: "B",
			Strike:          22950,
			LTP:             154.7,
			Quantity:        75,
		},
	}

	//Calculate PoP
	pop := service.CalculatePoPValue(TestSpot, TestDaysToExpiry, TestExpiry, TestSymbol, optionList)

	//Validate results
	assert.Greater(t, pop, 0.0, "PoP should be greater than 0")
	assert.LessOrEqual(t, pop, 1.0, "PoP should be less than or equal to 1")
	
	assert.Greater(t, pop, 0.3, "Long call PoP should be reasonable given strike vs spot")
	
	t.Logf("Long Call Strategy - PoP: %.2f (%.1f%%)", pop, pop*100)
}

// TestShortPutSpreadStrategy tests a bullish short put spread strategy
func TestShortPutSpreadStrategy(t *testing.T) {
	optionList := []model.OptionLeg{
		{
			OptionType:      "PE",
			TransactionType: "S", // Sell higher strike put
			Strike:          22900,
			LTP:             145.5,
			Quantity:        75,
		},
		{
			OptionType:      "PE",
			TransactionType: "B", // Buy lower strike put
			Strike:          22850,
			LTP:             98.2,
			Quantity:        75,
		},
	}

	pop := service.CalculatePoPValue(TestSpot, TestDaysToExpiry, TestExpiry, TestSymbol, optionList)

	assert.Greater(t, pop, 0.0, "PoP should be greater than 0")
	assert.LessOrEqual(t, pop, 1.0, "PoP should be less than or equal to 1")
	
	// Short put spread should have high PoP (both strikes below spot)

	assert.Greater(t, pop, 0.5, "Short put spread PoP should be high when both strikes are below spot")
	
	t.Logf("Short Put Spread Strategy - PoP: %.2f (%.1f%%)", pop, pop*100)
}