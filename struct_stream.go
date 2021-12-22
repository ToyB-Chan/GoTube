package gotube

import (
	"io"
	"net/http"
	"os"
	"reflect"
	"sort"
)

// This holds all important (and not so important) information about a video stream.
type SStream struct {
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

type SStreamSlice []*SStream

// IsDash returns true if the stream is a dash stream.
func (Me *SStream) IsDash() bool {
	return (Me.ACodec == "none") != (Me.VCodec == "none")
}

// HasAudio returns true if the stream contains audio data.
func (Me *SStream) HasAudio() bool {
	return Me.ACodec != "none"
}

// HasVideo returns true if the stream contains video data.
func (Me *SStream) HasVideo() bool {
	return Me.VCodec != "none"
}

// Download downloads the stream to the specified Path\File. File extension is automatically added. File is overwritten if it already exists. It returns any errors encounterd.
func (Me *SStream) Download(InDestPath string, InDestFile string) error {
	if InDestPath == "" {
		InDestPath = "."
	}
	if InDestFile == "" {
		InDestFile = "download"
	}

	CombinedPath := InDestPath + "\\" + InDestFile + "." + Me.Extension

	File, err := os.Create(CombinedPath)
	if err != nil {
		return err
	}
	defer File.Close()

	Resp, err := http.Get(Me.URL)
	if err != nil {
		return err
	}
	defer Resp.Body.Close()

	_, err = io.Copy(File, Resp.Body)
	return err
}

// GetFiltered returns a list of streams that match the specified filter function. It does not modify the original list.
func (Me SStreamSlice) GetFiltered(InFilterPredicate func(InStream *SStream) bool) SStreamSlice {
	OutStreams := SStreamSlice{}
	for _, ThisStream := range Me {
		if InFilterPredicate(ThisStream) {
			OutStreams = append(OutStreams, ThisStream)
		}
	}
	return OutStreams
}

// GetOrderedBy returns a list of streams sorted by the specified property. It does not modify the original list. It panics if the property is not found.
func (Me SStreamSlice) GetOrderedBy(InProperty string) SStreamSlice {
	OutStreams := SStreamSlice{}
	OutStreams = append(OutStreams, Me...)

	if !reflect.ValueOf(*OutStreams[0]).FieldByName(InProperty).IsValid() {
		panic("property '" + InProperty + "' does not exist")
	}

	IsPropertyFloat := reflect.ValueOf(*OutStreams[0]).FieldByName(InProperty).Kind() == reflect.Float32
	sort.Slice(OutStreams, func(i, j int) bool {
		RefElemI := reflect.ValueOf(*OutStreams[i])
		RefElemJ := reflect.ValueOf(*OutStreams[j])

		if IsPropertyFloat {
			return RefElemI.FieldByName(InProperty).Float() < RefElemJ.FieldByName(InProperty).Float()
		} else {
			return RefElemI.FieldByName(InProperty).Int() < RefElemJ.FieldByName(InProperty).Int()
		}
	})
	return OutStreams
}
