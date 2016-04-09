package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/robzienert/lever/api"
	"github.com/robzienert/lever/model"
	"github.com/robzienert/lever/router/middleware/context"
	"github.com/robzienert/lever/store"
	"github.com/robzienert/lever/store/memory"
	"github.com/robzienert/lever/store/mock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func jsonString(data interface{}) string {
	out, _ := json.Marshal(data)
	return string(out)
}

type FeaturesTestSuite struct {
	suite.Suite
}

func (suite *FeaturesTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
}

func (suite *FeaturesTestSuite) serveEndpoint(store store.Store, method string, endpoint string, endpointFn func(router *gin.Engine), body io.Reader) *httptest.ResponseRecorder {
	router := gin.New()
	router.Use(context.SetStore(store))
	endpointFn(router)

	req, _ := http.NewRequest(method, endpoint, body)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	return resp
}

func (suite *FeaturesTestSuite) TestFeatureGetAll_Empty() {
	memStore := memory.Load()
	resp := suite.serveEndpoint(memStore, "GET", "/features", func(router *gin.Engine) {
		router.GET("/features", GetAllFeatures)
	}, nil)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)

	expected := jsonString(api.GetFeatureListResponse{Features: []*model.Feature{}})
	assert.JSONEq(suite.T(), expected, resp.Body.String())
}

func (suite *FeaturesTestSuite) TestFeatureGetAll_StoreError() {
	mockStore := mock.LoadFeatureStore(&mock.FeatureStore{
		GetListFn: func(ns string) ([]*model.Feature, error) {
			return nil, errors.New("oh no")
		},
	})

	resp := suite.serveEndpoint(mockStore, "GET", "/features", func(router *gin.Engine) {
		router.GET("/features", GetAllFeatures)
	}, nil)

	assert.Equal(suite.T(), http.StatusInternalServerError, resp.Code)
}

func (suite *FeaturesTestSuite) TestFeatureGetAll_OK() {
	memStore := memory.Load()
	f1 := model.Feature{
		Key: "foo",
	}
	memStore.Features().Upsert(&f1)
	resp := suite.serveEndpoint(memStore, "GET", "/features", func(router *gin.Engine) {
		router.GET("/features", GetAllFeatures)
	}, nil)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)

	expected := jsonString(api.GetFeatureListResponse{Features: []*model.Feature{&f1}})
	assert.JSONEq(suite.T(), expected, resp.Body.String())
}

func (suite *FeaturesTestSuite) TestFeatureGet_NotFound() {
	memStore := memory.Load()
	resp := suite.serveEndpoint(memStore, "GET", "/features/foo", func(router *gin.Engine) {
		router.GET("/features/:key", GetFeature)
	}, nil)

	assert.Equal(suite.T(), http.StatusNotFound, resp.Code)
	assert.Empty(suite.T(), resp.Body.String())
}

func (suite *FeaturesTestSuite) TestFeatureGet_StoreError() {
	mockStore := mock.LoadFeatureStore(&mock.FeatureStore{
		GetFn: func(ns string, key string) (*model.Feature, error) {
			return nil, errors.New("oh no")
		},
	})
	resp := suite.serveEndpoint(mockStore, "GET", "/features/foo", func(router *gin.Engine) {
		router.GET("/features/:key", GetFeature)
	}, nil)

	assert.Equal(suite.T(), http.StatusInternalServerError, resp.Code)
	assert.Empty(suite.T(), resp.Body.String())
}

func (suite *FeaturesTestSuite) TestFeatureGet_OK() {
	memStore := memory.Load()
	f1 := model.Feature{
		Namespace: "",
		Key:       "foo",
	}
	memStore.Features().Upsert(&f1)
	resp := suite.serveEndpoint(memStore, "GET", "/features/foo", func(router *gin.Engine) {
		router.GET("/features/:key", GetFeature)
	}, nil)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)

	expected := jsonString(api.FeatureResponse{Feature: &f1})
	assert.JSONEq(suite.T(), expected, resp.Body.String())
}

