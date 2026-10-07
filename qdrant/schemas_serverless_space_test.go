package qdrant

import (
	"sort"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	commonv1 "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/common/v1"
	spacev1 "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/serverless/space/v1"
)

func TestFlattenServerlessSpace(t *testing.T) {
	created := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	space := &spacev1.Space{
		Id:                  "space-id",
		AccountId:           "account-id",
		Name:                "my-space",
		CloudRegionId:       "eu-central",
		CreatedAt:           timestamppb.New(created),
		Labels:              []*commonv1.KeyValue{{Key: "env", Value: "dev"}},
		CostAllocationLabel: newPointer("env"),
		Configuration: &spacev1.SpaceConfiguration{
			AllowedIpSourceRanges:  []string{"10.0.0.0/8"},
			AllowedOrigins:         []string{"https://app.example.com"},
			MaxCollectionsPerSpace: 100,
			CollectionSettings:     &spacev1.CollectionSettings{PlatformMaxSize: 2048, MaxSize: newPointer(uint64(1024))},
			SearcherSettings: &spacev1.SearcherSettings{
				IdleTimeout:        durationpb.New(10 * time.Minute),
				PlatformMaxWorkers: 8,
				MaxWorkers:         newPointer(uint64(4)),
			},
		},
		State: &spacev1.SpaceState{
			Phase:    spacev1.SpaceStatePhase_SPACE_STATE_PHASE_READY,
			Endpoint: &spacev1.SpaceEndpoint{Url: "https://space.example.com", RestPort: 443, GrpcPort: 6334},
		},
	}

	flat := flattenServerlessSpace(space)

	assert.Equal(t, "space-id", flat[slSpaceIDFieldName])
	assert.Equal(t, "account-id", flat[slSpaceAccountIDFieldName])
	assert.Equal(t, "my-space", flat[slSpaceNameFieldName])
	assert.Equal(t, "eu-central", flat[slSpaceCloudRegionIDFieldName])
	assert.Equal(t, "env", flat[slSpaceCostAllocationLabelFieldName])
	assert.Equal(t, formatTime(timestamppb.New(created)), flat[slSpaceCreatedAtFieldName])
	assert.Equal(t, "https://space.example.com", flat[slSpaceURLFieldName])
	assert.Equal(t, 443, flat[slSpaceRestPortFieldName])
	assert.Equal(t, 6334, flat[slSpaceGrpcPortFieldName])
	assert.Equal(t, []interface{}{map[string]interface{}{
		slStatePhaseFieldName:  "SPACE_STATE_PHASE_READY",
		slStateReasonFieldName: "",
	}}, flat[slSpaceStateFieldName])

	cfg := flat[slSpaceConfigurationFieldName].([]interface{})[0].(map[string]interface{})
	assert.Equal(t, []string{"10.0.0.0/8"}, cfg[slSpaceCfgAllowedIPSourceRangesFieldName])
	assert.Equal(t, []string{"https://app.example.com"}, cfg[slSpaceCfgAllowedOriginsFieldName])
	assert.Equal(t, 100, cfg[slSpaceCfgMaxCollectionsPerSpaceFieldName])
	cs := cfg[slSpaceCfgCollectionSettingsFieldName].([]interface{})[0].(map[string]interface{})
	assert.Equal(t, 2048, cs[slCollectionSettingsPlatformMaxSizeFieldName])
	assert.Equal(t, 1024, cs[slCollectionSettingsMaxSizeFieldName])
	ss := cfg[slSpaceCfgSearcherSettingsFieldName].([]interface{})[0].(map[string]interface{})
	assert.Equal(t, "10m0s", ss[slSearcherSettingsIdleTimeoutFieldName])
	assert.Equal(t, 8, ss[slSearcherSettingsPlatformMaxWorkersFieldName])
	assert.Equal(t, 4, ss[slSearcherSettingsMaxWorkersFieldName])

	t.Run("when the state is unspecified, it never emits UNSPECIFIED", func(t *testing.T) {
		flat := flattenServerlessSpace(&spacev1.Space{Id: "x"})
		state := flat[slSpaceStateFieldName].([]interface{})[0].(map[string]interface{})
		assert.Empty(t, state[slStatePhaseFieldName])
		assert.Equal(t, []interface{}{}, flat[slSpaceConfigurationFieldName])
	})

	t.Run("when flattened values are set on the resource, it is accepted by the schema", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, serverlessSpaceSchema(false), map[string]interface{}{})
		require.NoError(t, setFlattened(d, flat))
		assert.Equal(t, "my-space", d.Get(slSpaceNameFieldName))
	})
}

