package sources

import (
	"errors"
	"time"

	"github.com/animenotifier/anilist"
	"github.com/animenotifier/arn"
	"github.com/animenotifier/avatar"
)

// AniList - https://anilist.co/
type AniList struct {
	RequestLimiter *time.Ticker
}

// GetAvatar returns the Gravatar image for a user (if available).
func (source *AniList) GetAvatar(user *arn.User) (*avatar.Avatar, error) {
	alNick := user.Accounts.AniList.Nick

	// If the user has no username we can't get an avatar.
	if alNick == "" {
		return nil, errors.New("No Anilist nick")
	}

	// Wait for request limiter to allow us to send a request
	<-source.RequestLimiter.C

	// Authorize
	anilist.Authorize()

	// Get anilist user
	anilistUser, err := anilist.GetUser(alNick)

	if err != nil {
		return nil, err
	}

	// Download
	return avatar.FromURL(anilistUser.ImageURLLarge, user)
}
