package math

import (
	"errors"
	"math"
)

// normalCDF calculates cumulative standard normal distribution
func normalCDF(x float64) float64 {
	return 0.5 * (1.0 + math.Erf(x/math.Sqrt(2)))
}

// normalPDF calculates standard normal probability density function
func normalPDF(x float64) float64 {
	return math.Exp(-0.5*x*x) / math.Sqrt(2*math.Pi)
}

// blackScholesPrice calculates theoretical option price
func blackScholesPrice(S, K, T, r, sigma float64, isCall bool) float64 {
	if T <= 0 {
		if isCall {
			return math.Max(S-K, 0)
		}
		return math.Max(K-S, 0)
	}
	
	d1 := (math.Log(S/K) + (r+0.5*sigma*sigma)*T) / (sigma * math.Sqrt(T))
	d2 := d1 - sigma*math.Sqrt(T)
	
	if isCall {
		return S*normalCDF(d1) - K*math.Exp(-r*T)*normalCDF(d2)
	}
	return K*math.Exp(-r*T)*normalCDF(-d2) - S*normalCDF(-d1)
}

// vega calculates option's sensitivity to volatility
func vega(S, K, T, r, sigma float64) float64 {
	if T <= 0 {
		return 0
	}
	
	d1 := (math.Log(S/K) + (r+0.5*sigma*sigma)*T) / (sigma * math.Sqrt(T))
	return S * math.Sqrt(T) * normalPDF(d1)
}

// CalculateImpliedVolatility calculates Black-Scholes implied volatility using Newton-Raphson
func CalculateImpliedVolatility(S, K, T, r, marketPrice float64, isCall bool) (float64, error) {
	const (
		maxIter   = 100
		tolerance = 1e-8
		minVol    = 1e-6
		maxVol    = 5.0
	)
	
	// Validate inputs
	if S <= 0 || K <= 0 || T <= 0 || marketPrice <= 0 {
		return 0, errors.New("invalid input parameters")
	}
	
	// Initial guess
	sigma := 0.2
	
	// Newton-Raphson iteration
	for i := 0; i < maxIter; i++ {
		price := blackScholesPrice(S, K, T, r, sigma, isCall)
		vegaVal := vega(S, K, T, r, sigma)
		
		if vegaVal < minVol {
			return 0, errors.New("numerical instability")
		}
		
		priceDiff := price - marketPrice
		
		if math.Abs(priceDiff) < tolerance {
			return sigma, nil
		}
		
		sigma = sigma - priceDiff/vegaVal
		
		// Keep sigma in bounds
		if sigma < minVol {
			sigma = minVol
		} else if sigma > maxVol {
			sigma = maxVol
		}
	}
	
	return 0, errors.New("failed to converge")
}
