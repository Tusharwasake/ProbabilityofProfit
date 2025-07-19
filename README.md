# Probability of Profit Calculator

A Go-based REST API that calculates the Probability of Profit (PoP) for options trading strategies using Monte Carlo simulation and implied volatility data from Firstock API.

## Overview

This API calculates the probability that a given options strategy will end up profitable by expiry using statistical simulation methods.

### Features

- REST API endpoint for PoP calculations
- Multi-leg options strategy support
- Black-Scholes implied volatility calculations
- Monte Carlo price simulation (500,000 iterations)
- JSON input/output format
- Unit tests included

## Quick Start

### Prerequisites

- Go 1.24.4 or higher
- Git

### Installation

```bash
git clone <repository-url>
cd ProbabilityofProfit
go mod download
```

### Configuration

Create a `.env` file for Firstock API credentials (optional - used for authentication only):

```env
FIRSTOCK_USER_ID=your_user_id
FIRSTOCK_PASSWORD=your_password
FIRSTOCK_TOTP_SECRET=your_totp_secret
FIRSTOCK_API_KEY=your_api_key
FIRSTOCK_VENDOR_CODE=your_vendor_code
SERVER_PORT=8080
```

Note: The application uses Black-Scholes calculations for implied volatility. Firstock credentials are only required for authentication.

### Running the Server

```bash
go run main.go
```

The server will start on `http://localhost:8080`

## API Documentation

### Health Check

```bash
GET /status
```

### Calculate Probability of Profit

```bash
POST /pop
```

**Request Format:**

```json
{
  "spot": 22913.15,
  "expiry": "06-MAR-2025",
  "daysToExpiry": 8,
  "symbol": "NIFTY",
  "optionList": [
    {
      "optionType": "CE",
      "transactionType": "B",
      "strike": 22950,
      "ltp": 154.7,
      "quantity": 75
    }
  ]
}
```

**Response Format:**

```json
{
  "pop": 31.03
}
```

## Implementation Details

### Monte Carlo Simulation

The application uses Monte Carlo simulation with 500,000 iterations to calculate probability of profit. Price movements are modeled using **Geometric Brownian Motion** with log-normal distribution:

```
S(T) = S(0) × exp((r - 0.5 × σ²) × T + σ × √T × Z)
```

Where:

- S(T) = Simulated price at expiry
- S(0) = Current spot price
- r = Risk-free rate (6.5% annual)
- σ = Implied volatility (calculated per option)
- T = Time to expiry (years)
- Z = Standard normal random variable

**Note**: The drift term (r - 0.5 × σ²) is simplified to (-0.5 × σ²) assuming risk-neutral pricing.

### Black-Scholes Model for IV Calculation

The system calculates implied volatility using the **Black-Scholes-Merton model**:

**Call Option Price:**

```
C = S × N(d₁) - K × e^(-r×T) × N(d₂)
```

**Put Option Price:**

```
P = K × e^(-r×T) × N(-d₂) - S × N(-d₁)
```

Where:

```
d₁ = [ln(S/K) + (r + 0.5×σ²)×T] / (σ×√T)
d₂ = d₁ - σ×√T
```

**Variables:**

- C/P = Call/Put option price
- S = Current stock price
- K = Strike price
- r = Risk-free rate (6.5%)
- T = Time to expiry (years)
- σ = Volatility (what we solve for)
- N() = Cumulative standard normal distribution

### Newton-Raphson IV Calculation

Implied volatility is calculated using **Newton-Raphson iterative method**:

```
σₙ₊₁ = σₙ - (BS_Price(σₙ) - Market_Price) / Vega(σₙ)
```

**Parameters:**

- **Initial guess**: σ₀ = 0.2 (20%)
- **Maximum iterations**: 100
- **Convergence tolerance**: 1e-8
- **Volatility bounds**: 1e-6 ≤ σ ≤ 5.0
- **Vega calculation**: ∂C/∂σ = S × √T × φ(d₁)

### P&L Calculation Model

For each Monte Carlo simulation, the system calculates option payoffs and net P&L:

**Option Payoffs:**

- **Call Options**: max(S(T) - K, 0)
- **Put Options**: max(K - S(T), 0)

**Net P&L Calculation:**

- **Long positions**: (Payoff - Premium) × Quantity
- **Short positions**: (Premium - Payoff) × Quantity

**Final PoP**: Count of profitable simulations / Total simulations × 100

### Volatility Aggregation

When multiple options are present, the system:

1. Calculates individual IV for each option using its LTP
2. Computes weighted average IV across all valid options
3. Uses this average IV for all price simulations
4. Validates IV results (positive, within bounds)

## Project Structure

```
ProbabilityofProfit/
├── main.go                 # Application entry point
├── controller/
│   └── pop_controller.go   # HTTP request handlers
├── service/
│   └── pop_service.go      # Business logic and calculations
├── model/
│   └── pop_model.go        # Data structures
├── firstock/
│   └── client.go           # Firstock authentication & Black-Scholes calculations
├── test/
│   └── pop_test.go         # Unit tests
├── go.mod                  # Go module dependencies
└── .env                    # Configuration file
```

