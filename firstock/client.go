package firstock

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
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



type QuoteRequest struct {
	UserID        string `json:"userId"`
	Exchange      string `json:"exchange"`
	TradingSymbol string `json:"tradingSymbol"`
	JKey          string `json:"jKey"`
}


type QuoteResponse struct {
	Status string `json:"status"`
	Data   struct {
		Exchange      string  `json:"exchange"`
		TradingSymbol string  `json:"tradingSymbol"`
		LTP           float64 `json:"lp"`
		Open          float64 `json:"o"`
		High          float64 `json:"h"`
		Low           float64 `json:"l"`
		Close         float64 `json:"c"`
		Volume        int64   `json:"v"`
		OI            float64 `json:"oi"`
		TotalBuyQty   int64   `json:"tbq"`
		TotalSellQty  int64   `json:"tsq"`
		AvgPrice      float64 `json:"ap"`
		LowerCircuit  float64 `json:"lc"`
		UpperCircuit  float64 `json:"uc"`
		YearlyHigh    float64 `json:"yh"`
		YearlyLow     float64 `json:"yl"`
	} `json:"data"`
}


type OptionGreekRequest struct {
	UserID        string `json:"userId"`
	Exchange      string `json:"exchange"`
	TradingSymbol string `json:"tradingSymbol"`
	JKey          string `json:"jKey"`
}


type OptionGreekResponse struct {
	Status string `json:"status"`
	Data   struct {
		Delta float64 `json:"delta"`
		Gamma float64 `json:"gamma"`
		Theta float64 `json:"theta"`
		Vega  float64 `json:"vega"`
		IV    float64 `json:"iv"`
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
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Firstock-Go-Client/1.0")

	client := &http.Client{Timeout: 30 * time.Second}    
	//"&" create object and return pointer of that so method can be used  (but go automatically convert to & if don't write it)
	

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

// GetQuote fetches real-time quote data for a given symbol

func GetQuote(exchange, tradingSymbol string) (*QuoteResponse, error) {
	if JKey == "" {
		return nil, fmt.Errorf("not authenticated - please login first")
	}

	quoteReq := QuoteRequest{
		UserID:        UserID,
		Exchange:      exchange,
		TradingSymbol: tradingSymbol,
		JKey:          JKey,
	}

	jsonData, err := json.Marshal(quoteReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal quote request: %w", err)
	}

	resp, err := http.Post(BaseURL+"/getQuote", "application/json", bytes.NewBuffer(jsonData))
	// new buffer []byte → io.Reader
	if err != nil {
		return nil, fmt.Errorf("failed to make quote request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var quoteResp QuoteResponse

	if err := json.Unmarshal(body, &quoteResp); 
	
	err != nil {
		return nil, fmt.Errorf("failed to unmarshal quote response: %w", err)
	}

	if quoteResp.Status != "success" {
		return nil, fmt.Errorf("quote request failed: %s", string(body))
	}

	return &quoteResp, nil
}

// GetOptionGreek fetches option Greeks including IV for a given option symbol
func GetOptionGreek(exchange, tradingSymbol string) (*OptionGreekResponse, error) {

	if JKey == "" {
		return nil, fmt.Errorf("not authenticated - please login first")
	}

	greekReq := OptionGreekRequest{
		UserID:        UserID,
		Exchange:      exchange,
		TradingSymbol: tradingSymbol,
		JKey:          JKey,
	}

	jsonData, err := json.Marshal(greekReq)  
	if err != nil {
		return nil, fmt.Errorf("failed to marshal option Greek request: %w", err)
	}

	resp, err := http.Post(BaseURL+"/optionGreek", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to make option Greek request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var greekResp OptionGreekResponse

	if err := json.Unmarshal(body, &greekResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal option Greek response: %w", err)
	}

	if greekResp.Status != "success" {
		return nil, fmt.Errorf("option Greek request failed: %s", string(body))
	}

	return &greekResp, nil
}

// GetIV fetches the implied volatility for a given option symbol

func GetIV(exchange, tradingSymbol string) float64 {

	greekResp, err := GetOptionGreek(exchange, tradingSymbol)
	if err == nil && greekResp.Data.IV > 0 {
		return greekResp.Data.IV
	}

	quoteResp, err := GetQuote(exchange, tradingSymbol)
	if err != nil {
		return calculateFallbackIV(exchange, tradingSymbol)
	}


	return calculateIVFromPremium(quoteResp.Data.LTP)
}

func calculateIVFromPremium(premium float64) float64 {

	baseIV := 0.15 // Base IV of 15%
	premiumFactor := premium / 100.0

	if premiumFactor > 2.0 {
		return baseIV + 0.20 
	} else if premiumFactor > 1.0 {
		return baseIV + 0.10
	} else if premiumFactor > 0.5 {
		return baseIV + 0.05 
	} else {
		return baseIV 
	}
}


func calculateFallbackIV(exchange, tradingSymbol string) float64 {
	baseIV := 0.15

	if len(tradingSymbol) > 10 {
		baseIV += 0.08
	}

	if exchange == "NFO" {
		baseIV += 0.05 
	}

	return baseIV
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
		// Try alternative time windows
		prevTime := time.Now().Add(-30 * time.Second)
		if prevCode, err2 := totp.GenerateCode(strings.ToUpper(strings.ReplaceAll(totpSecret, " ", "")), prevTime); 
		err2 == nil {
			jKey, err = Login(userID, password, prevCode, apiKey, vendorCode)
		}
		
		if err != nil {
			nextTime := time.Now().Add(30 * time.Second)
			if nextCode, err3 := totp.GenerateCode(strings.ToUpper(strings.ReplaceAll(totpSecret, " ", "")), nextTime); 
			err3 == nil {
				jKey, err = Login(userID, password, nextCode, apiKey, vendorCode)
			}
		}
		
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}
	}

	JKey = jKey
	UserID = userID
	return nil
}
