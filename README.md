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

- Built a **POST API** `/pop` using **Gin** framework
- Parsed options data (`CE/PE`, `Buy/Sell`, `LTP`, `strike`, etc.)
- Connected with **Firstock API** to fetch **real-time IV** (or fallback)
- Simulated expiry prices using **log-normal distribution**
- Calculated **P\&L for each leg** of the strategy
- Counted simulations where overall payoff ≥ 0 to compute **PoP**
- Added **unit tests** for Long Call & Short Put Spread
- Wrote **detailed documentation**, health check, error handling, and `.env` support

---

## 🏁 How to Run (Step-by-Step for Beginners)

### 1. 🔧 Prerequisites

- Go 1.24.4+
- Git installed

### 2. 📦 Setup Project

```bash
git clone <repo-url>
cd pop-calculator
go mod download
```

### 3. 🔐 Setup .env (Optional, for real-time IV)

Create a `.env` file in the root directory:

```bash
# Copy from .env.example
cp .env.example .env
```

Or create manually:

```env
# Firstock API Configuration (Optional - for real-time IV data)
FIRSTOCK_USER_ID=your_user_id
FIRSTOCK_PASSWORD=your_password
FIRSTOCK_TOTP_SECRET=your_totp_secret_key
FIRSTOCK_API_KEY=your_api_key
FIRSTOCK_VENDOR_CODE=your_vendor_code

# Server Configuration
SERVER_PORT=8080
```

**Note:** Without these credentials, the app will use fallback IV calculations.

### 4. 🚀 Run the Server

```bash
go build -o pop-calculator
./pop-calculator
# OR
go run main.go
```

### 5. ✅ Test the API

**Health check:**

```bash
curl http://localhost:8080/status
```

**Calculate PoP:**

```bash
curl -X POST http://localhost:8080/pop \
  -H "Content-Type: application/json" \
  -d @strategy.json
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
├── main.go
├── controller/pop_controller.go
├── service/pop_service.go
├── model/pop_model.go
├── firstock/client.go
├── test/pop_test.go
├── go.mod
└── .env
```

---

## 🧪 Unit Tests

```bash
# Run all tests with verbose output
go test ./test -v

# Run specific test
go test ./test -run TestLongCallStrategy -v

# or
go test ./test -run TestShortPutSpreadStrategy -v



```
=== RUN   TestLongCallStrategy
    pop_test.go:42: Long Call PoP: 0.41
--- PASS: TestLongCallStrategy (0.00s)

=== RUN   TestShortPutSpreadStrategy
    pop_test.go:72: Short Put Spread PoP: 0.53
--- PASS: TestShortPutSpreadStrategy (0.00s)
```

---

## 📡 API Details

### POST `/pop`

**Request Body:**

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

**Response:**

```json
{
  "pop": 0.41
}
```

---

## 📚 Modeling Assumptions

| Assumption                    | Description                      |
| ----------------------------- | -------------------------------- |
| Risk-free rate = 0            | Simplified model                 |
| No dividends                  | Not factored into option prices  |
| Log-normal price distribution | More realistic than normal       |
| European-style options        | Exercised only at expiry         |
| Transaction costs ignored     | Net payoff = intrinsic - premium |
| Constant IV                   | Assumes IV remains till expiry   |

---

## 🛠️ Technologies Used

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

```bash
go run main.go
curl http://localhost:8080/status
curl -X POST http://localhost:8080/pop -d @strategy.json
```
