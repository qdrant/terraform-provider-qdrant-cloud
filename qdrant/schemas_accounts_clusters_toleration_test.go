package qdrant

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
