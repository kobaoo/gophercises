package models

type Story map[string]Chapter

type Options struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

type Chapter struct {
	Title      string    `json:"title"`
	Paragraphs []string  `json:"story"`
	Options    []Options `json:"options"`
}
