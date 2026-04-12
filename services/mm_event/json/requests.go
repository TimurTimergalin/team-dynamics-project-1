package json

type RequestType string

const (
	Register   RequestType = "Register"
	Unregister RequestType = "Unregister"
)

type Request struct {
	Type RequestType `json:"type"`
}
