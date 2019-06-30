package json

import "github.com/pquerna/ffjson/ffjson"

var DefaultMarshal = ffjson.Marshal
var Marshal = DefaultMarshal

var DefaultUnmarshal = ffjson.Unmarshal
var Unmarshal = DefaultUnmarshal
