package api

type UniversalAuthLoginRequest struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

type UniversalAuthLoginResponse struct {
	AccessToken string `json:"accessToken"`
}
