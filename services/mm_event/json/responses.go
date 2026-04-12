package json

type ResponseType string

const (
	Registered   ResponseType = "Registered"
	Unregistered ResponseType = "Unregistered"
	Match        ResponseType = "Match"
	Error        ResponseType = "Error"
)

type Response struct {
	Type         ResponseType `json:"type"`
	Address      *string      `json:"address,omitempty"`
	ErrorMessage *string      `json:"errorMessage,omitempty"`
}
