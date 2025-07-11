package model

type OptionLeg struct {
	OptionType      string  `json:"optionType"`
	TransactionType string  `json:"transactionType"`
	Strike          float64 `json:"strike"`
	LTP             float64 `json:"ltp"`
	Quantity        int     `json:"quantity"`
}

type PopRequest struct {
	Spot         float64     `json:"spot"`
	Expiry       string      `json:"expiry"`
	DaysToExpiry float64     `json:"daysToExpiry"`
	Symbol       string      `json:"symbol"`
	OptionList   []OptionLeg `json:"optionList"`
}

type PopResponse struct {
	Pop float64 `json:"pop"`
}
