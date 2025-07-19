package test

import (
	"testing"

	"pop-calculator/model"
	"pop-calculator/service"

	"github.com/stretchr/testify/assert"
)

const (
	TestSpot         = 22913.15
	TestDaysToExpiry = 8.0
	TestExpiry       = "06-MAR-2025"
	TestSymbol       = "NIFTY"
)

func TestLongCallStrategy(t *testing.T) {
	optionList := []model.OptionLeg{
		{
			OptionType:      "CE",
			TransactionType: "B",
			Strike:          22950,
			LTP:             154.7,
			Quantity:        75,
		},
	}

	pop := service.CalculatePoPValue(TestSpot, TestDaysToExpiry, TestExpiry, TestSymbol, optionList)

	assert.Greater(t, pop, 0.0, "PoP should be greater than 0")
	assert.LessOrEqual(t, pop, 100.0, "PoP should be less than or equal to 100")
	assert.Greater(t, pop, 20.0, "Long call PoP should be reasonable given strike vs spot")
	
	t.Logf("Long Call Strategy - PoP: %.2f%%", pop)
}

func TestShortPutSpreadStrategy(t *testing.T) {
	optionList := []model.OptionLeg{
		{
			OptionType:      "PE",
			TransactionType: "S",
			Strike:          22900,
			LTP:             145.5,
			Quantity:        75,
		},
		{
			OptionType:      "PE",
			TransactionType: "B",
			Strike:          22850,
			LTP:             98.2,
			Quantity:        75,
		},
	}

	pop := service.CalculatePoPValue(TestSpot, TestDaysToExpiry, TestExpiry, TestSymbol, optionList)

	assert.Greater(t, pop, 0.0, "PoP should be greater than 0")
	assert.LessOrEqual(t, pop, 100.0, "PoP should be less than or equal to 100")
	assert.Greater(t, pop, 40.0, "Short put spread PoP should be high when both strikes are below spot")
	
	t.Logf("Short Put Spread Strategy - PoP: %.2f%%", pop)
}

func TestATMCallStrategy(t *testing.T) {
	optionList := []model.OptionLeg{
		{
			OptionType:      "CE",
			TransactionType: "B",
			Strike:          22900, // Close to ATM
			LTP:             180.0,
			Quantity:        50,
		},
	}

	pop := service.CalculatePoPValue(TestSpot, TestDaysToExpiry, TestExpiry, TestSymbol, optionList)

	assert.Greater(t, pop, 0.0)
	assert.LessOrEqual(t, pop, 100.0)
	assert.Greater(t, pop, 30.0, "ATM call should have reasonable PoP")
	
	t.Logf("ATM Call Strategy - PoP: %.2f%%", pop)
}

func TestCoveredCallStrategy(t *testing.T) {
	optionList := []model.OptionLeg{
		{
			OptionType:      "CE",
			TransactionType: "S", // Sell call
			Strike:          23000,
			LTP:             120.0,
			Quantity:        50,
		},
	}

	pop := service.CalculatePoPValue(TestSpot, TestDaysToExpiry, TestExpiry, TestSymbol, optionList)

	assert.Greater(t, pop, 0.0)
	assert.LessOrEqual(t, pop, 100.0)
	assert.Greater(t, pop, 60.0, "Short OTM call should have high PoP")
	
	t.Logf("Short Call Strategy - PoP: %.2f%%", pop)
}

func TestInvalidInputs(t *testing.T) {
	// Test with zero LTP
	optionList := []model.OptionLeg{
		{
			OptionType:      "CE",
			TransactionType: "B",
			Strike:          22950,
			LTP:             0.0, // Invalid LTP
			Quantity:        75,
		},
	}

	pop := service.CalculatePoPValue(TestSpot, TestDaysToExpiry, TestExpiry, TestSymbol, optionList)
	assert.Equal(t, 0.0, pop, "Should return 0 for invalid LTP")
}