package http

type (
	Response struct {
		Meta Meta        `json:"meta"`
		Data interface{} `json:"data"`
	}

	Error struct {
		Meta Meta `json:"meta"`
	}

	Meta struct {
		Message string `json:"message" binding:"required"`
		Status  int    `json:"status"`
	}

	Data struct {
		Value interface{} `json:"data"`
	}
)
