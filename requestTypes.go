package main

type LoginUser struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password" `
}
type DematAccount struct {
	Code   string `json:"code" form:"code"`
	Broker string `json:"broker" form:"broker"`
}

type BankAccount struct {
	AccountNumber string `json:"account_number" form:"account_number"`
	Bank          string `json:"bank" form:"bank"`
}

type StockTrade struct {
	Symbol    string `json:"symbol" form:"symbol"`
	Demat     string `json:"demat" form:"demat"`
	Quantity  string `json:"quantity" form:"quantity"`
	Price     string `json:"price" form:"price"`
	TradeType string `json:"trade_type" form:"trade_type"`
	TradeDate string `json:"trade_date" form:"trade_date"`
}
type MFTrade struct {
	MutualFundId string `json:"mutual_fund_id" form:"mutual_fund_id"`
	Demat        string `json:"demat" form:"demat"`
	Quantity     string `json:"quantity" form:"quantity"`
	Price        string `json:"price" form:"price"`
	TradeType    string `json:"trade_type" form:"trade_type"`
	TradeDate    string `json:"trade_date" form:"trade_date"`
}
type ListedNCDTrade struct {
	Symbol    string `json:"symbol" form:"symbol"`
	Demat     string `json:"demat" form:"demat"`
	Quantity  string `json:"quantity" form:"quantity"`
	Price     string `json:"price" form:"price"`
	TradeType string `json:"trade_type" form:"trade_type"`
	TradeDate string `json:"trade_date" form:"trade_date"`
}
type UnlistedNCDTrade struct {
	Symbol    string `json:"symbol" form:"symbol"`
	Quantity  string `json:"quantity" form:"quantity"`
	Price     string `json:"price" form:"price"`
	TradeDate string `json:"trade_date" form:"trade_date"`
}

type FD struct {
	Bank        string `json:"bank" form:"bank"`
	Amount      string `json:"amount" form:"amount"`
	MtAmount    string `json:"maturity_amount" form:"maturity_amount"`
	IpRate      string `json:"ip_rate" form:"ip_rate"`
	IpFrequency string `json:"ip_frequency" form:"ip_frequency"`
	IpDate      string `json:"ip_date" form:"ip_date"`
	StartDate   string `json:"start_date" form:"start_date"`
	MtDate      string `json:"maturity_date" form:"maturity_date"`
	AccNumber   string `json:"account_number" form:"account_number"`
}
