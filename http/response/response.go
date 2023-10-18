package response

import (
	"encoding/json"
	"fmt"
	chttp "github.com/eidentitet/rest-go/http/meta"
	"io"
	"net/http"
)

type Response struct {
	Payload chttp.Payload
	Writer  http.ResponseWriter
}

// Json Marks a response as a json response
func (response *Response) Json() *Response {
	response.Writer.Header().Set("Content-Type", "application/json")
	return response
}

// Send writes the data to the response
func (response *Response) Send(statusCode int) {
	response.Payload.Success = true
	response.Writer.WriteHeader(statusCode)
	json.NewEncoder(response.Writer).Encode(response.Payload)
}

// SetPayload updates the response.Payload with any data
func (response *Response) SetPayload(data any) *Response {
	response.Payload.Data = data
	return response
}

// SetError sets the error code and message. It is used before
// calling Error(statusCode)
func (response *Response) SetError(code int, message string) *Response {
	response.Payload.Error = &chttp.Error{
		ErrorCode:    code,
		ErrorMessage: message,
	}
	return response
}

// Error just like Send it writes the error to the response
func (response *Response) Error(statusCode int) {
	response.Payload.Success = false
	response.Writer.WriteHeader(statusCode)
	json.NewEncoder(response.Writer).Encode(response.Payload)
}

// ParseBody reads the raw body data and parses it into the parameters reference
func ParseBody(body io.ReadCloser, parameters interface{}) error {

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(body).Decode(&parameters)
	if err != nil {
		return fmt.Errorf("Failed to parse json body: %s", err.Error())
	}

	return nil
}
