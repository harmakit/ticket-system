package message

import "errors"

var ErrUnknownMessageType = errors.New("unknown message type")
var ErrMessageTypeHeaderNotFound = errors.New("message type header not found")
