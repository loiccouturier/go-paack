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
	authenticateHost string
	host             string
	clientId         string
	clientSecret     string
	token            string
}

func (c *client) CreateOrder(order Order) (*OrderResponse, *ApiError) {
	var result *OrderResponse

	apiError := c.call(http.MethodPost,  fmt.Sprintf("%s/public/v3/orders", c.host), order, result, true)
	if apiError != nil {
		return nil, apiError
	}

	return result, nil
}

func (c *client) RetrieveOrder(orderId string) (*OrderResponse, *ApiError) {
	var result *OrderResponse

	apiError := c.call(http.MethodGet, fmt.Sprintf("%s/public/v3/orders/%s", c.host, orderId), &struct{}{}, result, true)
	if apiError != nil {
		return nil, apiError
	}

	return result, nil
}

func (c *client) UpdateOrder(orderId string, order Order) (*UpdateResponse, *ApiError) {
	var result *UpdateResponse

	apiError := c.call(http.MethodPut, fmt.Sprintf("%s/public/v3/orders/%s", c.host, orderId), order, result, true)
	if apiError != nil {
		return nil, apiError
	}

	return result, nil
}

func (c *client) ReplaceOrderParcels(orderId string, parcels []Parcel) (*ReplaceParcelsResponse, *ApiError) {
	var result *ReplaceParcelsResponse

	apiError := c.call(http.MethodPut, fmt.Sprintf("%s/public/v3/orders/%s/parcels", c.host, orderId), Parcels{Parcels: parcels}, result, true)
	if apiError != nil {
		return nil, apiError
	}

	return result, nil
}

func (c *client) CancelOrder(orderId string) (*CancelResponse, *ApiError) {
	var result *CancelResponse

	apiError := c.call(http.MethodDelete, fmt.Sprintf("%s/public/v3/orders/%s", c.host, orderId), &struct{}{}, result, true)
	if apiError != nil {
		return nil, apiError
	}

	return result, nil
}

func (c *client) CreateLabel(label Label) ([]byte, *ApiError) {
	var result []byte

	apiError := c.call(http.MethodPost, fmt.Sprintf("%s/public/v3/label", c.host), label, &result, true)
	if apiError != nil {
		return nil, apiError
	}

	return result, nil
}

func (c *client) authenticate() *ApiError {
	var result AuthenticateResponse

	fmt.Println("#########")
	apiError := c.call(http.MethodPost, fmt.Sprintf("%s/oauth/token", c.authenticateHost), &Authenticate{ClientId: c.clientId, ClientSecret: c.clientSecret, Audience: "https://ggl-stg-gcp-gw", GrantType: "client_credentials"}, &result, false)
	if apiError != nil {
		return apiError
	}
	fmt.Println("#########1")

	c.token = "Bearer " + result.Token

	return nil
}

func (c *client) call(method, url string, body, result interface{}, needAuthentication bool) *ApiError {
	fmt.Println(url)

	// Json encode body
	if body == nil {
		body = ""
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return &ApiError{Err: err}
	}

	fmt.Println(string(jsonBody))

	httpClient := http.DefaultClient
	req, err := http.NewRequest(method, url, bytes.NewReader(jsonBody))
	if err != nil {
		return &ApiError{Err: err}
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-type", "application/json")

	if needAuthentication {
		if c.token == "" {
			apiError := c.authenticate()
			if apiError != nil {
				return apiError
			}
		}

		if c.token != "" {
			req.Header.Set("Authorization", c.token)
		}
	}

	response, err := httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return &ApiError{Err: err}
	}
	fmt.Println(response)

	if response.StatusCode == http.StatusUnauthorized {
		// Clear token and retry call
		c.token = ""
		return c.call(method, url, body, result, needAuthentication)
	} else {
		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return &ApiError{Err: err}
		}

		fmt.Println(string(responseBody))

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
				err = json.Unmarshal(responseBody, result)
				if err != nil {
					return &ApiError{Err: err}
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

func NewClient(host, authenticateHost, clientId, clientSecret string) Client {
	return &client{
		host:             host,
		authenticateHost: authenticateHost,
		clientId:         clientId,
		clientSecret:     clientSecret,
	}
}
