package k8s

import "fmt"

type OpsErrorType int32

const (
	UnknownError OpsErrorType = iota
	FleetFullError
	ContentionError
	ConnectionError
	UniquenessViolationError
	NotFoundError
)

type OpsError struct {
	message string
	Type    OpsErrorType
	inner   error
}

func (e OpsError) Error() string {
	if e.message == "" {
		return fmt.Sprintf("K8sOps error %d", e.Type)
	}
	return fmt.Sprintf("K8sOps error %d: %s", e.Type, e.message)
}

func (e OpsError) Unwrap() error {
	return e.inner
}

func (e OpsErrorType) Make() *OpsError {
	return &OpsError{"", e, nil}
}

func (e OpsErrorType) MakeF(template string, args ...any) *OpsError {
	return &OpsError{fmt.Sprintf(template, args...), e, nil}
}

func (e OpsErrorType) Wrap(err error) *OpsError {
	return &OpsError{fmt.Sprintf("caused by error: %v", err), e, err}
}
