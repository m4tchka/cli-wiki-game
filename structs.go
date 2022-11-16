package main

type Response struct {
	Continue struct {
		Plcontinue  string `json:"plcontinue,omitempty"`
		Continue    string `json:"continue,omitempty"`
		GrnContinue string `json:"grncontinue,omitempty"`
	}
	Query struct {
		Pages map[int]Page
	}
	Limits struct {
		Links int `json:"links,omitempty"`
	}
}

type Page struct {
	PageId  int    `json:"pageId"`
	Ns      int    `json:"ns"`
	Title   string `json:"title"`
	Links   []Link `json:"links,omitempty"`
	Extract string `json:"extract,omitempty"`
}
type Link struct {
	Ns    int    `json:"ns"`
	Title string `json:"title"`
}
