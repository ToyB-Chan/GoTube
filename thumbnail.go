package gotube

import (
	"reflect"
	"sort"
)

// This holds all information about a thumbnail.
type Thumbnail struct {
	ID     int    `json:"id"`     // The thumbnail ID.
	URL    string `json:"url"`    // The thumbnail URL.
	Height int    `json:"height"` // The thumbnail height.
	Width  int    `json:"width"`  // The thumbnail width.
}

type Thumbnails []*Thumbnail

// Returns a list of thumbnails that match the specified filter function. It does not modify the original list.
func (t Thumbnails) Filtered(predicate func(t *Thumbnail) bool) Thumbnails {
	thumbnails := Thumbnails{}
	for _, ThisThumbnail := range t {
		if predicate(ThisThumbnail) {
			thumbnails = append(thumbnails, ThisThumbnail)
		}
	}
	return thumbnails
}

// Returns a list of thumbnails sorted by the specified property. It does not modify the original list. It panics if the property is not found.
func (t Thumbnails) OrderedBy(property string) Thumbnails {
	thumbnails := Thumbnails{}
	thumbnails = append(thumbnails, t...)

	if !reflect.ValueOf(*thumbnails[0]).FieldByName(property).IsValid() {
		panic("property '" + property + "' does not exist")
	}

	isInt := reflect.ValueOf(*thumbnails[0]).FieldByName(property).Kind() == reflect.Int

	sort.Slice(thumbnails, func(i, j int) bool {
		refElemI := reflect.ValueOf(*thumbnails[i])
		refElemJ := reflect.ValueOf(*thumbnails[j])

		if isInt {
			return refElemI.FieldByName(property).Int() < refElemJ.FieldByName(property).Int()
		} else {
			panic("property '" + property + "' is not a int value")
		}
	})
	return thumbnails
}
