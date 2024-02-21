package server

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/RabiesDev/request-helper"
	"github.com/tidwall/gjson"
)

func (server *Server) OAuth2Token(authorizeCode string) (string, string) {
	endpoint := "https://discord.com/api/v10/oauth2/token"
	payload := strings.NewReader(url.Values{
		"client_id":     {server.Environment.ClientID},
		"client_secret": {server.Environment.ClientSecret},
		"redirect_uri":  {"http://127.0.0.1:3030/authorize"},
		"grant_type":    {"authorization_code"},
		"code":          {authorizeCode},
	}.Encode())

	request := requests.Post(endpoint, payload)
	requests.SetHeaders(request, map[string]string{
		"content-type": "application/x-www-form-urlencoded",
	})

	body, response, err := requests.DoAndReadByte(http.DefaultClient, request)
	if err != nil || response.StatusCode >= 400 {
		return "", ""
	}

	return gjson.GetBytes(body, "access_token").String(), gjson.GetBytes(body, "refresh_token").String()
}

func (server *Server) Profile(accessToken string) (string, string) {
	request := requests.Get("https://discord.com/api/v10/users/@me")
	requests.SetHeaders(request, map[string]string{
		"authorization": "Bearer " + accessToken,
	})

	body, response, err := requests.DoAndReadByte(http.DefaultClient, request)
	if err != nil || response.StatusCode >= 400 {
		return "", ""
	}

	return gjson.GetBytes(body, "id").String(), gjson.GetBytes(body, "username").String()
}
