package auth

import (
	"net/http"
	"net/url"
	"os"
	"strconv"
	"tap-to-park/database"

	"github.com/carlmjohnson/requests"
	"github.com/gin-gonic/gin"
)

// GET https://github.com/login/oauth/authorize
// POST https://github.com/login/oauth/access_token

type Github struct{}

func (Github) Initialize(c *gin.Context) {
	vals := url.Values{}
	vals.Set("client_id", os.Getenv("GITHUB_CLIENT_ID"))
	vals.Set("scope", "user")
	c.Redirect(http.StatusMovedPermanently, "https://github.com/login/oauth/authorize?"+vals.Encode())
}

type GithubAccessTokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
	// RedirectUri  string `json:"redirect_uri"`
}

type GithubAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

type GithubUserResponse struct {
	Login string `json:"login"`
	ID    int64  `json:"id"`
}

func (Github) Callback(c *gin.Context) *database.User {

	code := c.Query("code")

	body := GithubAccessTokenRequest{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		Code:         code,
		// RedirectUri:  "",
	}

	var res GithubAccessTokenResponse
	err := requests.
		URL("https://github.com").
		Path("/login/oauth/access_token").
		Accept("application/json").
		BodyJSON(&body).
		ToJSON(&res).
		Fetch(c)

	if err != nil {
		return nil
	}

	var info GithubUserResponse
	err = requests.URL("https://api.github.com").
		Path("/user").
		Accept("application/json").
		Header("Authorization", "Bearer "+res.AccessToken).
		Header("X-OAuth-Scopes", "user").
		ToJSON(&info).
		Fetch(c)

	if err != nil {
		return nil
	}

	user := &database.User{
		Email:      info.Login,
		Type:       "github",
		ExternalID: strconv.FormatInt(info.ID, 10),
	}

	if result := database.Db.Model(user).Where("external_id = ?", user.ExternalID).First(user); result.Error != nil {
		if result := database.Db.Create(user); result.Error != nil {
			return nil
		}
	}

	return user
}
