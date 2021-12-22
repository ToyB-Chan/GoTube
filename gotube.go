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

type SGoTube struct {
	YTDLPath        string
	CustomArguments []string
}

func Init(InYouTubeDLPath string, InCustomArguments []string) *SGoTube {
	return &SGoTube{
		YTDLPath:        InYouTubeDLPath,
		CustomArguments: InCustomArguments,
	}
}

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

func (Me *SGoTube) NewPlaylist(InURL string) (*SPlaylist, error) {
	if !strings.Contains(InURL, "playlist?list=") {
		return nil, errors.New("this is a video, use NewVideo() instead")
	}

	Command := exec.Command(Me.YTDLPath, InURL, "-J", "-s", "-4", "--no-check-certificate", "--flat-playlist")
	Command.Args = append(Command.Args, Me.CustomArguments...)
	Command.Stdin = nil
	//Command.Stdout = nil
	//Command.Stderr = os.Stdout
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
		go parallelFetcher(i, "https://youtube.com/watch?v="+ThisVideo.ID)
	}

	WaitGroup.Wait()
	return &OutPlaylist, nil
}
