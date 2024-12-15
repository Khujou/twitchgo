package streamers

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/khujou/twitchgo/api/oauth"
)

type UserClips struct {
	Clips      []Clip `json:"data"`
	Pagination struct {
		Cursor string `json:"cursor"`
	} `json:"pagination"`
}

type Clip struct {
	ID              string  `json:"id"`
	URL             string  `json:"url"`
	EmbedURL        string  `json:"embed_url"`
	BroadcasterID   string  `json:"broadcaster_id"`
	BroadcasterName string  `json:"broadcaster_name"`
	CreatorID       string  `json:"creator_id"`
	CreatorName     string  `json:"creator_name"`
	VideoID         string  `json:"video_id"`
	GameID          string  `json:"game_id"`
	Language        string  `json:"language"`
	Title           string  `json:"title"`
	ViewCount       int     `json:"view_count"`
	CreationDate    string  `json:"created_at"`
	ThumbnailURL    string  `json:"thumbnail_url"`
	ClipDuration    float32 `json:"duration"`
	VodOffset       int32   `json:"vod_offset"`
	IsFeatured      bool    `json:"is_featured"`
}

/*
Fetches clips
*/
func fetchClips(broadcasterName string, broadcasterID string, st time.Time, et time.Time) []Clip { // Get top 10 clips from a channel from the previous month

	fmt.Printf("Attempting to fetch clips of %s\n", broadcasterName)

	startTimeParsed := st.Format(time.RFC3339)
	endTimeParsed := et.Format(time.RFC3339)

	params := url.Values{}
	params.Add("broadcaster_id", broadcasterID)
	params.Add("started_at", startTimeParsed)
	params.Add("ended_at", endTimeParsed)
	params.Add("first", "10")

	queries := fmt.Sprintf("clips?%s", params.Encode())

	body := oauth.FetchOAuthEndpoint(queries)

	var userClips UserClips
	check(json.Unmarshal(body, &userClips))

	return userClips.Clips
}

func GetClips(user User, startTime time.Time, endTime time.Time) []Clip {
	clips := fetchClips(user.DisplayName, user.ID, startTime, endTime)
	return clips
}
