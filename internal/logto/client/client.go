package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

type accessTokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type APIUsersPostReq struct {
	PrimaryPhone string `json:"primaryPhone"`
	PrimaryEmail string `json:"primaryEmail"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Name         string `json:"name"`
}

type APIUsersPostRes struct {
	Id string `json:"id"`
}

// NewClient creates a new instance of the Client struct.
// It takes an endpoint string as a parameter and returns a pointer to the Client struct.
func NewClient(endpoint string) *Client {
	return &Client{baseURL: endpoint, httpClient: &http.Client{}}
}

func (clnt *Client) AccessToken(user, pass, resource string) (string, error) {
	params := url.Values{}
	params.Set("grant_type", "client_credentials")
	params.Set("resource", resource)
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/oidc/token", clnt.baseURL),
		strings.NewReader(params.Encode()),
	)
	if err != nil {
		return "", fmt.Errorf("error in creating request %s ", err)
	}
	req.SetBasicAuth(user, pass)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := clnt.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error in making request %s", err)
	}
	if res.StatusCode != 200 {
		cnt, err := io.ReadAll(res.Body)
		if err != nil {
			return "", fmt.Errorf(
				"error in response and the reading the body %d %s",
				res.StatusCode,
				err,
			)
		}
		return "", fmt.Errorf(
			"unexpected error response %d %s",
			res.StatusCode,
			string(cnt),
		)
	}
	defer res.Body.Close()
	acresp := &accessTokenResp{}
	if err := json.NewDecoder(res.Body).Decode(acresp); err != nil {
		return "", fmt.Errorf("error in decoding json response %s", err)
func (clnt *Client) CreateUser(
	token string,
	user *APIUsersPostReq,
) (string, error) {
	var userId string
	content, err := json.Marshal(user)
	if err != nil {
		return userId, fmt.Errorf("error in converting to json %s", err)
	}
	ureq, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/api/users", clnt.baseURL),
		bytes.NewBuffer(content),
	)
	if err != nil {
		return userId, fmt.Errorf("error in making new request %s", err)
	}
	ureq.Header.Set("Content-Type", "application/json")
	ureq.Header.Set("Accept", "application/json")
	ureq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	uresp, err := clnt.reqToResponse(ureq)
	if err != nil {
		return userId, err
	}
	defer uresp.Body.Close()
	usr := &APIUsersPostRes{}
	if err := json.NewDecoder(uresp.Body).Decode(usr); err != nil {
		return userId, fmt.Errorf("error in decoding json response %s", err)
	}
	return usr.Id, nil
}
