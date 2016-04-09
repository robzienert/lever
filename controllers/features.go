package controllers

import (
	"errors"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/robzienert/gin-middleware/correlationid"
	"github.com/robzienert/lever/api"
	"github.com/robzienert/lever/metrics"
	"github.com/robzienert/lever/model"
	"github.com/robzienert/lever/processor"
	"github.com/robzienert/lever/router/middleware/session"
	"github.com/robzienert/lever/shared/strutil"
	"github.com/robzienert/lever/store"
)

const (
	namespaceQuery = "ns"
	actorsQuery    = "actors"
	groupsQuery    = "groups"
	keyParam       = "key"
)

// GetAllFeatures returns all features and their gate settings.
//
// Allows passing the "ns" query param to get all features just by namespace.
// Not passing the ns query param will return only those features without a
// namespace.
func GetAllFeatures(c *gin.Context) {
	features, err := store.GetFeatureList(c, c.Query(namespaceQuery))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if features == nil {
		features = make([]*model.Feature, 0)
	}

	c.IndentedJSON(http.StatusOK, api.GetFeatureListResponse{Features: features})
}

// GetFeature returns an individual feature by namespace or not.
func GetFeature(c *gin.Context) {
	feature, err := store.GetFeature(c, c.Query(namespaceQuery), c.Param(keyParam))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if feature == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.IndentedJSON(http.StatusOK, api.FeatureResponse{Feature: feature})
}

// PutFeature idemopotently upserts a feature and its associated gates.
func PutFeature(c *gin.Context) {
	var in model.Feature
	if err := c.BindJSON(&in); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if c.Param(keyParam) != in.Key {
		c.AbortWithError(http.StatusBadRequest, errors.New("key URL param does not match Feature key in body"))
		return
	}

	feature, err := store.GetFeature(c, in.Namespace, in.Key)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var breadcrumb *model.Breadcrumb
	now := time.Now().UTC()
	if feature == nil {
		feature = &in
		feature.DateCreated = now

		breadcrumb = model.NewBreadcrumb("create feature", session.AuditActor(c))
	} else {
		diff := feature.Diff(&in)

		feature.Type = in.Type
		feature.Value = in.Value
		feature.Gate = in.Gate

		breadcrumb = model.NewBreadcrumb("update feature", session.AuditActor(c)).WithFields(diff)
	}
	feature.LastUpdated = now

	if err = store.UpsertFeature(c, feature); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	go saveBreadcrumb(c.Copy(), breadcrumb)

	c.IndentedJSON(http.StatusOK, api.FeatureResponse{Feature: feature})
}

// DeleteFeature removes a feature from the service.
func DeleteFeature(c *gin.Context) {
	feature, err := store.GetFeature(c, c.Query(namespaceQuery), c.Param(keyParam))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if feature == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if err = store.DeleteFeature(c, feature); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	breadcrumb := model.NewBreadcrumb("delete feature", session.AuditActor(c)).WithFields(model.Fields{
		"key": feature.Key,
		"ns":  feature.Namespace,
	})

	go saveBreadcrumb(c.Copy(), breadcrumb)

	c.Writer.WriteHeader(http.StatusNoContent)
}

// GetFeatureState will retrieve the current gate state of a feature, given
// optional actors and groups.
func GetFeatureState(c *gin.Context) {
	metrics.WithTiming(c, "featureState.single", func() {
		feature, err := store.GetFeature(c, c.Query(namespaceQuery), c.Param(keyParam))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if feature == nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		enabled, err := processor.ProcessGate(feature.Gate, c.Query(actorsQuery), c.Query(groupsQuery))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		state := api.FeatureState{Key: feature.Key, Enabled: enabled}
		if feature.Namespace != "" {
			state.Namespace = feature.Namespace
		}

		c.IndentedJSON(http.StatusOK, state)
	})
}

// PostBatchFeatureState will return the current gate state of a group of features
// given optional actors and groups.
//
// TODO For now just iterating serially over features. This could easily be more
// efficient with a pipeline.
func PostBatchFeatureState(c *gin.Context) {
	metrics.WithTiming(c, "featureState.batch", func() {
		var in api.BatchFeatureStateRequest
		if err := c.BindJSON(&in); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		if len(in.NamespacedFeatures) == 0 && len(in.Features) == 0 {
			c.AbortWithError(http.StatusBadRequest, errors.New("no features to process"))
			return
		}

		resp := api.BatchFeatureState{}

		var all []*model.Feature
		for ns, requestedFeatures := range in.NamespacedFeatures {
			features, err := store.GetFeatureList(c, ns)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			for _, f := range features {
				if strutil.StringInSlice(f.Key, requestedFeatures) {
					all = append(all, f)
				}
			}
		}
		if len(in.Features) > 0 {
			features, err := store.GetFeatureList(c, "")
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			for _, f := range features {
				if strutil.StringInSlice(f.Key, in.Features) {
					all = append(all, f)
				}
			}
		}

		for _, f := range all {
			enabled, err := processor.ProcessGate(f.Gate, c.Query(actorsQuery), c.Query(groupsQuery))
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			state := api.FeatureState{Key: f.Key, Enabled: enabled}
			if f.Namespace != "" {
				state.Namespace = f.Namespace
			}
			resp.States = append(resp.States, state)
		}
		c.IndentedJSON(http.StatusOK, resp)
	})
}

func saveBreadcrumb(c *gin.Context, breadcrumb *model.Breadcrumb) {
	if err := store.SaveBreadcrumb(c, breadcrumb); err != nil {
		correlationid.Logger(c).WithFields(logrus.Fields{
			"err": err,
			"b":   *breadcrumb,
		}).Error("Could not save breadcrumb")
	}
}
