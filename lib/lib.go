package lib

import (
	"fmt"
	"reflect"
	"time"

	"github.com/aerogo/log"
	"github.com/animenotifier/arn"
	"github.com/animenotifier/avatar"
	"github.com/animenotifier/avatar/outputs"
	"github.com/animenotifier/avatar/sources"
	"github.com/fatih/color"
)

const (
	webPQuality = 80
)

// Log is the global log
var Log = log.New()

// Define the avatar sources
var avatarSources = []avatar.Source{
	&sources.Gravatar{
		Rating:         "pg",
		RequestLimiter: time.NewTicker(100 * time.Millisecond),
	},
	&sources.MyAnimeList{
		RequestLimiter: time.NewTicker(250 * time.Millisecond),
	},
	&sources.URL{},
	&sources.FileSystem{
		Directory: "images/avatars/large/",
	},
}

// Define the avatar outputs
var avatarOutputs = []avatar.Output{
	// Original - Large
	&outputs.OriginalFile{
		Directory: "images/avatars/large/",
		Size:      arn.AvatarMaxSize,
	},

	// Original - Small
	&outputs.OriginalFile{
		Directory: "images/avatars/small/",
		Size:      arn.AvatarSmallSize,
	},

	// WebP - Large
	&outputs.WebPFile{
		Directory: "images/avatars/large/",
		Size:      arn.AvatarMaxSize,
		Quality:   webPQuality,
	},

	// WebP - Small
	&outputs.WebPFile{
		Directory: "images/avatars/small/",
		Size:      arn.AvatarSmallSize,
		Quality:   webPQuality,
	},
}

// RefreshAvatar refreshes the avatar of a single user.
func RefreshAvatar(user *arn.User) {
	preferredSource := user.Settings().Avatar.Source
	user.Avatar.Extension = ""

	for _, source := range avatarSources {
		sourceName := reflect.TypeOf(source).Elem().Name()

		// Skip this source if it's not the user preferred source,
		// however FileSystem is always allowed.
		if preferredSource != "" && preferredSource != sourceName && sourceName != "FileSystem" {
			continue
		}

		avatar, err := source.GetAvatar(user)

		if err != nil {
			Log.Error(err)
			continue
		}

		if avatar == nil {
			// fmt.Println(color.RedString("✘"), reflect.TypeOf(source).Elem().Name(), user.Nick)
			continue
		}

		// Name of source
		user.Avatar.Source = sourceName

		// Log
		fmt.Println(color.GreenString("✔"), user.Avatar.Source, "|", user.Nick, "|", avatar)

		// Avoid JPG quality loss (if it's on the file system, we don't need to write it again)
		if user.Avatar.Source == "FileSystem" {
			user.Avatar.Extension = avatar.Extension()
			break
		}

		for _, writer := range avatarOutputs {
			err := writer.SaveAvatar(avatar)

			if err != nil {
				color.Red(err.Error())
			}
		}

		break
	}

	// Since this a very long running job, refresh user data before saving it.
	avatarExt := user.Avatar.Extension
	avatarSrc := user.Avatar.Source
	user, err := arn.GetUser(user.ID)

	if err != nil {
		Log.Error("Can't refresh user info:", user.ID, user.Nick)
		return
	}

	// Save avatar data
	user.Avatar.Extension = avatarExt
	user.Avatar.Source = avatarSrc
	user.Save()
}
