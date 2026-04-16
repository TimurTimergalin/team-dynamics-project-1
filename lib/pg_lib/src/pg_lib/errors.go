package pg_lib

import "fmt"

type PgLibErrorType = int32

const (
	LogicError PgLibErrorType = iota
	ConnectionError
	ServerError
)

func getStringRepr(type_ PgLibErrorType) string {
	switch type_ {
	case ConnectionError:
		return "Connection error"
	case LogicError:
		return "Logic error"
	case ServerError:
		return "Server error"
	}
	return "Unknown error"
}

type PgLibError struct {
	Type  PgLibErrorType
	Cause error
}

func (p *PgLibError) Error() string {
	return fmt.Sprintf("pglib %s: %v", getStringRepr(p.Type), p.Cause)
}

func makeConnectionError(cause error) *PgLibError {
	return &PgLibError{ConnectionError, cause}
}

func makeLogicError(cause error) *PgLibError {
	return &PgLibError{LogicError, cause}
}

func makeServerError(cause error) *PgLibError {
	return &PgLibError{ServerError, cause}
}
