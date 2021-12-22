package gotube

import (
	"io"
	"net/http"
	"os"
	"reflect"
	"sort"
)

// Holds information about a single stream/format.
type SStream struct {
	URL        string  `json:"url"`         //
	Extension  string  `json:"ext"`         //
	FileSize   int     `json:"filesize"`    //
	ASR        float32 `json:"asr"`         // Audio sample rate.
	TBR        float32 `json:"tbr"`         //
	VBR        float32 `json:"vbr"`         // Video bit rate.
	Quality    int     `json:"quality"`     //
	ACodec     string  `json:"acodec"`      // Audio codec.
	VCodec     string  `json:"vcodec"`      // Video Codec.
	Height     int     `json:"height"`      //
	Width      int     `json:"width"`       //
	FPS        float32 `json:"fps"`         //
	FormatID   string  `json:"format_id"`   //
	FormatNote string  `json:"format_note"` //
}

type SStreamSlice []*SStream

// Returns true if the stream is a dash stream.
func (Me *SStream) IsDash() bool {
	return (Me.ACodec == "none") != (Me.VCodec == "none")
}

// Returns true if the stream contains audio data.
func (Me *SStream) HasAudio() bool {
	return Me.ACodec != "none"
}

// Returns true if the stream contains video data.
func (Me *SStream) HasVideo() bool {
	return Me.VCodec != "none"
}

// Downloads the stream file, extension will be added automatically. Will overwrite file if one is already existing!
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

// Returns a filtered stream list based on 'InFilterPredicate'. Does not modify the original list.
func (Me SStreamSlice) GetFiltered(InFilterPredicate func(InStream *SStream) bool) SStreamSlice {
	OutStreams := SStreamSlice{}
	for _, ThisStream := range Me {
		if InFilterPredicate(ThisStream) {
			OutStreams = append(OutStreams, ThisStream)
		}
	}
	return OutStreams
}

// Returns an ordered stream list based on 'InProperty'. Does not modify the original list. Panics if the propertry does not exist.
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
