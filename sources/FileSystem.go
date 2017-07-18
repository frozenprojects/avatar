package sources

import (
	"bytes"
	"errors"
	"image"
	"io/ioutil"

	"github.com/animenotifier/arn"
	"github.com/animenotifier/avatar"
)

// FileSystem loads avatar from the local filesystem.
type FileSystem struct {
	Directory string
}

// GetAvatar returns the local image for the user.
func (source *FileSystem) GetAvatar(user *arn.User) (*avatar.Avatar, error) {
	fullPath := arn.FindFileWithExtension(user.ID, source.Directory, arn.OriginalImageExtensions)

	if fullPath == "" {
		return nil, errors.New("Not found on file system")
	}

	data, err := ioutil.ReadFile(fullPath)

	if err != nil {
		return nil, err
	}

	// Decode
	img, format, decodeErr := image.Decode(bytes.NewReader(data))

	if decodeErr != nil {
		return nil, decodeErr
	}

	return &avatar.Avatar{
		User:   user,
		Image:  img,
		Data:   data,
		Format: format,
	}, nil
}
