package outputs

import (
	"github.com/animenotifier/arn"
	"github.com/animenotifier/avatar"
	"github.com/nfnt/resize"
)

// WebPFile ...
type WebPFile struct {
	Directory string
	Size      int
	Quality   float32
}

// SaveAvatar writes the avatar in WebP format to the file system.
func (output *WebPFile) SaveAvatar(avatar *avatar.Avatar) error {
	img := avatar.Image

	// Resize if needed
	if img.Bounds().Dx() > output.Size {
		img = resize.Resize(uint(output.Size), 0, img, resize.Lanczos3)
	}

	// Write to file
	fileName := output.Directory + avatar.User.ID + ".webp"
	return arn.SaveWebP(img, fileName, output.Quality)
}
