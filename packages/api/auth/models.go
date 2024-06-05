package api

type UniversalAuthLoginRequest struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

type UniversalAuthLoginResponse struct {
	AccessToken string `json:"accessToken"`
}
