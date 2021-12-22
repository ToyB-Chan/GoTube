package gotube

import (
	"io"
	"net/http"
	"os"
	"reflect"
	"sort"
)

type SStream struct {
	URL        string  `json:"url"`
	Extension  string  `json:"ext"`
	FileSize   int     `json:"filesize"`
	ASR        float32 `json:"asr"`
	TBR        float32 `json:"tbr"`
	VBR        float32 `json:"vbr"`
	Quality    int     `json:"quality"`
	ACodec     string  `json:"acodec"`
	VCodec     string  `json:"vcodec"`
	Height     int     `json:"height"`
	Width      int     `json:"width"`
	FPS        float32 `json:"fps"`
	FormatID   string  `json:"format_id"`
	FormatNote string  `json:"format_note"`
}

type SStreamSlice []*SStream

func (Me *SStream) IsDash() bool {
	return (Me.ACodec == "none") != (Me.VCodec == "none")
}

func (Me *SStream) HasAudio() bool {
	return Me.ACodec != "none"
}

func (Me *SStream) HasVideo() bool {
	return Me.VCodec != "none"
}

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

func (Me SStreamSlice) GetFiltered(InFilterPredicate func(InStream *SStream) bool) SStreamSlice {
	OutStreams := SStreamSlice{}
	for _, ThisStream := range Me {
		if InFilterPredicate(ThisStream) {
			OutStreams = append(OutStreams, ThisStream)
		}
	}
	return OutStreams
}

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
