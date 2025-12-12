package helpers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strings"

	"github.com/1107-adishjain/golang-jwt/internal/config"
)

// Generate a random code_verifier for PKCE
func GenerateCodeVerifier() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~" //set of allowed characters 
	length := 64
	b := make([]byte, length)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[num.Int64()]
	}
	return string(b), nil
}

// Create a code_challenge from code_verifier
func GenerateCodeChallenge(verifier string) string {
	h := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(h[:])
}

// Returns the Google OAuth URL and the code_verifier
func GetGoogleOAuthURL(cfg *config.Config) (string, string, error) {
	codeVerifier, err := GenerateCodeVerifier()
	if err != nil {
		return "", "", err
	}
	codeChallenge := GenerateCodeChallenge(codeVerifier)
	base := "https://accounts.google.com/o/oauth2/v2/auth"
	params := url.Values{}
	params.Add("client_id", cfg.Client_ID)
	params.Add("redirect_uri", cfg.Redirect_URI)
	params.Add("response_type", "code")
	params.Add("scope", "openid email profile")
	params.Add("state", "random_state_string")
	params.Add("code_challenge", codeChallenge)
	params.Add("code_challenge_method", "S256")
	return base + "?" + params.Encode(), codeVerifier, nil
}

func ExchangeCodeForTokens(code string, codeVerifier string, cfg *config.Config) (string, string, string, error) {
	// here we will make a post request to google oauth2 token endpoint
	// to exchange the code for access token and id token

	endpoint := "https://oauth2.googleapis.com/token"
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", cfg.Client_ID)
	data.Set("client_secret", cfg.Client_Secret)
	data.Set("redirect_uri", cfg.Redirect_URI)
	data.Set("grant_type", "authorization_code")
	data.Set("code_verifier", codeVerifier)

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return "", "", "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 3. Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", "", err
	}
	defer resp.Body.Close()

	// 4. Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", err
	}

	// 5. Parse the JSON response to extract tokens
	type tokenResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		IdToken      string `json:"id_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
	}
	var tr tokenResponse
	if err := json.Unmarshal(body, &tr); err != nil {
		return "", "", "", err
	}

	// 6. Return the token These can be used to get user info and manage sessions
	return tr.AccessToken, tr.RefreshToken, tr.IdToken, nil
}

type UserInfo struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
}

func GetUserInfoFromGoogle(accessToken string) (UserInfo, error) {
	// here we will make a get request to google oauth2 userinfo endpoint
	// to get the user info using the access token

	endpoint := "https://www.googleapis.com/oauth2/v2/userinfo"
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return UserInfo{}, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return UserInfo{}, err
	}
	defer resp.Body.Close()
	// read the response.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return UserInfo{}, err
	}

	// we will use struct here to store the user info

	var userInfo UserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return UserInfo{}, err
	}
	return userInfo, nil
}
