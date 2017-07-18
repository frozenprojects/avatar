package sources

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/animenotifier/arn"
	"github.com/animenotifier/avatar"
	"github.com/parnurzeal/gorequest"
)

var userIDRegex = regexp.MustCompile(`<user_id>(\d+)<\/user_id>`)

// MyAnimeList - https://myanimelist.net/
type MyAnimeList struct {
	RequestLimiter *time.Ticker
}

// GetAvatar returns the Gravatar image for a user (if available).
func (source *MyAnimeList) GetAvatar(user *arn.User) (*avatar.Avatar, error) {
	malNick := user.Accounts.MyAnimeList.Nick

	// If the user has no username we can't get an avatar.
	if malNick == "" {
		return nil, errors.New("No MAL nick")
	}

	// Download user info
	userInfoURL := "https://myanimelist.net/malappinfo.php?u=" + malNick
	response, xml, networkErrs := gorequest.New().Get(userInfoURL).End()

	if len(networkErrs) > 0 {
		return nil, networkErrs[0]
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("Unexpected status code: " + strconv.Itoa(response.StatusCode))
	}

	// Build URL
	matches := userIDRegex.FindStringSubmatch(xml)

	if matches == nil || len(matches) < 2 {
		return nil, errors.New("Could not find MAL user ID")
	}

	malID := matches[1]
	malAvatarURL := "https://myanimelist.cdn-dena.com/images/userimages/" + malID + ".jpg"

	// Wait for request limiter to allow us to send a request
	<-source.RequestLimiter.C

	// Download
	return avatar.FromURL(malAvatarURL, user)
}
