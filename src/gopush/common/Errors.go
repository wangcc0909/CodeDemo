package common

import "errors"

var (
	ERR_CONNECTION_LOSS             = errors.New("ERR_CONNECTION_LESS")
	ERR_SEND_MESSAGE_FULL           = errors.New("ERR_SEND_MESSAGE_FULL")
	ERR_ROOM_ID_INVALID             = errors.New("ERR_ROOM_ID_INVALID")
	ERR_MAX_ROOM                    = errors.New("ERR_MAX_ROOM")
	ERR_JOINED_ROOM                 = errors.New("ERR_JOINED_ROOM")
	ERR_LEAVE_ROOM_UNEXIST          = errors.New("ERR_LEAVE_ROOM_UNEXIST")
	ERR_JOIN_ROOM_TWICE             = errors.New("ERR_JOIN_ROOM_TWICE")
	ERR_DISPATCH_CHANNEL_FULL       = errors.New("ERR_DISPATCH_CHANNEL_FULL")
	ERR_CERT_NOT_INVALID            = errors.New("ERR_CERT_NOT_INVALID")
	ERR_SERVER_RUN_FAIL             = errors.New("ERR_SERVER_RUN_FAIL")
	ERR_MERGER_CHANNEL_FULL         = errors.New("ERR_MERGER_CHANNEL_FULL")
	ERR_LOGIC_DISPATCH_CHANNEL_FULL = errors.New("ERR_LOGIC_DISPATCH_CHANNEL_FULL")
)
