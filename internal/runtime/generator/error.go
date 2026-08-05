package generator

import "errors"

var (
	ErrModelNotSupportToolCall = errors.New("model does not support tool call")
)
