package heroku

import (
	"fmt"
	"net/http"
)

// ListRange describes a range.
type ListRange struct {
	Field      string
	Max        int
	Descending bool
	FirstID    string
	LastID     string
}

// SetHeader set headers on the given Request.
func (lr *ListRange) SetHeader(req *http.Request) {
	var hdrval string
	if lr.Field != "" {
		hdrval += lr.Field + " "
	}
	hdrval += lr.FirstID + ".." + lr.LastID
	if lr.Max != 0 {
		hdrval += fmt.Sprintf("; max=%d", lr.Max)
		if lr.Descending {
			hdrval += ", "
		}
	}
	if lr.Descending {
		hdrval += ", order=desc"
	}
	req.Header.Set("Range", hdrval)
	return
}
