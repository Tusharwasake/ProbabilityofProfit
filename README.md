# 🧠 PoP Calculator (Go) — Options Strategy Profitability API

A high-performance **Go-based backend API** that calculates the **Probability of Profit (PoP)** for multi-leg options trading strategies using **Monte Carlo simulation** and **live implied volatility** from the **Firstock API**.

---

## ✅ Task: What Was Asked

**Assignment Goal:**

> Build an API that calculates the probability that a given options strategy will end up profitable by expiry.

### Requirements Given:

| Feature / Deliverable                  | Required? | Status  |
| -------------------------------------- | --------- | ------- |
| REST API POST `/pop`                   | ✅ Yes    | ✅ Done |
| Accepts options strategy as JSON input | ✅ Yes    | ✅ Done |
| Uses implied volatility (IV)           | ✅ Yes    | ✅ Done |
| Simulates 10,000 expiry prices         | ✅ Yes    | ✅ Done |
| Computes net P\&L for each simulation  | ✅ Yes    | ✅ Done |
| Returns `{"pop": float}`               | ✅ Yes    | ✅ Done |
| Unit tests for 2 strategies            | ✅ Yes    | ✅ Done |
| Documentation for assumptions & model  | ✅ Yes    | ✅ Done |

---

## 👨‍💻 What I Did (Highlights)

✅ Built a **POST API** `/pop` using **Gin** framework
✅ Parsed options data (`CE/PE`, `Buy/Sell`, `LTP`, `strike`, etc.)
✅ Connected with **Firstock API** to fetch **real-time IV** (or fallback)
✅ Simulated expiry prices using **log-normal distribution**
✅ Calculated **P\&L for each leg** of the strategy
✅ Counted simulations where overall payoff ≥ 0 to compute **PoP**
✅ Added **unit tests** for Long Call & Short Put Spread
✅ Wrote **detailed documentation**, health check, error handling, and `.env` support

---

## 🏁 How to Run (Step-by-Step for Beginners)

### 1. 🔧 Prerequisites

- Go 1.24.4+
- Git installed
- (Optional) Firstock account credentials for real IV

---

### 2. 📦 Setup Project

```bash
# Clone the repository
git clone <repo-url>
cd pop-calculator

# Install dependencies
go mod download
```

---

### 3. 🔐 Setup .env (Optional, for real-time IV)

```bash
cat > .env << EOF
FIRSTOCK_USER_ID=your_user_id
FIRSTOCK_PASSWORD=your_password
FIRSTOCK_TOTP_SECRET=your_totp_key
FIRSTOCK_API_KEY=your_api_key
FIRSTOCK_VENDOR_CODE=your_vendor_code
EOF
```

---

### 4. 🚀 Run the Server

```bash
# Build binary
go build -o pop-calculator

# Run binary
./pop-calculator

# OR run directly with go
go run main.go
```

---

### 5. ✅ Test the API

**Health check:**

```bash
curl http://localhost:8080/status
```

**Calculate PoP for a sample strategy:**

```bash
curl -X POST http://localhost:8080/pop \
  -H "Content-Type: application/json" \
  -d '{
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
  }'
```

**Response:**

```json
{
  "pop": 0.41
}
```

---

## 🧠 What It Does Internally

### 🔢 Core Simulation Logic

- **Log-normal simulation:**
  $S(T) = S(0) \times \exp\left(-\frac{\sigma^2}{2} \cdot T + \sigma \cdot \sqrt{T} \cdot Z\right)$

- **Monte Carlo engine:**
  Simulates 10,000 expiry prices using above formula

- **P\&L Calculation:**

  - Long Call: `(max(S - K, 0) - premium) * quantity`
  - Short Put: `(premium - max(K - S, 0)) * quantity`

- **PoP = Profitable Simulations / Total Simulations**

---

## 📁 Project Structure

```bash
pop-calculator/
├── main.go                 # App entrypoint
├── controller/
│   └── pop_controller.go   # Handles HTTP POST /pop
├── service/
│   └── pop_service.go      # Business logic (IV fetch, simulation)
├── model/
│   └── pop_model.go        # Data models (request/response)
├── firstock/
│   └── client.go           # Connects to Firstock for IV
├── test/
│   └── pop_test.go         # Unit tests for strategies
├── go.mod                  # Go dependencies
└── .env                    # Your credentials (optional)
```

---

## 🧪 Unit Tests

You can run all tests using:

```bash
go test ./...
```

Example test output:

```
=== RUN   TestLongCallStrategy
    pop_test.go:42: Long Call PoP: 0.41
--- PASS: TestLongCallStrategy (0.00s)

=== RUN   TestShortPutSpreadStrategy
    pop_test.go:72: Short Put Spread PoP: 0.53
--- PASS: TestShortPutSpreadStrategy (0.00s)
```

✅ Tests for:

- Long Call
- Short Put Spread
- Fallback IV calculation
- Net payoff logic

---

## 📡 API Details

### POST `/pop`

| Field               | Type      | Description                      |
| ------------------- | --------- | -------------------------------- |
| `spot`              | `float64` | Current price of underlying      |
| `expiry`            | `string`  | Expiry date (e.g., 06-MAR-2025)  |
| `daysToExpiry`      | `int`     | Days to expiry                   |
| `symbol`            | `string`  | Trading symbol (NIFTY/BANKNIFTY) |
| `optionList`        | `array`   | List of option legs              |
| ↳ `optionType`      | `string`  | "CE" or "PE"                     |
| ↳ `transactionType` | `string`  | "B" (Buy) or "S" (Sell)          |
| ↳ `strike`          | `float64` | Strike price of the option       |
| ↳ `ltp`             | `float64` | Premium (last traded price)      |
| ↳ `quantity`        | `int`     | Number of contracts              |

### Response:

```json
{
  "pop": 0.65
}
```

---

## 📚 Modeling Assumptions

| Assumption                    | Description                              |
| ----------------------------- | ---------------------------------------- |
| Risk-free rate = 0            | Simplified model                         |
| No dividends                  | Not factored into option prices          |
| Log-normal price distribution | More realistic than normal               |
| European-style options        | Exercised only at expiry                 |
| Transaction costs ignored     | Net payoff = intrinsic - premium         |
| Constant IV                   | Assumes IV remains unchanged till expiry |

---

## 🔧 Advanced Features

- Real-time IV via **Firstock API**
- Fallback IV: Base 15% + premium factor
- `.env` support for environment configs
- Built-in `/status` endpoint
- Custom strategies: Iron Condors, Straddles, Spreads supported

---

## 🛠 Technologies Used

| Tech         | Role                     |
| ------------ | ------------------------ |
| Go           | Backend language         |
| Gin          | HTTP API framework       |
| Firstock API | IV data source           |
| rand/math    | Price simulations        |
| TOTP Auth    | Secure login to Firstock |
| JSON         | Input/output format      |
| go test      | Unit testing             |

---

## 🚀 Ready to Use?

**Start it:**

```bash
go run main.go
```

**Ping it:**

```bash
curl http://localhost:8080/status
```

**Try a strategy:**

```bash
curl -X POST http://localhost:8080/pop -d @strategy.json
```
