package api

import "github.com/robzienert/lever/model"

// GetAuditResponse is the HTTP response wrapper for audit breadcrumbs.
type GetAuditResponse struct {
	Breadcrumbs []*model.Breadcrumb `json:"breadcrumbs"`
}
