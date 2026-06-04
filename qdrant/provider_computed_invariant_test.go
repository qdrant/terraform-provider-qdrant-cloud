package qdrant

import (
	"sort"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

// computedInvariantAllowlist lists Optional-not-Computed fields that are
// intentionally exempt from the provider-wide invariant below. Each entry must
// be either genuinely client-side (the backend never echoes it) or pending
// confirmation of backend behavior. Do NOT add a field here to silence the test
// without understanding whether the backend can populate it — that is exactly
// the perpetual-diff bug (CP-552). Prefer Computed: true.
const reasonUnconfirmedBackend = "unconfirmed: verify backend default behavior"

var computedInvariantAllowlist = map[string]string{
	// Top-level fields in other resources whose backend round-trip behavior is
	// not yet confirmed. If the API returns a value when the user leaves them
	// unset, they must become Computed (see CP-552 follow-up). Listed here so the
	// invariant covers the whole provider without silently flipping behavior.
	"qdrant-cloud_accounts_cluster.labels":                   "unconfirmed: may be server-augmented; verify before flipping to Computed",
	"qdrant-cloud_accounts_backup_schedule.retention_period": reasonUnconfirmedBackend,
	"qdrant-cloud_accounts_manual_backup.retention_period":   reasonUnconfirmedBackend,
	"qdrant-cloud_accounts_role.description":                 reasonUnconfirmedBackend,
}

// TestProviderOptionalConfigFieldsAreComputed is the provider-wide generalization
// of the cluster-only TestOptionalFieldsMustBeComputed (#186) and the env-only
// TestHCEnvOptionalFieldsMustBeComputed. With a gRPC/protobuf backend, any field
// the server can populate must be Computed or it produces a perpetual diff
// (CP-552). This walks every registered resource and asserts that invariant for
// scalar and set fields the backend reflects, with principled exemptions:
//   - Required / Default / Deprecated fields are not candidates.
//   - Unbounded TypeList / TypeMap fields are ordered/keyed user input, not a
//     Computed concern (and Computed on them is rarely correct).
//   - Fields *inside* a TypeSet element are governed by the set's hash func, not
//     Computed, so we do not descend into TypeSet elements.
//
// New server-reflected fields added to any resource are now caught automatically.
func TestProviderOptionalConfigFieldsAreComputed(t *testing.T) {
	p := Provider()
	var offenders []string
	for name, res := range p.ResourcesMap {
		walkComputedInvariant(name, res.Schema, &offenders)
	}
	sort.Strings(offenders)
	assert.Empty(t, offenders,
		"these Optional fields are not Computed; if the backend can populate them they will perpetually diff (CP-552). "+
			"Set Computed: true, or add to computedInvariantAllowlist with justification:\n%v", offenders)
}

func walkComputedInvariant(path string, m map[string]*schema.Schema, offenders *[]string) {
	for k, s := range m {
		fp := path + "." + k

		// Only scalar and set fields are Computed candidates here.
		checkable := s.Type == schema.TypeString || s.Type == schema.TypeBool ||
			s.Type == schema.TypeInt || s.Type == schema.TypeFloat || s.Type == schema.TypeSet

		if checkable && s.Optional && !s.Computed && !s.Required && s.Default == nil && s.Deprecated == "" {
			if _, ok := computedInvariantAllowlist[fp]; !ok {
				*offenders = append(*offenders, fp)
			}
		}

		// Descend into nested config objects (TypeList blocks) only. Do NOT descend
		// into TypeSet elements: their inner fields are keyed by the set hash func,
		// not by Computed.
		if s.Type == schema.TypeList {
			if res, ok := s.Elem.(*schema.Resource); ok {
				walkComputedInvariant(fp, res.Schema, offenders)
			}
		}
	}
}
