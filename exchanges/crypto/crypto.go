package crypto

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/thrasher-corp/gocryptotrader/common/crypto"
	exchange "github.com/thrasher-corp/gocryptotrader/exchanges"
	"github.com/thrasher-corp/gocryptotrader/exchanges/request"
	"github.com/thrasher-corp/gocryptotrader/log"
	"net/http"
	"time"
)

// Crypto is the overarching type across this package
type Crypto struct {
	exchange.Base
}

type responseCode int

const (
	success               responseCode = iota
	sysError                           = 10001
	unauthorized                       = 10002
	ipIllegal                          = 10003
	badRequest                         = 10004
	userTierInvalid                    = 10005
	tooManyRequests                    = 10006
	invalidNonce                       = 10007
	methodNotFound                     = 10008
	invalidDateRange                   = 10009
	duplicateRecord                    = 20001
	negativeBalance                    = 20002
	symbolNotFound                     = 30003
	sideNotSupported                   = 30004
	orderTypeNotSupported              = 30004
)
const (
	cryptoAPIURL     = "https://api.crypto.com"
	cryptoAPIVersion = "v2"

	// Rate Limit constants

	// Public endpoints
	cryptoInstruments = "public/instruments"

	// Authenticated endpoints
	cryptoWithdrawalHistory = "private/get-withdrawal-history"
)

// Start implementing public and private exchange API funcs below

// GetInstruments Provides information on all supported instruments (e.g. BTC_USDT)
func (cr *Crypto) GetInstruments() (InstrumentResult, error) {
	var result InstrumentResponse
	err := cr.SendHTTPRequest(exchange.RestSpot, cryptoInstruments, nil, false, &result)
	return result.Result, err
}

func (cr *Crypto) SendHTTPRequest(ep exchange.URL, apiRequest string, params map[string]interface{}, authenticated bool, result interface{}) (err error) {
	if !cr.API.AuthenticatedSupport && authenticated {
		return fmt.Errorf(exchange.WarningAuthenticatedRequestWithoutCredentialsSet, cr.Name)
	}

	endpoint, err := cr.API.Endpoints.GetURL(ep)
	if err != nil {
		return err
	}

	if params == nil {
		params = map[string]interface{}{}
	}

	params["nonce"] = getNonce()
	params["request"] = apiRequest

	payload, err := json.Marshal(params)
	if err != nil {
		return errors.New("sendHTTPRequest: Unable to JSON request")
	}

	if cr.Verbose {
		log.Debugf(log.ExchangeSys, "Request JSON: %s", payload)
	}

	headers := make(map[string]string)
	if authenticated {
		headers["X-USER"] = cr.API.Credentials.ClientID
		hmac := crypto.GetHMAC(crypto.HashSHA256, payload, []byte(cr.API.Credentials.Key))
		headers["X-SIGNATURE"] = crypto.HexEncodeToString(hmac)
	}
	headers["Content-Type"] = "application/json"

	var rawMsg json.RawMessage
	err = cr.SendPayload(context.Background(), &request.Item{
		Method:        http.MethodPost,
		Path:          endpoint,
		Headers:       headers,
		Body:          bytes.NewBuffer(payload),
		Result:        &rawMsg,
		AuthRequest:   authenticated,
		NonceEnabled:  true,
		Verbose:       cr.Verbose,
		HTTPDebugging: cr.HTTPDebugging,
		HTTPRecording: cr.HTTPRecording,
	})
	if err != nil {
		return err
	}

	var genResp GenericResponse
	err = json.Unmarshal(rawMsg, &genResp)
	if err != nil {
		return err
	}

	if genResp.Code != success {
		return fmt.Errorf("%s SendHTTPRequest error: %s", cr.Name,
			genResp.Message)
	}

	return json.Unmarshal(rawMsg, result)
}
func getNonce() int64 {
	return time.Now().Unix()
}
