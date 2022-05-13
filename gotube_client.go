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

// Use this structs methods to extract video/playlist data.
type GoTubeClient struct {
	YoutubeDlPath   string
	CustomArguments []string
}

// This is a constructor for 'GoTubeClient'. Here you can set the path to the youtube-dl binary and custom arguments. Use this structs methods to extract video/playlist data. Returns a new 'SGoTube' struct.
func New(youtubeDlPath string, customArgs []string) *GoTubeClient {
	return &GoTubeClient{
		YoutubeDlPath:   youtubeDlPath,
		CustomArguments: customArgs,
	}
}

// Extracts video information from given video URL. Returns a new 'Video' struct and any errors encountered.
func (gtc *GoTubeClient) NewVideo(url string) (*Video, error) {
	if strings.Contains(url, "playlist?list=") {
		return nil, errors.New("this is a playlist, use NewPlaylist() instead")
	}

	cmd := exec.Command(gtc.YoutubeDlPath, url, "-J", "-s", "-4", "--no-check-certificate")
	cmd.Args = append(cmd.Args, gtc.CustomArguments...)
	cmd.Stdin = nil
	cmd.Stdout = nil
	jsonBytes, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	video := Video{}
	json.Unmarshal(jsonBytes, &video)
	json.Unmarshal(jsonBytes, &video.Uploader)

	ConvertedTime, _ := time.Parse("20060102", video.UploadDate)
	video.UploadDate = fmt.Sprint(ConvertedTime.Unix())
	return &video, nil
}

// Extracts playlist information from given playlist URL. Returns a new 'Playlist' struct and any errors encountered.
func (gtc *GoTubeClient) NewPlaylist(url string, extractParallel bool) (*Playlist, error) {
	if !strings.Contains(url, "playlist?list=") {
		return nil, errors.New("this is a video, use NewVideo() instead")
	}

	cmd := exec.Command(gtc.YoutubeDlPath, url, "-J", "-s", "-4", "--no-check-certificate", "--flat-playlist")
	cmd.Args = append(cmd.Args, gtc.CustomArguments...)
	cmd.Stdin = nil
	cmd.Stdout = nil
	jsonBytes, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	playlist := Playlist{}

	json.Unmarshal(jsonBytes, &playlist)
	json.Unmarshal(jsonBytes, &playlist.Uploader)

	wg := sync.WaitGroup{}
	parallelFetcher := func(InIndex int, InURL string) {
		wg.Add(1)
		playlist.Videos[InIndex], _ = gtc.NewVideo(InURL)
		wg.Done()
	}

	for i, vid := range playlist.Videos {
		if extractParallel {
			go parallelFetcher(i, "https://youtube.com/watch?v="+vid.ID)
		} else {
			parallelFetcher(i, "https://youtube.com/watch?v="+vid.ID)
		}
	}

	wg.Wait()
	return &playlist, nil
}
