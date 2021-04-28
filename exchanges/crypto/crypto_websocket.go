package crypto

import (
	"encoding/json"
	"errors" //nolint:gci internal package
	"github.com/gorilla/websocket"
	exchange "github.com/thrasher-corp/gocryptotrader/exchanges"
	"github.com/thrasher-corp/gocryptotrader/exchanges/stream"
	"net/http" //nolint:gci internal package
)

func (cr *Crypto) GetWsInstruments() (list InstrumentResult, err error) {
	var response InstrumentResponse
	request := wsRequest{
		Method: cryptoInstruments,
		Nonce:  getNonce(),
	}
	resp, err := cr.Websocket.Conn.SendMessageReturnResponse(request.Nonce, request)
	if err != nil {
		return
	}
	err = json.Unmarshal(resp, &response)
	if err == nil {
		list = response.Result
	}
	return
}

func (cr *Crypto) wsReadData() {
	cr.Websocket.Wg.Add(1)
	defer cr.Websocket.Wg.Done()

	for {
		select {
		case <-cr.Websocket.ShutdownC:
			return
		default:
			resp := cr.Websocket.Conn.ReadMessage()
			if resp.Raw == nil {
				return
			}

			err := cr.wsHandleData(resp.Raw)
			if err != nil {
				cr.Websocket.DataHandler <- err
			}
		}
	}
}

func (cr *Crypto) WsConnect() error {
	if !cr.Websocket.IsEnabled() || !cr.IsEnabled() {
		return errors.New(stream.WebsocketNotEnabled)
	}
	var dialer websocket.Dialer
	err := cr.Websocket.Conn.Dial(&dialer, http.Header{})
	if err != nil {
		return err
	}

	go cr.wsReadData()
	if cr.GetAuthenticatedAPISupport(exchange.WebsocketAuthentication) {
		err = cr.WsAuth()
		if err != nil {
			cr.Websocket.DataHandler <- err
			cr.Websocket.SetCanUseAuthenticatedEndpoints(false)
		}
	}
	return nil
}

func (cr *Crypto) GenerateDefaultSubscriptions() (channels []stream.ChannelSubscription, err error) {
	return
}

func (cr *Crypto) WsAuth() (err error) {
	return
}

func (cr *Crypto) wsHandleData(raw []byte) (err error) {
	return
}