func (suite *FeaturesTestSuite) TestFeaturePut_BadRequestData() {
	// Invalid JSON
	resp := suite.serveEndpoint(memory.Load(), "PUT", "/features/foo", func(router *gin.Engine) {
		router.PUT("/features/:key", PutFeature)
	}, strings.NewReader("bad data"))

	assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)

	// Insufficient fields
	f, err := json.Marshal(model.Feature{Key: "one"})
	assert.NoError(suite.T(), err)

	resp = suite.serveEndpoint(memory.Load(), "PUT", "/features/foo", func(router *gin.Engine) {
		router.PUT("/features/:key", PutFeature)
	}, strings.NewReader(string(f)))

	assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)
}

func (suite *FeaturesTestSuite) TestFeaturePut_NonmatchingKeys() {
	f, err := json.Marshal(model.Feature{Key: "two"})
	assert.NoError(suite.T(), err)

	resp := suite.serveEndpoint(memory.Load(), "PUT", "/features/one", func(router *gin.Engine) {
		router.PUT("/features/:key", PutFeature)
	}, strings.NewReader(string(f)))

	assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)
}

func (suite *FeaturesTestSuite) TestFeaturePut_GetStoreError() {
	mockStore := mock.LoadFeatureStore(&mock.FeatureStore{
		GetFn: func(ns string, key string) (*model.Feature, error) {
			return nil, errors.New("oh no")
		},
	})

	f, err := json.Marshal(model.Feature{
		Key:   "one",
		Type:  "java.lang.Boolean",
		Value: "true",
		Gate:  &model.Gate{},
	})
	assert.NoError(suite.T(), err)

	resp := suite.serveEndpoint(mockStore, "PUT", "/features/one", func(router *gin.Engine) {
		router.PUT("/features/:key", PutFeature)
	}, strings.NewReader(string(f)))

	assert.Equal(suite.T(), http.StatusInternalServerError, resp.Code)
}

func (suite *FeaturesTestSuite) TestFeaturePut_SaveStoreError() {
	mockStore := mock.LoadFeatureStore(&mock.FeatureStore{
		GetFn: func(ns string, key string) (*model.Feature, error) {
			return nil, nil
		},
		UpsertFn: func(feature *model.Feature) error {
			return errors.New("oh no")
		},
	})

	f, err := json.Marshal(model.Feature{
		Key:   "one",
		Type:  "java.lang.Boolean",
		Value: "true",
		Gate:  &model.Gate{},
	})
	assert.NoError(suite.T(), err)

	resp := suite.serveEndpoint(mockStore, "PUT", "/features/one", func(router *gin.Engine) {
		router.PUT("/features/:key", PutFeature)
	}, strings.NewReader(string(f)))

	assert.Equal(suite.T(), http.StatusInternalServerError, resp.Code)
}

