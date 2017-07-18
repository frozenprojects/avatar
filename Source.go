package avatar

import (
	"github.com/animenotifier/arn"
)

// Source describes a source where we can find avatar images for a user.
type Source interface {
	GetAvatar(*arn.User) *Avatar
}
