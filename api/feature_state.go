package api

// NamespacedFeatures represents a list of namespaced features. The key of the
// map being the namespace, and the string slice a list of feature keys belonging
// to that namespace.
type NamespacedFeatures map[string][]string

// FeatureState represents an individual feature's gate state.
type FeatureState struct {
	Namespace string `json:"namespace,omitempty"`
	Key       string `json:"key"`
	Enabled   bool   `json:"enabled"`
}

// BatchFeatureState presents a collection of FeatureStates.
type BatchFeatureState struct {
	States []FeatureState `json:"states"`
}

// BatchFeatureStateRequest is used to get a collection of feature states.
type BatchFeatureStateRequest struct {
	NamespacedFeatures NamespacedFeatures `json:"namespacedFeatures,omitempty"`
	Features           []string           `json:"features,omitempty"`
}
