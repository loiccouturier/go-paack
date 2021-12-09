package paack

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var ErrBadRequest = errors.New("ErrBadRequest")
var ErrForbidden = errors.New("ErrForbidden")
var ErrNotFound = errors.New("ErrNotFound")

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
}

func (c *client) CreateOrder(order Order) (*OrderResponse, *ApiError) {
	var result OrderResponse

	apiError := c.call(http.MethodPost, fmt.Sprintf("%s/public/v3/orders", c.host), order, &result, true, false)
	if apiError != nil {
		return nil, apiError
	}

	return &result, nil
}

func (c *client) RetrieveOrder(orderId string) (*OrderResponse, *ApiError) {
	var result OrderResponse

	apiError := c.call(http.MethodGet, fmt.Sprintf("%s/public/v3/orders/%s", c.host, orderId), &struct{}{}, &result, true, false)
	if apiError != nil {
		return nil, apiError
	}

	return &result, nil
}

func (c *client) UpdateOrder(orderId string, order Order) (*UpdateResponse, *ApiError) {
	var result UpdateResponse

	apiError := c.call(http.MethodPut, fmt.Sprintf("%s/public/v3/orders/%s", c.host, orderId), order, &result, true, false)
	if apiError != nil {
		return nil, apiError
	}

	return &result, nil
}

func (c *client) ReplaceOrderParcels(orderId string, parcels []Parcel) (*ReplaceParcelsResponse, *ApiError) {
	var result ReplaceParcelsResponse

	apiError := c.call(http.MethodPut, fmt.Sprintf("%s/public/v3/orders/%s/parcels", c.host, orderId), Parcels{Parcels: parcels}, &result, true, false)
	if apiError != nil {
		return nil, apiError
	}

	return &result, nil
}

func (c *client) CancelOrder(orderId string) (*CancelResponse, *ApiError) {
	var result CancelResponse

	apiError := c.call(http.MethodDelete, fmt.Sprintf("%s/public/v3/orders/%s", c.host, orderId), &struct{}{}, &result, true, false)
	if apiError != nil {
		return nil, apiError
	}

	return &result, nil
}

func (c *client) CreateLabel(label Label) ([]byte, *ApiError) {
	var result LabelResponse

	apiError := c.call(http.MethodPost, fmt.Sprintf("%s/v3/labels", c.labelHost), label, &result, true, true)
	if apiError != nil {
		return nil, apiError
	}

	return result.Content, nil
}

func (c *client) authenticate() *ApiError {
	var result AuthenticateResponse

	apiError := c.call(http.MethodPost, fmt.Sprintf("%s/oauth/token", c.authenticateHost), &Authenticate{ClientId: c.clientId, ClientSecret: c.clientSecret, Audience: c.audience, GrantType: "client_credentials"}, &result, false, false)
	if apiError != nil {
		return apiError
	}

	c.token = "Bearer " + result.Token

	return nil
}

func (c *client) authenticateForLabel() *ApiError {
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

	httpClient := http.DefaultClient
	req, err := http.NewRequest(method, url, bytes.NewReader(jsonBody))
	if err != nil {
		return &ApiError{Err: err}
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-type", "application/json")

	if needAuthentication {
		if !isLabel && c.token == "" {
			apiError := c.authenticate()
			if apiError != nil {
				return apiError
			}
		}

		if isLabel && c.tokenForLabel == "" {
			apiError := c.authenticateForLabel()
			if apiError != nil {
				return apiError
			}
		}

		if !isLabel && c.token != "" {
			req.Header.Set("Authorization", c.token)
		}

		if isLabel && c.tokenForLabel != "" {
			req.Header.Set("Authorization", c.tokenForLabel)
		}
	}

	response, err := httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return &ApiError{Err: err}
	}

	if response.StatusCode == http.StatusUnauthorized {
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
			err = ErrBadRequest
		} else if response.StatusCode == http.StatusForbidden {
			err = ErrForbidden
		} else if response.StatusCode == http.StatusNotFound {
			err = ErrNotFound
		}

		if err == nil {
			if result != nil {
				if !isLabel {
					err = json.Unmarshal(responseBody, result)
					if err != nil {
						return &ApiError{Err: err}
					}
				} else {
					r := result.(*LabelResponse)
					r.Content = responseBody
				}
			}
		} else {
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

func NewClient(host, authenticateHost, authenticateHostForLabel, labelHost, audience, audienceForLabel, clientId, clientSecret, clientIdForLabel, clientSecretForLabel string) Client {
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
	}
}
