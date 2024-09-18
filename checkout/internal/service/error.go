package service

import "errors"

var ErrMultipleCartsReceived = errors.New("multiple carts received")
var ErrEmptyCart = errors.New("empty cart")
