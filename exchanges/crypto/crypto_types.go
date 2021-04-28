package crypto

type wsRequest struct {
	ID     int                    `json:"id"`
	Method string                 `json:"method"`
	Params map[string]interface{} `json:"params"`
	Nonce  int64                  `json:"nonce"`
	Sig    string                 `json:"sig"`
}

type GenericResponse struct {
	Code    responseCode `json:"code"`
	Method  string       `json:"method"`
	ID      int          `json:"id"`
	Message string       `json:"message"`
}
type InstrumentResponse struct {
	Result InstrumentResult `json:"result"`
}
type InstrumentResult struct {
	Instruments []Instruments `json:"instruments"`
}
type Instruments struct {
	InstrumentName       string `json:"instrument_name"`
	QuoteCurrency        string `json:"quote_currency"`
	BaseCurrency         string `json:"base_currency"`
	PriceDecimals        int    `json:"price_decimals"`
	QuantityDecimals     int    `json:"quantity_decimals"`
	MarginTradingEnabled bool   `json:"margin_trading_enabled"`
}
