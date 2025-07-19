package firstock

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/pquerna/otp/totp"
)

// Global variables for authentication
var JKey string
var UserID string

const BaseURL = "https://connect.thefirstock.com/api/V4"

type LoginRequest struct {
	UserID     string `json:"userId"`
	Password   string `json:"password"`
	TOTP       string `json:"TOTP"`
	ApiKey     string `json:"apiKey"`
	VendorCode string `json:"vendorCode"`
}


type LoginResponse struct {
	Status string `json:"status"`
	Data   struct {
		ActID      string `json:"actid"`
		UserName   string `json:"userName"`
		SUserToken string `json:"susertoken"`
		Email      string `json:"email"`
	} `json:"data"`
}

func hashPassword(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}


func generateTOTP(secret string) (string, error) {
	return totp.GenerateCode(secret, time.Now())
}

func Login(userID, password, totp, apiKey, vendorCode string) (string, error) {
	
	// Hash the password as required by Firstock API
	hashedPassword := hashPassword(password)
	
	loginReq := LoginRequest{
		UserID:     userID,
		Password:   hashedPassword, 
		TOTP:       totp,
		ApiKey:     apiKey,
		VendorCode: vendorCode,
	}
	
	jsonData, err := json.Marshal(loginReq)
	if err != nil {
		return "", fmt.Errorf("failed to marshal login request in json []byte: %w", err)
	}

	req, err := http.NewRequest("POST", BaseURL+"/login", bytes.NewBuffer(jsonData)) 
	if err != nil {
		return "", fmt.Errorf("failed to create login request: %w", err)
	}

	// Set proper headers
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make login request: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var loginResp LoginResponse

	if err := json.Unmarshal(body, &loginResp); 
	
	err != nil {
		return "", 
		fmt.Errorf("failed to unmarshal login response: %w", err)
	}

	if loginResp.Status != "success" {
		return "", 
		fmt.Errorf("login failed: %s", string(body))
	}

	// Set global variables
	JKey = loginResp.Data.SUserToken
	UserID = loginResp.Data.ActID

	return JKey, nil
}

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

// InitializeFromEnv initializes Firstock client from environment variables

func InitializeFromEnv() error {
	userID := os.Getenv("FIRSTOCK_USER_ID")
	password := os.Getenv("FIRSTOCK_PASSWORD")
	totpSecret := os.Getenv("FIRSTOCK_TOTP_SECRET")
	apiKey := os.Getenv("FIRSTOCK_API_KEY")
	vendorCode := os.Getenv("FIRSTOCK_VENDOR_CODE")

	if userID == "" || password == "" || totpSecret == "" || apiKey == "" || vendorCode == "" {
		return fmt.Errorf("missing required environment variables")
	}

	totpCode, err := generateTOTP(totpSecret)
	if err != nil {
		return fmt.Errorf("failed to generate TOTP: %w", err)
	}

	jKey, err := Login(userID, password, totpCode, apiKey, vendorCode)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	JKey = jKey
	UserID = userID
	return nil
}
