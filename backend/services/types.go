package services

type GitHubRes struct {
	Id        *int    `json:"id"`
	OauthId   *string `copier:"OauthId"`
	AvatarURL *string `json:"avatar_url" copier:"Picture"`
	Email     *string `json:"email" copier:"Email"`
}

type GithubEmailsRes struct {
	Email []GithubEmails
}
type GithubEmails struct {
	Email    *string `json:"email"`
	Primary  *bool   `json:"primary"`
	Verified *bool   `json:"verified"`
}

type GoogleRes struct {
	OauthId *string `copier:"OauthId"`
	Id      *string `json:"id"`
	Picture *string `json:"picture" copier:"Picture"`
	Email   *string `json:"email" copier:"Email"`
}