func (suite *FeaturesTestSuite) TestFeaturePut_OK() {
	memStore := memory.Load()
	f := model.Feature{
		Key:   "one",
		Type:  "java.lang.Boolean",
		Value: "true",
		Gate: &model.Gate{
			Value: "true",
		},
	}

	featureJSON, err := json.Marshal(&f)
	assert.NoError(suite.T(), err)

	resp := suite.serveEndpoint(memStore, "PUT", "/features/one", func(router *gin.Engine) {
		router.PUT("/features/:key", PutFeature)
	}, strings.NewReader(string(featureJSON)))

	assert.Equal(suite.T(), http.StatusOK, resp.Code)

	featureResp := &api.FeatureResponse{}
	assert.NoError(suite.T(), json.Unmarshal(resp.Body.Bytes(), featureResp))
	insertFeature := featureResp.Feature
	assert.Equal(suite.T(), f.Key, insertFeature.Key, "unexpected Key")
	assert.Equal(suite.T(), f.Type, insertFeature.Type, "unexpected Type")
	assert.Equal(suite.T(), f.Value, insertFeature.Value, "unexpected Value")
	assert.NotNil(suite.T(), insertFeature.Gate, "Gate missing")
	assert.Equal(suite.T(), f.Gate.Value, insertFeature.Gate.Value, "unexpected Gate Value")
	assert.Empty(suite.T(), insertFeature.Gate.Actors, "unexpected Gate Actors")
	assert.Empty(suite.T(), insertFeature.Gate.Groups, "unexpected Gate Groups")
	assert.Zero(suite.T(), insertFeature.Gate.ActorPercent, "unexpected Gate ActorPercent")
	assert.Zero(suite.T(), insertFeature.Gate.PercentOfTime, "unexpected Gate PercentOfTime")
	assert.NotEmpty(suite.T(), insertFeature.DateCreated, "unexpected DateCreated")
	assert.Equal(suite.T(), insertFeature.DateCreated, insertFeature.LastUpdated, "DateCreated and LastUpdated do not match")

	f = model.Feature{
		Key:   "one",
		Type:  "java.lang.Boolean",
		Value: "true",
		Gate: &model.Gate{
			Value:  "",
			Actors: []string{"one", "two"},
		},
	}
	featureJSON, err = json.Marshal(&f)

	resp = suite.serveEndpoint(memStore, "PUT", "/features/one", func(router *gin.Engine) {
		router.PUT("/features/:key", PutFeature)
	}, strings.NewReader(string(featureJSON)))

	assert.Equal(suite.T(), http.StatusOK, resp.Code)

	featureResp = &api.FeatureResponse{}
	assert.NoError(suite.T(), json.Unmarshal(resp.Body.Bytes(), featureResp))
	updateFeature := featureResp.Feature
	assert.Empty(suite.T(), updateFeature.Gate.Value, "updated Gate Value not empty")
	assert.NotEmpty(suite.T(), updateFeature.Gate.Actors, "updated Gate Actors empty")
	assert.Equal(suite.T(), insertFeature.DateCreated, updateFeature.DateCreated, "updated Feature changed DateCreated")
	assert.NotEqual(suite.T(), updateFeature.DateCreated, updateFeature.LastUpdated, "updated Feature did not change LastUpdated")

	// Allow goroutines to finish.
	time.Sleep(10 * time.Millisecond)

	breadcrumbs, err := memStore.Breadcrumbs().GetList()
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), breadcrumbs, 2, "no Breadcrumbs found after destructive actions")
}

func (suite *FeaturesTestSuite) TestFeatureDelete_GetStoreError() {
	mockStore := mock.LoadFeatureStore(&mock.FeatureStore{
		GetFn: func(ns string, key string) (*model.Feature, error) {
			return nil, errors.New("oh no")
		},
	})

	resp := suite.serveEndpoint(mockStore, "DELETE", "/features/one", func(router *gin.Engine) {
		router.DELETE("/features/:key", DeleteFeature)
	}, nil)

	assert.Equal(suite.T(), http.StatusInternalServerError, resp.Code)
}

func (suite *FeaturesTestSuite) TestFeatureDelete_NotFound() {
	resp := suite.serveEndpoint(memory.Load(), "DELETE", "/features/one", func(router *gin.Engine) {
		router.DELETE("/features/:key", DeleteFeature)
	}, nil)

	assert.Equal(suite.T(), http.StatusNotFound, resp.Code)
}

func (suite *FeaturesTestSuite) TestFeatureDelete_DeleteStoreError() {
	mockStore := mock.LoadFeatureStore(&mock.FeatureStore{
		GetFn: func(ns string, key string) (*model.Feature, error) {
			return &model.Feature{}, nil
		},
		DeleteFn: func(feature *model.Feature) error {
			return errors.New("oh no")
		},
	})

	resp := suite.serveEndpoint(mockStore, "DELETE", "/features/one", func(router *gin.Engine) {
		router.DELETE("/features/:key", DeleteFeature)
	}, nil)

	assert.Equal(suite.T(), http.StatusInternalServerError, resp.Code)
}

