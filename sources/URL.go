package sources

import (
	"errors"

	"github.com/animenotifier/arn"
	"github.com/animenotifier/avatar"
)

// URL ...
type URL struct{}

// GetAvatar returns the Gravatar image for a user (if available).
func (source *URL) GetAvatar(user *arn.User) (*avatar.Avatar, error) {
	// If the user has no Email registered we can't get a Gravatar.
	if user.Settings().Avatar.SourceURL == "" {
		return nil, errors.New("Avatar source URL is empty")
	}

	// Download
	return avatar.FromURL(user.Settings().Avatar.SourceURL, user)
}
