package response

type AuthorizeResponse struct {
	AccessToken string `json:"access_token"`
	Username    string `json:"username"`
}
