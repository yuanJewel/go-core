package object

import "errors"

var (
	SelectOverOneError = errors.New("the query results in more than one content")
	AffectedRowsError  = errors.New("operation expected to affect the content does not match the actual")
	InputError         = errors.New("there is a problem with the code, the input does not meet the specifications")
)
