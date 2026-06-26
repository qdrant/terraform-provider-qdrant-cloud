package qdrant

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	qcCluster "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/cluster/v1"
)

// TestTolerationSchemaValidatesEnumNames ensures the toleration operator/effect
// fields accept the protobuf enum names and reject anything else (notably the
// Kubernetes short forms like "Exists"/"NoSchedule", which expandTolerations
// would otherwise silently drop).
func TestTolerationSchemaValidatesEnumNames(t *testing.T) {
	s := tolerationSchema(false)

	opValidate := s[tolerationOperatorFieldName].ValidateFunc
	effValidate := s[tolerationEffectFieldName].ValidateFunc
	require.NotNil(t, opValidate, "operator must have a ValidateFunc in resource mode")
	require.NotNil(t, effValidate, "effect must have a ValidateFunc in resource mode")

	// Accepted: protobuf enum names.
	_, errs := opValidate("TOLERATION_OPERATOR_EXISTS", tolerationOperatorFieldName)
	assert.Empty(t, errs)
	_, errs = effValidate("TOLERATION_EFFECT_NO_SCHEDULE", tolerationEffectFieldName)
	assert.Empty(t, errs)

	// Rejected: Kubernetes short forms (the bug this guards against).
	_, errs = opValidate("Exists", tolerationOperatorFieldName)
	assert.NotEmpty(t, errs, "k8s short form 'Exists' must be rejected")
	_, errs = effValidate("NoSchedule", tolerationEffectFieldName)
	assert.NotEmpty(t, errs, "k8s short form 'NoSchedule' must be rejected")

	// Data-source mode is computed-only and must not validate.
	ds := tolerationSchema(true)
	assert.Nil(t, ds[tolerationOperatorFieldName].ValidateFunc)
	assert.Nil(t, ds[tolerationEffectFieldName].ValidateFunc)
}

// TestExpandTolerationsOmitsEmptyKeyValue ensures empty key/value strings (which
// the SDK always materializes for set elements) are not sent as empty pointers.
// Sending an empty value pointer makes the API reject an Exists toleration with
// "value must not be set when operator is Exists".
func TestExpandTolerationsOmitsEmptyKeyValue(t *testing.T) {
	in := []interface{}{
		map[string]interface{}{
			tolerationKeyFieldName:      "example.com/node-panther-static",
			tolerationValueFieldName:    "", // not set by the user
			tolerationOperatorFieldName: "TOLERATION_OPERATOR_EXISTS",
			tolerationEffectFieldName:   "TOLERATION_EFFECT_NO_SCHEDULE",
			tolerationSecondsFieldName:  0,
		},
	}

	// A toleration with a non-empty value: the value must be preserved. This
	// guards against a future change that drops all values (not just empty
	// ones) while still passing the nil-on-empty assertion above.
	in = append(in, map[string]interface{}{
		tolerationKeyFieldName:      "example.com/node-panther-static",
		tolerationValueFieldName:    "v", // set by the user
		tolerationOperatorFieldName: "TOLERATION_OPERATOR_EQUAL",
		tolerationEffectFieldName:   "TOLERATION_EFFECT_NO_SCHEDULE",
		tolerationSecondsFieldName:  0,
	})

	out := expandTolerations(in)
	require.Len(t, out, 2)

	tol := out[0]
	assert.Nil(t, tol.Value, "empty value must not be sent (nil pointer), else API rejects Exists tolerations")
	require.NotNil(t, tol.Key)
	assert.Equal(t, "example.com/node-panther-static", tol.GetKey())
	require.NotNil(t, tol.Operator)
	assert.Equal(t, qcCluster.TolerationOperator_TOLERATION_OPERATOR_EXISTS, tol.GetOperator())
	require.NotNil(t, tol.Effect)
	assert.Equal(t, qcCluster.TolerationEffect_TOLERATION_EFFECT_NO_SCHEDULE, tol.GetEffect())

	// Positive assertion: a value the user did set must survive expansion.
	withValue := out[1]
	require.NotNil(t, withValue.Value, "a non-empty value must be sent")
	assert.Equal(t, "v", withValue.GetValue())
	require.NotNil(t, withValue.Operator)
	assert.Equal(t, qcCluster.TolerationOperator_TOLERATION_OPERATOR_EQUAL, withValue.GetOperator())
}
