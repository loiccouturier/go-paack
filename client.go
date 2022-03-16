package paack

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var ErrBadRequest = errors.New("ErrBadRequest")
var ErrForbidden = errors.New("ErrForbidden")
var ErrNotFound = errors.New("ErrNotFound")
var ErrUndefined = errors.New("ErrUndefined")

type Client interface {
	CreateOrder(order Order) (*OrderResponse, *ApiError)
	RetrieveOrder(orderId string) (*OrderResponse, *ApiError)
	UpdateOrder(orderId string, order Order) (*UpdateResponse, *ApiError)
	ReplaceOrderParcels(orderId string, parcels []Parcel) (*ReplaceParcelsResponse, *ApiError)
	CancelOrder(orderId string) (*CancelResponse, *ApiError)
	CreateLabel(label Label) ([]byte, *ApiError)
}

type client struct {
	host                     string
	authenticateHost         string
	token                    string
	clientId                 string
	clientSecret             string
	audience                 string
	authenticateHostForLabel string
	tokenForLabel            string
	clientIdForLabel         string
	clientSecretForLabel     string
	audienceForLabel         string
	labelHost                string
	debug                    bool
}

func (c *client) CreateOrder(order Order) (*OrderResponse, *ApiError) {
	if c.debug {
		log.Println("paack.client", "CreateOrder")
	}

	var result OrderResponse

	apiError := c.call(http.MethodPost, fmt.Sprintf("%s/public/v3/orders", c.host), order, &result, true, false)
	if apiError != nil {
		return nil, apiError
	}

	return &result, nil
}

func (c *client) RetrieveOrder(orderId string) (*OrderResponse, *ApiError) {
	if c.debug {
		log.Println("paack.client", "RetrieveOrder")
	}

	var result OrderResponse

	apiError := c.call(http.MethodGet, fmt.Sprintf("%s/public/v3/orders/%s", c.host, orderId), &struct{}{}, &result, true, false)
	if apiError != nil {
		return nil, apiError
	}

	return &result, nil
}

func (c *client) UpdateOrder(orderId string, order Order) (*UpdateResponse, *ApiError) {
	if c.debug {
		log.Println("paack.client", "UpdateOrder")
	}

	var result UpdateResponse

	apiError := c.call(http.MethodPut, fmt.Sprintf("%s/public/v3/orders/%s", c.host, orderId), order, &result, true, false)
	if apiError != nil {
		return nil, apiError
	}

	return &result, nil
}

func (c *client) ReplaceOrderParcels(orderId string, parcels []Parcel) (*ReplaceParcelsResponse, *ApiError) {
	if c.debug {
		log.Println("paack.client", "ReplaceOrderParcels")
	}

	var result ReplaceParcelsResponse

	apiError := c.call(http.MethodPut, fmt.Sprintf("%s/public/v3/orders/%s/parcels", c.host, orderId), Parcels{Parcels: parcels}, &result, true, false)
	if apiError != nil {
		return nil, apiError
	}

	return &result, nil
}

func (c *client) CancelOrder(orderId string) (*CancelResponse, *ApiError) {
	if c.debug {
		log.Println("paack.client", "CancelOrder")
	}

	var result CancelResponse

	apiError := c.call(http.MethodDelete, fmt.Sprintf("%s/public/v3/orders/%s", c.host, orderId), &struct{}{}, &result, true, false)
	if apiError != nil {
		return nil, apiError
	}

	return &result, nil
}

func (c *client) CreateLabel(label Label) ([]byte, *ApiError) {
	if c.debug {
		log.Println("paack.client", "CreateLabel")
	}

	var result LabelResponse

	apiError := c.call(http.MethodPost, fmt.Sprintf("%s/v3/labels", c.labelHost), label, &result, true, true)
	if apiError != nil {
		return nil, apiError
	}

	return result.Content, nil
}

func (c *client) authenticate() *ApiError {
	if c.debug {
		log.Println("paack.client", "authenticate")
	}

	var result AuthenticateResponse

	apiError := c.call(http.MethodPost, fmt.Sprintf("%s/oauth/token", c.authenticateHost), &Authenticate{ClientId: c.clientId, ClientSecret: c.clientSecret, Audience: c.audience, GrantType: "client_credentials"}, &result, false, false)
	if apiError != nil {
		return apiError
	}

	c.token = "Bearer " + result.Token

	return nil
}

