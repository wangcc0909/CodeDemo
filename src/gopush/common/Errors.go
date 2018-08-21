package common

import "errors"

var (
	ERR_CONNECTION_LOSS = errors.New("ERR_CONNECTION_LESS")
	ERR_SEND_MESSAGE_FULL= errors.New("ERR_SEND_MESSAGE_FULL")
)