func TestExpandServerlessSpace(t *testing.T) {
	d := schema.TestResourceDataRaw(t, serverlessSpaceSchema(false), map[string]interface{}{
		slSpaceNameFieldName:                "my-space",
		slSpaceCloudRegionIDFieldName:       "eu-central",
		slSpaceCostAllocationLabelFieldName: "env",
		slSpaceLabelsFieldName: []interface{}{
			map[string]interface{}{"key": "env", "value": "dev"},
		},
		slSpaceConfigurationFieldName: []interface{}{map[string]interface{}{
			slSpaceCfgAllowedIPSourceRangesFieldName: []interface{}{"10.0.0.0/8"},
			slSpaceCfgAllowedOriginsFieldName:        []interface{}{"https://app.example.com"},
			slSpaceCfgCollectionSettingsFieldName: []interface{}{map[string]interface{}{
				slCollectionSettingsMaxSizeFieldName: 1024,
			}},
			slSpaceCfgSearcherSettingsFieldName: []interface{}{map[string]interface{}{
				slSearcherSettingsIdleTimeoutFieldName: "10m",
				slSearcherSettingsMaxWorkersFieldName:  4,
			}},
		}},
	})

	space := expandServerlessSpace(d, "account-id")

	assert.Equal(t, "account-id", space.GetAccountId())
	assert.Equal(t, "my-space", space.GetName())
	assert.Equal(t, "eu-central", space.GetCloudRegionId())
	assert.Equal(t, "env", space.GetCostAllocationLabel())
	require.Len(t, space.GetLabels(), 1)
	assert.Equal(t, "env", space.GetLabels()[0].GetKey())
	assert.Equal(t, "dev", space.GetLabels()[0].GetValue())
	cfg := space.GetConfiguration()
	assert.Equal(t, []string{"10.0.0.0/8"}, cfg.GetAllowedIpSourceRanges())
	assert.Equal(t, []string{"https://app.example.com"}, cfg.GetAllowedOrigins())
	assert.Equal(t, uint64(1024), cfg.GetCollectionSettings().GetMaxSize())
	assert.Equal(t, 10*time.Minute, cfg.GetSearcherSettings().GetIdleTimeout().AsDuration())
	assert.Equal(t, uint64(4), cfg.GetSearcherSettings().GetMaxWorkers())

	t.Run("when optional fields are not set, it leaves them unset", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, serverlessSpaceSchema(false), map[string]interface{}{
			slSpaceNameFieldName:          "my-space",
			slSpaceCloudRegionIDFieldName: "eu-central",
		})
		space := expandServerlessSpace(d, "account-id")
		assert.Nil(t, space.CostAllocationLabel)
		assert.Empty(t, space.GetLabels())
		assert.Nil(t, space.GetConfiguration().GetCollectionSettings())
		assert.Nil(t, space.GetConfiguration().GetSearcherSettings())
	})
}

func TestServerlessSpaceUpdateMask(t *testing.T) {
	t.Run("when nothing changed, it returns no paths", func(t *testing.T) {
		assert.Empty(t, serverlessSpaceUpdateMask(func(string) bool { return false }))
	})
	t.Run("when everything changed, it returns all updatable paths", func(t *testing.T) {
		paths := serverlessSpaceUpdateMask(func(string) bool { return true })
		sort.Strings(paths)
		assert.Equal(t, []string{
			"configuration.allowed_ip_source_ranges",
			"configuration.allowed_origins",
			"configuration.collection_settings.max_size",
			"configuration.searcher_settings.idle_timeout",
			"configuration.searcher_settings.max_workers",
			"cost_allocation_label",
			"labels",
			"name",
		}, paths)
	})
	t.Run("when a nested attribute changed, it maps to the nested proto path", func(t *testing.T) {
		paths := serverlessSpaceUpdateMask(func(attr string) bool {
			return attr == "configuration.0.searcher_settings.0.max_workers"
		})
		assert.Equal(t, []string{"configuration.searcher_settings.max_workers"}, paths)
	})
	t.Run("when the mask attributes are used, they exist in the resource schema", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, serverlessSpaceSchema(false), map[string]interface{}{})
		for _, p := range serverlessSpaceUpdateMaskPaths {
			assert.NotPanics(t, func() { d.HasChange(p.attribute) }, p.attribute)
			assert.NotNil(t, d.Get(p.attribute), p.attribute)
		}
	})
}
