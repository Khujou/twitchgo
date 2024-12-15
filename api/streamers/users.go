package streamers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/khujou/twitchgo/api/oauth"
)

type UsersData struct {
	Users []User `json:"data"`
}

type User struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	UserType        string `json:"type"`
	BroadcasterType string `json:"broadcaster_type"`
	Description     string `json:"description"`
	ProfileImageURL string `json:"profile_image_url"`
	OfflineImageURL string `json:"offline_image_url"`
	ViewCount       int    `json:"view_count"`
	CreationDate    string `json:"created_at"`
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

func fetchUsers(logins []string) []User {

	fmt.Printf("Attempting to fetch users %s\n", logins)

	params := url.Values{}

	for _, login := range logins {
		params.Add("login", login)
	}
	queries := fmt.Sprintf("users?%s", params.Encode())

	body := oauth.FetchOAuthEndpoint(queries)

	var data UsersData
	check(json.Unmarshal(body, &data))

	return data.Users

}

func GetUsers(logins []string) []User {
	users := fetchUsers(logins)

	return users
}
