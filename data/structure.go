package data

type AudioData struct {
	Name  string `json:"filename"`
	Title string `json:"title"`
}

type Playlist struct {
	ID          string      `json:"id"`
	Album       string      `json:"album"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	AudioList   []AudioData `json:"audioList"`
}
