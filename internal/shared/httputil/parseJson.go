package httputil

import (
	"encoding/json"
	"net/http"
)

func ParseJson(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	//dec.DisallowUnknownFields()

	return dec.Decode(dst)
}
