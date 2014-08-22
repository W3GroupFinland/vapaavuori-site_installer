package app_base

import (
	"encoding/json"
	"io"
)

type JSON struct {
}

func (a *JSON) NewDecoder(r io.Reader) *json.Decoder {
	return json.NewDecoder(r)
}

func (a *JSON) NewEncoder(w io.Writer) *json.Encoder {
	return json.NewEncoder(w)
}
