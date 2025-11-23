package pr

import "time"

type PrList struct {
	Number   int       `json:"number"`
	MergedAt time.Time `json:"mergedAt"`
}

type Author struct {
	ID    string `json:"id"`
	IsBot bool   `json:"is_bot"`
	Login string `json:"login"`
	Name  string `json:"name"`
}

type MergedPrList struct {
	Author Author `json:"author"`
	PrList
	Title string `json:"title"`
}
