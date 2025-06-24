package httputil

import (
	"encoding/json"
	"net/http"
)

func ParseJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	//dec.DisallowUnknownFields()

	return dec.Decode(dst)
}