func (suite *FeaturesTestSuite) TestFeatureDelete_OK() {
	memStore := memory.Load()

	f := model.Feature{
		Key:   "one",
		Type:  "java.lang.Boolean",
		Value: "true",
		Gate: &model.Gate{
			Value:  "",
			Actors: []string{"one", "two"},
		},
	}
	assert.NoError(suite.T(), memStore.Features().Upsert(&f))

	resp := suite.serveEndpoint(memStore, "DELETE", "/features/one", func(router *gin.Engine) {
		router.DELETE("/features/:key", DeleteFeature)
	}, nil)

	assert.Equal(suite.T(), http.StatusNoContent, resp.Code)

	features, err := memStore.Features().GetList()
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), features, 0)
}

func (suite *FeaturesTestSuite) TestFeatureSingleState_NotFound() {
	memStore := memory.Load()
	resp := suite.serveEndpoint(memStore, "GET", "/features/one/state", func(router *gin.Engine) {
		router.GET("/features/:key/state", GetFeatureState)
	}, nil)

	assert.Equal(suite.T(), http.StatusNotFound, resp.Code)
}

func (suite *FeaturesTestSuite) TestFeatureSingleState_GetStoreError() {
	mockStore := mock.LoadFeatureStore(&mock.FeatureStore{
		GetFn: func(ns string, key string) (*model.Feature, error) {
			return nil, errors.New("oh no")
		},
	})

	resp := suite.serveEndpoint(mockStore, "GET", "/features/one/state", func(router *gin.Engine) {
		router.GET("/features/:key/state", GetFeatureState)
	}, nil)

	assert.Equal(suite.T(), http.StatusInternalServerError, resp.Code)
}

func (suite *FeaturesTestSuite) TestFeatureSingleState_OK() {
	memStore := memory.Load()

	f := model.Feature{
		Key:   "one",
		Type:  "java.lang.Boolean",
		Value: "true",
		Gate: &model.Gate{
			Value: "true",
		},
	}
	assert.NoError(suite.T(), memStore.Features().Upsert(&f))

	resp := suite.serveEndpoint(memStore, "GET", "/features/one/state", func(router *gin.Engine) {
		router.GET("/features/:key/state", GetFeatureState)
	}, nil)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)

	stateResp := &api.FeatureState{}
	assert.NoError(suite.T(), json.Unmarshal(resp.Body.Bytes(), stateResp))
	assert.Equal(suite.T(), f.Key, stateResp.Key)
	assert.Equal(suite.T(), f.Namespace, stateResp.Namespace)
	assert.True(suite.T(), stateResp.Enabled)
}

func (suite *FeaturesTestSuite) TestFeatureBatchState_BadRequest() {
	resp := suite.serveEndpoint(memory.Load(), "POST", "/features", func(router *gin.Engine) {
		router.POST("/features", PutFeature)
	}, strings.NewReader("bad data"))

	assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)
}

func (suite *FeaturesTestSuite) TestFeatureBatchState_EmptyRequest() {
	req, err := json.Marshal(&api.BatchFeatureStateRequest{
		NamespacedFeatures: api.NamespacedFeatures{},
		Features:           []string{},
	})
	assert.NoError(suite.T(), err)

	resp := suite.serveEndpoint(memory.Load(), "POST", "/features", func(router *gin.Engine) {
		router.POST("/features", PutFeature)
	}, strings.NewReader(string(req)))

	assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)
}

func (suite *FeaturesTestSuite) TestFeatureBatchState_ListStoreError() {
	mockStore := mock.LoadFeatureStore(&mock.FeatureStore{
		GetListFn: func(ns string) ([]*model.Feature, error) {
			return nil, errors.New("oh no")
		},
	})

	req, err := json.Marshal(&api.BatchFeatureStateRequest{
		Features: []string{"one"},
	})
	assert.NoError(suite.T(), err)

	resp := suite.serveEndpoint(mockStore, "POST", "/features", func(router *gin.Engine) {
		router.POST("/features", PostBatchFeatureState)
	}, strings.NewReader(string(req)))

	assert.Equal(suite.T(), http.StatusInternalServerError, resp.Code)
}

