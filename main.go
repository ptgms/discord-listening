package main

import (
	"discordlistening/helpers"
	"fmt"
	"github.com/hugolgst/rich-go/client"
	"github.com/shkh/lastfm-go/lastfm"
	"strconv"
	"time"
)

func main() {
	// Check if the configuration file exists in order to create en empty one
	if !helpers.FileExists("config.yaml") {
		helpers.MakeEmptyFile()
	}

	// Load config yaml content into a struct readable by this program
	config := helpers.LoadSettings()

	api := lastfm.New(config.LastFM.ApiKey, config.LastFM.ApiKeySec) // create an api session
	// test call to check creds
	_, err := api.User.GetRecentTracks(lastfm.P{
		"user": config.LastFM.Username,
	})

	if err != nil {
		panic(err)
	}

	errD := client.Login(config.Discord.Client) // create the discord session
	if errD != nil {
		panic(err)
	}

	var current = "" // placeholder values for the checker if new song started
	var currentTime = time.Now()
	var finishTime = time.Now().Add(time.Duration(10000))

	for {
		song, err := api.User.GetRecentTracks(lastfm.P{ // refresh status
			"user": config.LastFM.Username,
		})

		if err != nil {
			continue
		}

		if song.Tracks[0].NowPlaying == "true" { // if latest track reported as now playing
			duration, _ := api.Track.GetInfo(lastfm.P{
				"artist": song.Tracks[0].Artist.Name, "track": song.Tracks[0].Name,
			})
			if current != song.Tracks[0].Name {
				currentTime = time.Now()
				i, errI := strconv.Atoi(duration.Duration)
				if errI == nil && i != 0 { // end time fetch failed, will display elapsed in discord
					finishTime = time.Now().Add(time.Millisecond * time.Duration(i))
				} else {
					finishTime = time.Unix(1, 1)
				}
				current = song.Tracks[0].Name
			}

			if finishTime == time.Unix(1, 1) {
				_ = client.SetActivity(client.Activity{
					Details:    song.Tracks[0].Name,
					State:      fmt.Sprintf("Playing %s by %s", song.Tracks[0].Name, song.Tracks[0].Artist.Name),
					LargeImage: config.Discord.ImageName,
					LargeText:  "Scrobbling now on last.fm",
					Timestamps: &client.Timestamps{Start: &currentTime},
				})
			} else {
				_ = client.SetActivity(client.Activity{
					Details:    song.Tracks[0].Name,
					State:      fmt.Sprintf("Playing %s by %s", song.Tracks[0].Name, song.Tracks[0].Artist.Name),
					LargeImage: config.Discord.ImageName,
					LargeText:  "Scrobbling now on last.fm",
					Timestamps: &client.Timestamps{
						Start: &currentTime,
						End:   &finishTime,
					},
				})
			}
		} else {
			_ = client.SetActivity(client.Activity{
				Details: "None",
			})
		}
		time.Sleep(time.Second * 1)
	}
}
