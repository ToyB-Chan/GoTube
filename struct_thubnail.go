package gotube

import (
	"reflect"
	"sort"
)

type SThumbnail struct {
	ID     int    `json:"id"`
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type SThumbnailSlice []*SThumbnail

func (Me SThumbnailSlice) GetFiltered(InFilterPredicate func(InThumbnail *SThumbnail) bool) SThumbnailSlice {
	OutThumbnails := SThumbnailSlice{}
	for _, ThisThumbnail := range Me {
		if InFilterPredicate(ThisThumbnail) {
			OutThumbnails = append(OutThumbnails, ThisThumbnail)
		}
	}
	return OutThumbnails
}

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