func (suite *FeaturesTestSuite) TestFeatureBatchState_OK() {
	memStore := memory.Load()

	f1 := model.Feature{
		Key:   "globalKeyWithInternalGroup",
		Type:  "java.lang.Boolean",
		Value: "true",
		Gate: &model.Gate{
			Groups: []string{"internal"},
		},
	}
	assert.NoError(suite.T(), memStore.Features().Upsert(&f1))
	f2 := model.Feature{
		Namespace: "mobile.ios",
		Key:       "namespacedKeyWithFalseBoolean",
		Type:      "java.lang.Boolean",
		Value:     "true",
		Gate: &model.Gate{
			Value: "false",
		},
	}
	assert.NoError(suite.T(), memStore.Features().Upsert(&f2))
	f3 := model.Feature{
		Namespace: "mobile.ios",
		Key:       "namespacedKeyWithTrueBoolean",
		Type:      "java.lang.Boolean",
		Value:     "true",
		Gate: &model.Gate{
			Value: "true",
		},
	}
	assert.NoError(suite.T(), memStore.Features().Upsert(&f3))
	f4 := model.Feature{
		Namespace: "mobile.ios",
		Key:       "unrequestedNamespacedKeyWithTrueBoolean",
		Type:      "java.lang.Boolean",
		Value:     "true",
		Gate: &model.Gate{
			Value: "true",
		},
	}
	assert.NoError(suite.T(), memStore.Features().Upsert(&f4))
	f5 := model.Feature{
		Key:   "globalKeyWithUnrequestedActors",
		Type:  "java.lang.Boolean",
		Value: "true",
		Gate: &model.Gate{
			Actors: []string{"robzienert"},
		},
	}
	assert.NoError(suite.T(), memStore.Features().Upsert(&f5))

	req, err := json.Marshal(&api.BatchFeatureStateRequest{
		NamespacedFeatures: api.NamespacedFeatures{
			"mobile.ios": []string{"namespacedKeyWithFalseBoolean", "namespacedKeyWithTrueBoolean"},
		},
		Features: []string{"globalKeyWithInternalGroup", "globalKeyWithUnrequestedActors"},
	})
	assert.NoError(suite.T(), err)

	resp := suite.serveEndpoint(memStore, "POST", "/features?groups=internal", func(router *gin.Engine) {
		router.POST("/features", PostBatchFeatureState)
	}, strings.NewReader(string(req)))

	assert.Equal(suite.T(), http.StatusOK, resp.Code)

	stateResp := &api.BatchFeatureState{}
	assert.NoError(suite.T(), json.Unmarshal(resp.Body.Bytes(), stateResp))
	assert.Len(suite.T(), stateResp.States, 4)

	expectedStates := []api.FeatureState{
		{
			Namespace: "mobile.ios",
			Key:       "namespacedKeyWithFalseBoolean",
			Enabled:   false,
		},
		{
			Namespace: "mobile.ios",
			Key:       "namespacedKeyWithTrueBoolean",
			Enabled:   true,
		},
		{
			Key:     "globalKeyWithInternalGroup",
			Enabled: true,
		},
		{
			Key:     "globalKeyWithUnrequestedActors",
			Enabled: false,
		},
	}

	checked := 0
	for i, expected := range expectedStates {
		for _, state := range stateResp.States {
			if expected.Namespace == state.Namespace && expected.Key == state.Key {
				checked++
				if expected.Enabled {
					assert.True(suite.T(), state.Enabled, fmt.Sprintf("expected state #%d to be true", i))
				} else {
					assert.False(suite.T(), state.Enabled, fmt.Sprintf("expected state #%d to be false", i))
				}
			}
		}
	}

	assert.Equal(suite.T(), 4, checked, "unexpected state check count")

	suite.T().SkipNow()
}

func TestFeaturesTestSuite(t *testing.T) {
	suite.Run(t, new(FeaturesTestSuite))
}