## Testing

Run the test suite:

```bash
go test ./test -v
```

The test suite includes 5 comprehensive test cases:

```
=== RUN   TestLongCallStrategy
    pop_test.go:36: Long Call Strategy - PoP: 31.03%
--- PASS: TestLongCallStrategy (0.03s)
=== RUN   TestShortPutSpreadStrategy
    pop_test.go:63: Short Put Spread Strategy - PoP: 55.70%
--- PASS: TestShortPutSpreadStrategy (0.04s)
=== RUN   TestATMCallStrategy
    pop_test.go:83: ATM Call Strategy - PoP: 33.19%
--- PASS: TestATMCallStrategy (0.02s)
=== RUN   TestCoveredCallStrategy
    pop_test.go:103: Short Call Strategy - PoP: 71.70%
--- PASS: TestCoveredCallStrategy (0.03s)
=== RUN   TestInvalidInputs
--- PASS: TestInvalidInputs (0.00s)
PASS
ok      pop-calculator/test     1.381s
```

### Test Coverage

- **Long Call Strategy**: Tests bullish call buying with OTM strike
- **Short Put Spread**: Tests bullish put spread with both strikes below spot
- **ATM Call Strategy**: Tests at-the-money call option behavior
- **Short Call Strategy**: Tests covered call writing with OTM strike
- **Invalid Inputs**: Tests error handling with zero LTP values

## Mathematical Models & Assumptions

### Black-Scholes Model Assumptions

1. **Constant Risk-Free Rate**: 6.5% annual (fixed in code)
2. **Constant Volatility**: Implied volatility remains constant until expiry
3. **Log-Normal Price Distribution**: Underlying follows Geometric Brownian Motion
4. **European Exercise**: Options exercised only at expiration
5. **No Dividends**: Dividend yield = 0% for all calculations
6. **Continuous Trading**: No gaps or trading halts
7. **No Transaction Costs**: Commissions and bid-ask spreads ignored
8. **Perfect Liquidity**: Can trade any quantity at market prices

### Numerical Method Parameters

**Newton-Raphson IV Solver:**

- **Convergence Tolerance**: 1e-8 (0.00000001)
- **Maximum Iterations**: 100
- **Initial Volatility Guess**: 20%
- **Volatility Bounds**: 0.0001% to 500%
- **Numerical Stability Check**: Vega > 1e-6

**Monte Carlo Simulation:**

- **Number of Simulations**: 500,000 iterations
- **Random Seed**: Fixed seed (12345) for reproducible results
- **Distribution**: Standard normal (Box-Muller transformation)
- **Precision**: Float64 (double precision)

### Data Input Assumptions

1. **LTP Validation**: Last Traded Price must be > 0
2. **Time Conversion**: Days to expiry converted to years (T = days/365)
3. **Option Types**: "CE" for calls, "PE" for puts
4. **Transaction Types**: "B" for buy, "S" for sell
5. **Strike Prices**: Must be positive values
6. **Quantities**: Integer values representing lot sizes

### Error Handling & Edge Cases

- **Invalid IV Inputs**: Returns error for S≤0, K≤0, T≤0, or MarketPrice≤0
- **Convergence Failure**: Returns "failed to converge" after 100 iterations
- **Numerical Instability**: Returns error when Vega < 1e-6
- **No Valid Options**: Returns 0% PoP when all options have invalid LTP
- **Volatility Bounds**: Clips calculated IV to [1e-6, 5.0] range

### Model Limitations

1. **Static IV**: Doesn't account for volatility smile/skew
2. **No Greeks Hedging**: Doesn't consider delta hedging or gamma risk
3. **Single Underlying**: Multi-asset correlations not modeled
4. **American Options**: Early exercise features not supported
5. **Interest Rate Risk**: Fixed risk-free rate assumption
6. **Liquidity Risk**: Perfect liquidity assumption may not hold

## Dependencies

- **gin-gonic/gin**: HTTP web framework
- **joho/godotenv**: Environment variable management
- **pquerna/otp**: TOTP authentication for Firstock API
- **stretchr/testify**: Testing framework for assertions

## Performance

- **Simulation Speed**: 500,000 Monte Carlo iterations in ~0.03 seconds
- **IV Calculation**: Newton-Raphson convergence typically within 5-10 iterations
- **Memory Efficient**: Minimal memory footprint with streaming calculations
- **Numerical Precision**: Float64 precision for all mathematical operations
- **Deterministic Results**: Fixed random seed ensures reproducible outputs

## Example Usage

```bash
# Start the server
go run main.go

# Test health endpoint
curl http://localhost:8080/status

# Calculate PoP for a long call strategy
curl -X POST http://localhost:8080/pop \
  -H "Content-Type: application/json" \
  -d '{
    "spot": 25000,
    "daysToExpiry": 10,
    "symbol": "NIFTY",
    "optionList": [{
      "optionType": "CE",
      "transactionType": "B",
      "strike": 25100,
      "ltp": 150,
      "quantity": 1
    }]
  }'
```
