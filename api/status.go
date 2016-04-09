package api

import "github.com/robzienert/http-healthcheck"

// HealthStatusResponse is returned by the health status endpoint.
type HealthStatusResponse struct {
	Status map[string]string `json:"status"`
}

// MarshalHealthStatusResponse converts a health.Status model to a
// HealthStatusResponse object.
func MarshalHealthStatusResponse(status healthcheck.Status) HealthStatusResponse {
	r := HealthStatusResponse{Status: make(map[string]string, len(status.Statuses))}
	for p, s := range status.Statuses {
		var v string
		if s == nil {
			v = "OK"
		} else {
			v = s.Error()
		}
		r.Status[p] = v
	}
	return r
}
