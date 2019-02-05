package main

type Story struct {
	ID      int         `json:"id"`
	Author  string      `json:"author"`
	Created string      `json:"created"`
	Updated string      `json:"updated"`
	Title   string      `json:"title"`
	Teaser  interface{} `json:"teaser"`
	Body    interface{} `json:"body"`
	URL     string      `json:"url"`
	Image   []struct {
		Size string `json:"size"`
		URL  string `json:"url"`
	} `json:"image"`
	Channels []struct {
		Name string `json:"name"`
	} `json:"channels"`
	Stocks []struct {
		Name string `json:"name"`
	} `json:"stocks"`
	Tags []struct {
		Name string `json:"name"`
	} `json:"tags"`
}
