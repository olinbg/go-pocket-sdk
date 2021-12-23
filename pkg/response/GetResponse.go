package response

type GetResponse struct {
	Status     int                    `json:"status"`
	Complete   int                    `json:"complete"`
	List       map[string]ListElement `json:"list"`
	Error      interface{}            `json:"error"`
	SearchMeta struct {
		SearchType string `json:"search_type"`
	} `json:"search_meta"`
	Since int `json:"since"`
}

type ListElement struct {
	ItemId         string `json:"item_id"`
	ResolvedId     string `json:"resolved_id"`
	GivenUrl       string `json:"given_url"`
	GivenTitle     string `json:"given_title"`
	Favorite       string `json:"favorite"`
	Status         string `json:"status"`
	TimeAdded      string `json:"time_added"`
	TimeUpdated    string `json:"time_updated"`
	TimeRead       string `json:"time_read"`
	TimeFavorited  string `json:"time_favorited"`
	SortId         int    `json:"sort_id"`
	ResolvedTitle  string `json:"resolved_title"`
	ResolvedUrl    string `json:"resolved_url"`
	Excerpt        string `json:"excerpt"`
	IsArticle      string `json:"is_article"`
	IsIndex        string `json:"is_index"`
	HasVideo       string `json:"has_video"`
	HasImage       string `json:"has_image"`
	WordCount      string `json:"word_count"`
	Lang           string `json:"lang"`
	DomainMetadata struct {
		Name          string `json:"name"`
		Logo          string `json:"logo"`
		GreyscaleLogo string `json:"greyscale_logo"`
	} `json:"domain_metadata"`
	ListenDurationEstimate int `json:"listen_duration_estimate"`
}
