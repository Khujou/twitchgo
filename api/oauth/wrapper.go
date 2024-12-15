package oauth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type OAuthToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	ExpiryTime  int64  `json:"expiry_time"`
	//TokenType   string `json:"token_type"`
}

func check(reason string, err error) {
	if err != nil {
		log.Fatal(reason, err)
		panic(err)
	}
}

func fetch(req *http.Request) []byte {
	//fmt.Println("fetching...")
	client := &http.Client{}
	resp, err := client.Do(req)
	check("failed doing", err)
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)

	body, err := io.ReadAll(resp.Body)
	check("failed reading resp body", err)

	return body
}

/*
 */
func makeHTTPReq(method, url string, body io.Reader) *http.Request {
	req, err := http.NewRequest(strings.ToUpper(method), url, body)
	check("error creating request:", err)

	return req
}

/*
 */
func makeHTTPReqWithAuthToken(method, url string, body io.Reader, authToken string) *http.Request {
	req := makeHTTPReq(method, url, body)

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))

	return req
}

/*
The client credentials grant flow is meant only for server-to-server API requests that use an app access token.
*/
func getOAuthTokenUsingClientCredentialsGrantFlow() OAuthToken {
	//fmt.Print("Attempting to get OAuth Token....\n")

	url := "https://id.twitch.tv/oauth2/token"

	params := fmt.Sprintf("client_id=%s&client_secret=%s&grant_type=client_credentials",
		os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"))

	req := makeHTTPReq("POST", url, bytes.NewBufferString(params))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	body := fetch(req)

	var oAuthToken OAuthToken
	check("unable to unmarshal oauthtoken body", json.Unmarshal(body, &oAuthToken))

	return oAuthToken
}

/*
Used for most api endpoints
*/
func FetchOAuthEndpoint(queries string) []byte {
	//fmt.Print("Attempting to Fetch OAuth Endpoint....\n")

	url := fmt.Sprintf("https://api.twitch.tv/helix/%s", queries)

	oAuthToken := getOAuthTokenUsingClientCredentialsGrantFlow()

	req := makeHTTPReqWithAuthToken("GET", url, nil, oAuthToken.AccessToken)
	req.Header.Add("Client-Id", os.Getenv("CLIENT_ID"))

	body := fetch(req)

	return body
}
