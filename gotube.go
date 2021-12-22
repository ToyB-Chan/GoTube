package gotube

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// Use this structs methods to request videos/playlists.
type SGoTube struct {
	YTDLPath        string
	CustomArguments []string
}

// Returns a new 'SGoTube' struct. You can set the path to 'youtube-dl' here as well as set custom arguments. Use this structs methods to request videos/playlists.
func Init(InYouTubeDLPath string, InCustomArguments []string) *SGoTube {
	return &SGoTube{
		YTDLPath:        InYouTubeDLPath,
		CustomArguments: InCustomArguments,
	}
}

// Returns a new 'SVideo' struct of given video URL.
func (Me *SGoTube) NewVideo(InURL string) (*SVideo, error) {
	if strings.Contains(InURL, "playlist?list=") {
		return nil, errors.New("this is a playlist, use NewPlaylist() instead")
	}

	Command := exec.Command(Me.YTDLPath, InURL, "-J", "-s", "-4", "--no-check-certificate")
	Command.Args = append(Command.Args, Me.CustomArguments...)
	Command.Stdin = nil
	Command.Stdout = nil
	JSONBytes, err := Command.Output()

	if err != nil {
		return nil, err
	}

	OutVideo := SVideo{}
	json.Unmarshal(JSONBytes, &OutVideo)
	json.Unmarshal(JSONBytes, &OutVideo.Uploader)

	ConvertedTime, _ := time.Parse("20060102", OutVideo.UploadDate)
	OutVideo.UploadDate = fmt.Sprint(ConvertedTime.Unix())
	return &OutVideo, nil
}

// Returns a new 'SPlaylist' struct of given playlist URL.
func (Me *SGoTube) NewPlaylist(InURL string, InExtractParallel bool) (*SPlaylist, error) {
	if !strings.Contains(InURL, "playlist?list=") {
		return nil, errors.New("this is a video, use NewVideo() instead")
	}

	Command := exec.Command(Me.YTDLPath, InURL, "-J", "-s", "-4", "--no-check-certificate", "--flat-playlist")
	Command.Args = append(Command.Args, Me.CustomArguments...)
	Command.Stdin = nil
	Command.Stdout = nil
	JSONBytes, err := Command.Output()

	if err != nil {
		return nil, err
	}

	OutPlaylist := SPlaylist{}

	json.Unmarshal(JSONBytes, &OutPlaylist)
	json.Unmarshal(JSONBytes, &OutPlaylist.Uploader)

	WaitGroup := sync.WaitGroup{}
	parallelFetcher := func(InIndex int, InURL string) {
		WaitGroup.Add(1)
		OutPlaylist.Videos[InIndex], _ = Me.NewVideo(InURL)
		WaitGroup.Done()
	}

	for i, ThisVideo := range OutPlaylist.Videos {
		if InExtractParallel {
			go parallelFetcher(i, "https://youtube.com/watch?v="+ThisVideo.ID)
		} else {
			parallelFetcher(i, "https://youtube.com/watch?v="+ThisVideo.ID)
		}
	}

	WaitGroup.Wait()
	return &OutPlaylist, nil
}
