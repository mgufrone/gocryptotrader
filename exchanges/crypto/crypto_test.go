package crypto

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/thrasher-corp/gocryptotrader/config"
	exchange "github.com/thrasher-corp/gocryptotrader/exchanges"
)

// Please supply your own keys here to do authenticated endpoint testing
const (
	apiKey                  = ""
	apiSecret               = ""
	apiSandbox              = "https://uat-api.3ona.co/"
	canManipulateRealOrders = false
)

var cr Crypto

func TestMain(m *testing.M) {
	cr.SetDefaults()
	cfg := config.GetConfig()
	err := cfg.LoadConfig("../../testdata/configtest.json", true)
	if err != nil {
		log.Fatal(err)
	}

	exchCfg, err := cfg.GetExchangeConfig("Crypto")
	if err != nil {
		log.Fatal(err)
	}

	exchCfg.API.AuthenticatedSupport = true
	exchCfg.API.AuthenticatedWebsocketSupport = true
	exchCfg.API.Credentials.Key = apiKey
	exchCfg.API.Credentials.Secret = apiSecret

	err = cr.Setup(exchCfg)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

// Ensures that this exchange package is compatible with IBotExchange
func TestInterface(t *testing.T) {
	var e exchange.IBotExchange
	if e = new(Crypto); e == nil {
		t.Fatal("unable to allocate exchange")
	}
}

func areTestAPIKeysSet() bool {
	return cr.ValidateAPICredentials()
}

func TestCrypto_GetInstruments(t *testing.T) {
	_, err := cr.GetInstruments()
	if err != nil {
		t.Error("GetInstruments() error", err)
	}
}

func TestCrypto_FetchTradablePairs(t *testing.T) {
	pairs, err := cr.FetchTradablePairs("BTC-USD")
	if err != nil {
		t.Error("GetInstruments() error", err)
	}
	fmt.Println(pairs)
}

// Implement tests for API endpoints below
