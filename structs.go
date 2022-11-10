package main

type Response struct {
	Continue struct {
		Plcontinue string `json:"plcontinue"`
		Continue   string `json:"continue"`
	}
	Query struct {
		Pages struct {
			int struct {
				PageId int    `json:"pageId"`
				Ns     int    `json:"ns"`
				Title  string `json:"title"`
				Links  []Link `json:"links"`
			}
		}
	}
	Limits struct {
		Links int `json:"links"`
	}
}

type Link struct {
	Ns    int    `json:"ns"`
	Title string `json:"title"`
}
