package paack

type Authenticate struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Audience     string `json:"audience"`
	GrantType    string `json:"grant_type"`
}

type AuthenticateResponse struct {
	Token string `json:"access_token"`
}
