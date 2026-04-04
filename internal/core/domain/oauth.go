package domain

type OauthProvider string

const (
	OauthProviderGoogle   OauthProvider = "google"
	OauthProviderFacebook OauthProvider = "facebook"
	OauthProviderApple    OauthProvider = "apple"
	OauthProviderBasic    OauthProvider = "basic"
)

type BasicLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SocialLogin struct {
	Provider OauthProvider `json:"provider"`
	Token    string        `json:"token"`
	Nonce    string        `json:"nonce"`
}

type OauthUser struct {
	ID          int64   `json:"id"`
	Email       string  `json:"email"`
	Password    *string `json:"password"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	DisplayName string  `json:"display_name"`
	Bio         *string `json:"bio,omitempty"`
	Avatar      *string `json:"avatar,omitempty"`
}

type GoogleUserInfo struct {
	ProviderID  string
	Email       string
	FirstName   string
	LastName    string
	DisplayName string
	Avatar      *string
}
