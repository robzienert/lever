package api

import "github.com/robzienert/lever/model"

// GetFeatureListResponse is the HTTP response wrapper for feature lists.
type GetFeatureListResponse struct {
	Features []*model.Feature `json:"features"`
}

// FeatureResponse is the HTTP response wrapper for a single feature.
type FeatureResponse struct {
	Feature *model.Feature `json:"feature"`
}
