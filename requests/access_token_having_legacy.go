package requests

type AccessTokenHavingLegacy struct {
	AccessToken string `json:"accessToken" validate:"required"`
}
