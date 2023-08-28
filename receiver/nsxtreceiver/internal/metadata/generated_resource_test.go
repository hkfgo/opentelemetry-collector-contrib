// Code generated by mdatagen. DO NOT EDIT.

package metadata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResourceBuilder(t *testing.T) {
	for _, test := range []string{"default", "all_set", "none_set"} {
		t.Run(test, func(t *testing.T) {
			cfg := loadResourceAttributesConfig(t, test)
			rb := NewResourceBuilder(cfg)
			rb.SetDeviceID("device.id-val")
			rb.SetNsxtNodeID("nsxt.node.id-val")
			rb.SetNsxtNodeName("nsxt.node.name-val")
			rb.SetNsxtNodeType("nsxt.node.type-val")

			res := rb.Emit()
			assert.Equal(t, 0, rb.Emit().Attributes().Len()) // Second call should return empty Resource

			switch test {
			case "default":
				assert.Equal(t, 4, res.Attributes().Len())
			case "all_set":
				assert.Equal(t, 4, res.Attributes().Len())
			case "none_set":
				assert.Equal(t, 0, res.Attributes().Len())
				return
			default:
				assert.Failf(t, "unexpected test case: %s", test)
			}

			val, ok := res.Attributes().Get("device.id")
			assert.True(t, ok)
			if ok {
				assert.EqualValues(t, "device.id-val", val.Str())
			}
			val, ok = res.Attributes().Get("nsxt.node.id")
			assert.True(t, ok)
			if ok {
				assert.EqualValues(t, "nsxt.node.id-val", val.Str())
			}
			val, ok = res.Attributes().Get("nsxt.node.name")
			assert.True(t, ok)
			if ok {
				assert.EqualValues(t, "nsxt.node.name-val", val.Str())
			}
			val, ok = res.Attributes().Get("nsxt.node.type")
			assert.True(t, ok)
			if ok {
				assert.EqualValues(t, "nsxt.node.type-val", val.Str())
			}
		})
	}
}
