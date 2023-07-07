package models

import (
	"fmt"
)

const (
	fmtMsgErrBadRequest = "error processing request: %v"
	fmtMsgNotFound      = "requested entity %q not found"
	fmtMsgBadData       = "malformed request: %v"
)

type ErrBadData struct {
	req []byte
}

func (e ErrBadData) Error() string {
	return fmt.Sprintf(fmtMsgBadData, e.req)
}
