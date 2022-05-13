package gotube

import (
	"io"
	"net/http"
	"os"
	"reflect"
	"sort"
)

// This holds all important (and not so important) information about a video stream.
type Stream struct {
	URL        string  `json:"url"`         // The URL of the stream.
	Extension  string  `json:"ext"`         // The file extension of the stream.
	FileSize   int     `json:"filesize"`    // The file size of the stream.
	ASR        float32 `json:"asr"`         // The audio sample rate used.
	TBR        float32 `json:"tbr"`         // I don't know. But it's there and sure useful for someone.
	VBR        float32 `json:"vbr"`         // The video bit rate used.
	Quality    int     `json:"quality"`     // The quality of the stream. This corresponds to the quality setting when watching youtube videos.
	ACodec     string  `json:"acodec"`      // The audio codec used.
	VCodec     string  `json:"vcodec"`      // The video codec used.
	Height     int     `json:"height"`      // The height of the video.
	Width      int     `json:"width"`       // The width of the video.
	FPS        float32 `json:"fps"`         // The frames per second of the video.
	FormatID   string  `json:"format_id"`   // The Format ID of the stream. See https://gist.github.com/AgentOak/34d47c65b1d28829bb17c24c04a0096f for more.
	FormatNote string  `json:"format_note"` // The format note of the stream.
}

type Streams []*Stream

// Returns true if the stream is a dash stream.
func (s *Stream) IsDash() bool {
	return (s.ACodec == "none") != (s.VCodec == "none")
}

// Returns true if the stream contains audio data.
func (s *Stream) HasAudio() bool {
	return s.ACodec != "none"
}

// Returns true if the stream contains video data.
func (s *Stream) HasVideo() bool {
	return s.VCodec != "none"
}

// Downloads the stream file to the specified Path\File. File extension is automatically added. File is overwritten if it already exists. Returns any errors encountered.
func (s *Stream) Download(destPath string, destFile string) error {
	if destPath == "" {
		destPath = "."
	}
	if destFile == "" {
		destFile = "download"
	}

	combinedPath := destPath + "\\" + destFile + "." + s.Extension

	file, err := os.Create(combinedPath)
	if err != nil {
		return err
	}
	defer file.Close()

	resp, err := http.Get(s.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

// Returns a list of streams that match the specified filter function. It does not modify the original list.
func (s Streams) Filtered(predicate func(s *Stream) bool) Streams {
	streams := Streams{}
	for _, stream := range s {
		if predicate(stream) {
			streams = append(streams, stream)
		}
	}
	return streams
}

// Returns a list of streams sorted by the specified property. It does not modify the original list. It panics if the property is not found.
func (s Streams) OrderedBy(property string) Streams {
	streams := Streams{}
	streams = append(streams, s...)

	if !reflect.ValueOf(*streams[0]).FieldByName(property).IsValid() {
		panic("property '" + property + "' does not exist")
	}

	isFloat := reflect.ValueOf(*streams[0]).FieldByName(property).Kind() == reflect.Float32
	isInt := reflect.ValueOf(*streams[0]).FieldByName(property).Kind() == reflect.Int

	sort.Slice(streams, func(i, j int) bool {
		refElemI := reflect.ValueOf(*streams[i])
		refElemJ := reflect.ValueOf(*streams[j])

		if isFloat {
			return refElemI.FieldByName(property).Float() < refElemJ.FieldByName(property).Float()
		} else if isInt {
			return refElemI.FieldByName(property).Int() < refElemJ.FieldByName(property).Int()
		} else {
			panic("property '" + property + "' is not a int/float value")
		}
	})
	return streams
}
