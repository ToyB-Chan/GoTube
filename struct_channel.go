package gotube

// Holds basic information about a single channel.
type SChannel struct {
	ID   string `json:"uploader_id"`  //
	URL  string `json:"uploader_url"` //
	Name string `json:"uploader"`     //
}
