package avatar

import (
	"bytes"
	"errors"
	"image"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/animenotifier/arn"
	"github.com/parnurzeal/gorequest"
)

// ImageFromURL ...
func ImageFromURL(url string) (img image.Image, data []byte, format string, err error) {
	// Download
	response, data, networkErrs := gorequest.New().Get(url).EndBytes()

	// Network errors
	if len(networkErrs) > 0 {
		return nil, nil, "", networkErrs[0]
	}

	// Retry HTTP only version after 5 seconds if service unavailable
	if response == nil || response.StatusCode == http.StatusServiceUnavailable {
		time.Sleep(5 * time.Second)
		response, data, networkErrs = gorequest.New().Get(strings.Replace(url, "https://", "http://", 1)).EndBytes()
	}

	// Network errors on 2nd try
	if len(networkErrs) > 0 {
		return nil, nil, "", networkErrs[0]
	}

	// Bad status codes
	if response.StatusCode != http.StatusOK {
		return nil, nil, "", errors.New("Unexpected status code: " + strconv.Itoa(response.StatusCode))
	}

	// Decode
	img, format, decodeErr := image.Decode(bytes.NewReader(data))

	if decodeErr != nil {
		return nil, nil, "", decodeErr
	}

	return img, data, format, nil
}

// FromURL downloads and decodes the image from an URL and creates an Avatar.
func FromURL(url string, user *arn.User) (*Avatar, error) {
	img, data, format, err := ImageFromURL(url)

	if err != nil {
		return nil, err
	}

	return &Avatar{
		User:   user,
		Image:  img,
		Data:   data,
		Format: format,
	}, nil
}
