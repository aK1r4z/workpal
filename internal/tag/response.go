package tag

import "time"

type tagResponse struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type getResponse struct {
	tagResponse
}

type listResponse struct {
	Tags []tagResponse `json:"tags"`
}
