package common

import "errors"

var (
	ERR_CONNECTION_LOSS   = errors.New("ERR_CONNECTION_LESS")
	ERR_SEND_MESSAGE_FULL = errors.New("ERR_SEND_MESSAGE_FULL")
	ERR_ROOM_ID_INVALID   = errors.New("ERR_ROOM_ID_INVALID")
	ERR_MAX_ROOM          = errors.New("ERR_MAX_ROOM")
	ERR_JOINED_ROOM       = errors.New("ERR_JOINED_ROOM")
	ERR_LEAVE_ROOM        = errors.New("ERR_LEAVE_ROOM")
)
