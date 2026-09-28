package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/bhumika019579/Vouch/server/internal/config"
	"github.com/bhumika019579/Vouch/server/internal/models"
	"github.com/bhumika019579/Vouch/server/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type githubTokenResponse struct {
	AccessToken string `json:"access_token"`
}

func GithubLogin(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		redirectURL := "https://github.com/login/oauth/authorize" +
			"?client_id=" + cfg.GithubClientID +
			"&redirect_uri=" + cfg.GithubRedirectURL +
			"&scope=read:user,user:email,public_repo"
		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
	}
}
func GithubCallback(cfg *config.Config, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Query("code")
		if code == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing code"})
			return
		}
		tokenReqURL := "https://github.com/login/oauth/access_token"
		form := url.Values{}
		form.Set("client_id", cfg.GithubClientID)
		form.Set("client_secret", cfg.GithubClientSecret)
		form.Set("code", code)
		form.Set("redirect_uri", cfg.GithubRedirectURL)
		req, _ := http.NewRequest("POST", tokenReqURL, nil)
		req.URL.RawQuery = form.Encode()
		req.Header.Set("Accept", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reach github"})
			return
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		var tokenResp githubTokenResponse
		json.Unmarshal(body, &tokenResp)
		if tokenResp.AccessToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to get access token"})
			return
		}
		req2, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
		req2.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)
		req2.Header.Set("Accept", "application/json")
		resp2, err := http.DefaultClient.Do(req2)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch github user"})
			return
		}
		defer resp2.Body.Close()
		body2, _ := io.ReadAll(resp2.Body)
		var githubUser struct {
			ID        int64  `json:"id"`
			Login     string `json:"login"`
			Name      string `json:"name"`
			Email     string `json:"email"`
			AvatarURL string `json:"avatar_url"`
			Bio       string `json:"bio"`
		}
		json.Unmarshal(body2, &githubUser)
		var user models.User
		result := db.Where("github_id = ?", githubUser.ID).First(&user)

		if result.Error != nil {
			user = models.User{
				Name:              githubUser.Name,
				Email:             githubUser.Email,
				GithubID:          githubUser.ID,
				GithubUsername:    githubUser.Login,
				GithubAccessToken: tokenResp.AccessToken,
				AvatarURL:         githubUser.AvatarURL,
				Bio:               githubUser.Bio,
			}
			db.Create(&user)
		} else {
			user.GithubAccessToken = tokenResp.AccessToken
			user.AvatarURL = githubUser.AvatarURL
			user.Bio = githubUser.Bio
			db.Save(&user)
		}
		jwtToken, err := utils.GenerateJWT(user.ID, cfg.JWTSecret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
			return
		}

		frontendURL := "http://localhost:5173/?token=" + jwtToken
		c.Redirect(http.StatusTemporaryRedirect, frontendURL)

	}
}
