package auth

import (
	"net/http"
	"net/url"
	"os"
	"tap-to-park/database"

	"github.com/carlmjohnson/requests"
	"github.com/gin-gonic/gin"
)

// GET https://github.com/login/oauth/authorize
// POST https://github.com/login/oauth/access_token

type Google struct{}

func (Google) Initialize(c *gin.Context) {
	vals := url.Values{}
	vals.Set("scope", "openid https://www.googleapis.com/auth/userinfo.email")
	vals.Set("response_type", "code")
	vals.Set("client_id", os.Getenv("GOOGLE_CLIENT_ID"))
	vals.Set("redirect_uri", os.Getenv("FRONTEND_HOST")+"/auth/google")
	vals.Set("access_type", "offline")
	c.Redirect(http.StatusMovedPermanently, "https://accounts.google.com/o/oauth2/v2/auth?"+vals.Encode())
}

type GoogleAccessTokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
	GrantType    string `json:"grant_type"`
	RedirectUri  string `json:"redirect_uri"`
}

type GoogleAccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    uint64 `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

type GoogleTokenInfoResponse struct {
	Email string `json:"email"`
	Sub   string `json:"sub"`
}

func (Google) Callback(c *gin.Context) *database.User {

	code := c.Query("code")

	body := GoogleAccessTokenRequest{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Code:         code,
		GrantType:    "authorization_code",
		RedirectUri:  os.Getenv("FRONTEND_HOST") + "/auth/google",
	}

	var res GoogleAccessTokenResponse
	err := requests.
		URL("https://oauth2.googleapis.com").
		Path("/token").
		ContentType("application/json").
		Accept("application/json").
		BodyJSON(&body).
		ToJSON(&res).
		Fetch(c)

	if err != nil {
		return nil
	}

	var tokenInfo GoogleTokenInfoResponse
	err = requests.
		URL("https://oauth2.googleapis.com").
		Path("/tokeninfo").
		Header("Authorization", "Bearer "+res.AccessToken).
		ContentType("application/json").
		Accept("application/json").
		ToJSON(&tokenInfo).
		Fetch(c)

	user := &database.User{
		Email:      tokenInfo.Email,
		Type:       "google",
		ExternalID: tokenInfo.Sub,
	}

	if result := database.Db.Model(user).Where("external_id = ?", user.ExternalID).First(user); result.Error != nil {
		if result := database.Db.Create(user); result.Error != nil {
			return nil
		}
	}

	return user
}
