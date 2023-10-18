package meta

import (
	"net/http"
)

type Route struct {
	Method      string
	Name        string
	Path        string
	Version     string
	Handler     func(http.ResponseWriter, *http.Request)
	Middlewares []string
}

type Payload struct {
	Success bool        `json:"success,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   *Error       `json:"error,omitempty"`
}

type Error struct {
	ErrorCode    int    `json:"error_code,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}

type CreateParameters struct {
	Identifier string `json:"identifier"`
	Version    int    `json:"version"`
	Path       string `json:"path"`
	Data       string `json:"data"`
}

type UpdateParameters struct {
	ID      int
	Version int
	Path    string
	Data    string
}

type DeleteParameters struct {
	ID      int
	Version int
	Path    string
}

type GetBagParameters struct {
	ID int
}
