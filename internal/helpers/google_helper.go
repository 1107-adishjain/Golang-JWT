package helpers

import (
	"net/http"
	"net/url"
	"strings"
	"encoding/json"
	"io/ioutil"

	"github.com/1107-adishjain/golang-jwt/internal/config"
)

func GetGoogleOAuthURL(cfg *config.Config) string {
	base:= "https://accounts.google.com/o/oauth2/v2/auth"
	params:= url.Values{}
	params.Add("client_id", cfg.Client_ID)
	params.Add("redirect_uri", cfg.Redirect_URI)
	params.Add("response_type", "code")
	params.Add("scope", "openid email profile")
	params.Add("state", "random_state_string") 
	params.Add("code_challenge", "CODE_CHALLENGE")
	params.Add("code_challenge_method", "S256")
	return base + "?" + params.Encode()
}

func ExchangeCodeForTokens(code string, cfg *config.Config) (string, string, string, error){
	// here we will make a post request to google oauth2 token endpoint
	// to exchange the code for access token and id token

	endpoint:= "https://oauth2.googleapis.com/token"
	data:= url.Values{}
	data.Set("code", code)
	data.Set("client_id", cfg.Client_ID)
	data.Set("client_secret", cfg.Client_Secret)
	data.Set("redirect_uri", cfg.Redirect_URI)
	data.Set("grant_type", "authorization_code")
	data.Set("code_verifier", "CODE_VERIFIER")

	req,err:= http.NewRequest("POST",endpoint,strings.NewReader(data.Encode()))
	if err!=nil{
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
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", "", "", err
    }


	// 5. Parse the JSON response to extract tokens
    //    Google returns access_token, refresh_token, id_token, etc.
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

    // 6. Return the tokens
    //    These can be used to get user info and manage sessions
    return tr.AccessToken, tr.RefreshToken, tr.IdToken, nil
}

func GetUserInfoFromGoogle(accessToken string) (map[string]interface{}, error){
	// here we will make a get request to google oauth2 userinfo endpoint
	// to get the user info using the access token

	endpoint:= "https://www.googleapis.com/oauth2/v2/userinfo"
	req,err:= http.NewRequest("GET", endpoint, nil)
	if err!= nil{
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+ accessToken)
	client:= &http.Client{}
	resp, err:= client.Do(req)
	if err!= nil{
		return nil, err
	}
	defer resp.Body.Close()

	body, err:= ioutil.ReadAll(resp.Body)
	if err!= nil{
		return nil, err
	}		

	var userInfo map[string]interface{}
	if err:= json.Unmarshal(body, &userInfo); err!= nil{
		return nil, err
	}
	return userInfo, nil	
}