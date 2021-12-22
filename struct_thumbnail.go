package gotube

import (
	"reflect"
	"sort"
)

// This holds all information about a thumbnail.
type SThumbnail struct {
	ID     int    `json:"id"`     // The thumbnail ID.
	URL    string `json:"url"`    // The thumbnail URL.
	Height int    `json:"height"` // The thumbnail height.
	Width  int    `json:"width"`  // The thumbnail width.
}

type SThumbnailSlice []*SThumbnail

// GetFiltered returns a list of thumbnails that match the specified filter function. It does not modify the original list.
func (Me SThumbnailSlice) GetFiltered(InFilterPredicate func(InThumbnail *SThumbnail) bool) SThumbnailSlice {
	OutThumbnails := SThumbnailSlice{}
	for _, ThisThumbnail := range Me {
		if InFilterPredicate(ThisThumbnail) {
			OutThumbnails = append(OutThumbnails, ThisThumbnail)
		}
	}
	return OutThumbnails
}

// GetOrderBy returns a list of thumbnails sorted by the specified property. It does not modify the original list. It panics if the property is not found.
func (Me SThumbnailSlice) GetOrderedBy(InProperty string) SThumbnailSlice {
	OutThumbnails := SThumbnailSlice{}
	OutThumbnails = append(OutThumbnails, Me...)

	if !reflect.ValueOf(*OutThumbnails[0]).FieldByName(InProperty).IsValid() {
		panic("property '" + InProperty + "' does not exist")
	}

	sort.Slice(OutThumbnails, func(i, j int) bool {
		RefElemI := reflect.ValueOf(*OutThumbnails[i])
		RefElemJ := reflect.ValueOf(*OutThumbnails[j])

		return RefElemI.FieldByName(InProperty).Int() < RefElemJ.FieldByName(InProperty).Int()
	})
	return OutThumbnails
}
