package pkg

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"golang.org/x/oauth2"
)

var verifier = oauth2.GenerateVerifier()

type GoogleOAuth struct {
	config oauth2.Config
}

type UserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
	Name          string `json:"name"`
}

func NewGoogleOAuth(clientID, clientSecret, redirectURL string) *GoogleOAuth {
	return &GoogleOAuth{
		config: oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
			Endpoint: oauth2.Endpoint{
				AuthURL:   "https://accounts.google.com/o/oauth2/auth",
				TokenURL:  "https://oauth2.googleapis.com/token",
				AuthStyle: oauth2.AuthStyleInParams,
			},
		},
	}
}

func (g *GoogleOAuth) GetAuthCodeURL(state string) string {
	return g.config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))
}

func (g *GoogleOAuth) GetUserInfo(ctx context.Context, code string) (UserInfo, error) {
	token, err := g.config.Exchange(ctx, code, oauth2.AccessTypeOffline, oauth2.VerifierOption(verifier))
	if err != nil {
		return UserInfo{}, err
	}
	request, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	if err != nil {
		return UserInfo{}, err
	}
	request.Header.Set("Authorization", "Bearer "+token.AccessToken)
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return UserInfo{}, err
	}
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return UserInfo{}, err
	}
	var userInfo UserInfo
	err = json.Unmarshal(contents, &userInfo)
	if err != nil {
		return UserInfo{}, err
	}
	return userInfo, nil
}
