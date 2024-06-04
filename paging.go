package linkedin

type Paging struct {
	Count int    `json:"count"`
	Start int    `json:"start"`
	Links []Link `json:"links"`
	//Total int    `json:"total"`
}

type MetaData struct {
	NextPageToken string `json:"nextPageToken"`
}

type Link struct {
	Type string `json:"type"`
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

func (paging Paging) HasLink(rel string) bool {
	for _, link := range paging.Links {
		if link.Rel == rel {
			return true
		}
	}

	return false
}
