package main

import (
	"discordlistening/helpers"
	"fmt"
	"github.com/gen2brain/beeep"
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

	debugPrint(fmt.Sprintf("Service to sync discord to LastFM started successfully for user %s!", config.LastFM.Username), config)
	if !config.ServiceSettings.BlockNoti {
		errNotify := beeep.Notify(fmt.Sprintf("LastFM sync started for user %s!", config.LastFM.Username), "Service is running!", "")
		if errNotify != nil {
			debugPrint(err.Error(), config)
		}
	}

	var current = "" // placeholder values for the checker if new song started
	var currentTime = time.Now()
	var finishTime = time.Now().Add(time.Duration(10000))

	var lastStatus = false
	for {
		debugPrint("Starting loop", config)

		song, err := api.User.GetRecentTracks(lastfm.P{ // refresh status
			"user": config.LastFM.Username,
		})

		if err != nil {
			debugPrint("Check-in with LastFM failed, retrying next loop!", config)
			continue
		}

		if song.Tracks[0].NowPlaying == "true" { // if latest track reported as now playing
			debugPrint("Latest track is now-playing!", config)
			if !lastStatus {
				errD := client.Login(config.Discord.Client) // create the discord session
				if errD != nil {
					panic(err)
				}
			}
			duration, _ := api.Track.GetInfo(lastfm.P{
				"artist": song.Tracks[0].Artist.Name, "track": song.Tracks[0].Name,
			})
			if current != song.Tracks[0].Name {
				debugPrint(fmt.Sprintf("Song changed to %s", song.Tracks[0].Name), config)
				currentTime = time.Now()
				i, errI := strconv.Atoi(duration.Duration)
				if errI == nil && i != 0 { // end time fetch failed, will display elapsed in discord
					finishTime = time.Now().Add(time.Millisecond * time.Duration(i))
				} else {
					debugPrint(fmt.Sprintf("No time duration for %s", song.Tracks[0].Name), config)
					finishTime = time.Unix(1, 1)
				}

				var timestamp = &client.Timestamps{
					Start: nil,
					End:   nil,
				}

				if config.ServiceSettings.NoTime {
					timestamp = nil
				} else if finishTime == time.Unix(1, 1) {
					debugPrint("Song didn't have duration in LastFM database, so I will only display elapsed time.", config)
					timestamp = &client.Timestamps{Start: &currentTime}
				} else {
					debugPrint(fmt.Sprintf("Start time is %s and stop time is %s.", currentTime.String(), finishTime.String()), config)
					timestamp = &client.Timestamps{
						Start: &currentTime,
						End:   &finishTime,
					}
				}
				_ = client.SetActivity(client.Activity{
					Details:    song.Tracks[0].Name,
					State:      fmt.Sprintf("Playing %s by %s", song.Tracks[0].Name, song.Tracks[0].Artist.Name),
					LargeImage: config.Discord.ImageName,
					LargeText:  "Scrobbling now on last.fm",
					Timestamps: timestamp,
				})
				current = song.Tracks[0].Name
			}
			lastStatus = true
		} else {
			debugPrint("None playing on lastfm right now, displaying none.", config)
			client.Logout()
			lastStatus = false
		}
		time.Sleep(time.Second * 1)
	}
}

func debugPrint(message string, config helpers.Config) {
	if config.ServiceSettings.DebugMode {
		fmt.Printf("[D] %s\n", message)
	}
}