func (c *client) authenticateForLabel() *ApiError {
	if c.debug {
		log.Println("paack.client", "authenticateForLabel")
	}

	var result AuthenticateResponse

	apiError := c.call(http.MethodPost, fmt.Sprintf("%s/oauth/token", c.authenticateHostForLabel), &Authenticate{ClientId: c.clientIdForLabel, ClientSecret: c.clientSecretForLabel, Audience: c.audienceForLabel, GrantType: "client_credentials"}, &result, false, false)
	if apiError != nil {
		return apiError
	}

	c.tokenForLabel = "Bearer " + result.Token

	return nil
}

func (c *client) call(method, url string, body, result interface{}, needAuthentication bool, isLabel bool) *ApiError {
	// Json encode body
	if body == nil {
		body = ""
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return &ApiError{Err: err}
	}
	if c.debug {
		log.Println("paack.client", "payload: ", string(jsonBody))
	}

	httpClient := http.DefaultClient
	req, err := http.NewRequest(method, url, bytes.NewReader(jsonBody))
	if err != nil {
		return &ApiError{Err: err}
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-type", "application/json")

	if needAuthentication {
		if c.debug {
			log.Println("paack.client", "call need authentication")
		}

		if !isLabel && c.token == "" {
			if c.debug {
				log.Println("paack.client", "start authentication")
			}
			apiError := c.authenticate()
			if apiError != nil {
				return apiError
			}
		}

		if isLabel && c.tokenForLabel == "" {
			if c.debug {
				log.Println("paack.client", "start authentication for label")
			}
			apiError := c.authenticateForLabel()
			if apiError != nil {
				return apiError
			}
		}

		if !isLabel && c.token != "" {
			if c.debug {
				log.Println("paack.client", "set token")
			}
			req.Header.Set("Authorization", c.token)
		}

		if isLabel && c.tokenForLabel != "" {
			if c.debug {
				log.Println("paack.client", "set token for label")
			}
			req.Header.Set("Authorization", c.tokenForLabel)
		}
	}

	response, err := httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return &ApiError{Err: err}
	}

	if response.StatusCode == http.StatusUnauthorized {
		if c.debug {
			log.Println("paack.client", "Http status 401 received for url ", url)
		}
		// Clear token and retry call
		c.token = ""
		c.tokenForLabel = ""
		return c.call(method, url, body, result, needAuthentication, isLabel)
	} else {
		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return &ApiError{Err: err}
		}

		defer func() {
			err := response.Body.Close()
			if err != nil {
				panic(fmt.Sprintf("can not close response body: %v\n", err))
			}
		}()

		err = nil
		if response.StatusCode == http.StatusBadRequest {
			if c.debug {
				log.Println("paack.client", "Http status 400 received for url ", url, "response", string(responseBody))
			}
			err = ErrBadRequest
		} else if response.StatusCode == http.StatusForbidden {
			if c.debug {
				log.Println("paack.client", "Http status 403 received for url ", url)
			}
			err = ErrForbidden
		} else if response.StatusCode == http.StatusNotFound {
			if c.debug {
				log.Println("paack.client", "Http status 404 received for url ", url)
			}
			err = ErrNotFound
		} else if response.StatusCode != http.StatusOK {
			if c.debug {
				log.Println("paack.client", fmt.Sprintf("Http status %d received for url %s", response.StatusCode, url))
			}
			err = ErrUndefined
		}

		if err == nil {
			if c.debug {
				log.Println("paack.client", "No error for url ", url)
			}
			if result != nil {
				if !isLabel {
					if c.debug {
						log.Println("paack.client", "Request is not for a label, unmarshal the result")
					}
					err = json.Unmarshal(responseBody, result)
					if err != nil {
						return &ApiError{Err: err}
					}
				} else {
					if c.debug {
						log.Println("paack.client", "Request is for a label, return content")
					}
					r := result.(*LabelResponse)
					r.Content = responseBody
				}
			}
		} else {
			if c.debug {
				log.Println("paack.client", "error received for url ", url)
			}
			var apiError ApiError
			err2 := json.Unmarshal(responseBody, &apiError)
			if err2 != nil {
				return &ApiError{Err: err2}
			}
			apiError.Err = err
			return &apiError
		}
	}

	return nil
}

func NewClient(host, authenticateHost, authenticateHostForLabel, labelHost, audience, audienceForLabel, clientId, clientSecret, clientIdForLabel, clientSecretForLabel string, debug bool) Client {
	return &client{
		host:                     host,
		authenticateHost:         authenticateHost,
		authenticateHostForLabel: authenticateHostForLabel,
		labelHost:                labelHost,
		clientId:                 clientId,
		clientSecret:             clientSecret,
		clientIdForLabel:         clientIdForLabel,
		clientSecretForLabel:     clientSecretForLabel,
		audience:                 audience,
		audienceForLabel:         audienceForLabel,
		debug:                    debug,
	}
}
