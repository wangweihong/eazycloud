package iharbor

type (
	SearchRequest struct {
		Query string `json:"query" form:"query"`
	}

	SearchResponse struct {
		Projects   []Project    `json:"project"`
		Repository []Repository `json:"repository"`
	}
)
