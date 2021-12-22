package gotube

type SChannel struct {
	ID   string `json:"uploader_id"`
	URL  string `json:"uploader_url"`
	Name string `json:"uploader"`
}
