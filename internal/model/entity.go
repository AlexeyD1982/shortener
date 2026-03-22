package model

type URL struct {
	Origin string
	Short  string
}

type ApiShortenRequest struct {
	URL string `json:"url"`
}
type ApiShortenResponse struct {
	Result string `json:"result"`
}
