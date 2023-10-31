package socials

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/agustfricke/x-clone/config"
	"github.com/agustfricke/x-clone/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func ConfigGoogle() *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     config.Config("CLIENT_GOOGLE"),
		ClientSecret: config.Config("SECRET_GOOGLE"),
		RedirectURL:  config.Config("REDIRECT_GOOGLE"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint: google.Endpoint,
	}
	return conf
}

func GetGoogleResponse(token string) models.GoogleResponse {
    reqURL, err := url.Parse("https://www.googleapis.com/oauth2/v1/userinfo")
    if err != nil {
        panic(err)
    }
    ptoken := fmt.Sprintf("Bearer %s", token)
    res := &http.Request{
        Method: "GET",
        URL:    reqURL,
        Header: map[string][]string{
            "Authorization": {ptoken}},
    }
    req, err := http.DefaultClient.Do(res)
    if err != nil {
        panic(err)
    }
    defer req.Body.Close()
    body, err := io.ReadAll(req.Body)
    if err != nil {
        panic(err)
    }
    var data models.GoogleResponse
    err = json.Unmarshal(body, &data)
    if err != nil {
        panic(err)
    }
    return data
}

func ConfigGitHub() *oauth2.Config {
	conf := &oauth2.Config{
        ClientID:     config.Config("CLIENT_GITHUB"),
        ClientSecret: config.Config("SECRET_GITHUB"),
        RedirectURL:  "http://localhost:8080/auth/github/callback",
        Scopes:       []string{"user"},
        Endpoint:     oauth2.Endpoint{
            AuthURL:  "https://github.com/login/oauth/authorize",
            TokenURL: "https://github.com/login/oauth/access_token",
        },
      }
	return conf
}

func GetGitHubResponse(token string) models.GitHubResponse {

    reqURL, err := url.Parse("https://api.github.com/user")

    if err != nil {
        panic(err)
    }
    ptoken := fmt.Sprintf("Bearer %s", token)
    res := &http.Request{
        Method: "GET",
        URL:    reqURL,
        Header: map[string][]string{
            "Authorization": {ptoken}},
    }
    req, err := http.DefaultClient.Do(res)
    if err != nil {
        panic(err)
    }
    defer req.Body.Close()
    body, err := io.ReadAll(req.Body)
    if err != nil {
        panic(err)
    }
    var data models.GitHubResponse
    err = json.Unmarshal(body, &data)
    if err != nil {
        panic(err)
    }
    return data
}
