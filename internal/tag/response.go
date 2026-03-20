package tag

import "time"

type getResponse struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type listResponse struct {
	Tags []string `json:"tags"`
}
