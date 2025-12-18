package planmodifiers

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestUseStateForUnknownInt32(t *testing.T) {
	tests := []struct {
		name         string
		planValue    types.Int32
		stateValue   types.Int32
		expectedPlan types.Int32
	}{
		{
			name:         "unknown plan value with known state",
			planValue:    types.Int32Unknown(),
			stateValue:   types.Int32Value(42),
			expectedPlan: types.Int32Value(42),
		},
		{
			name:         "known plan value",
			planValue:    types.Int32Value(10),
			stateValue:   types.Int32Value(42),
			expectedPlan: types.Int32Value(10),
		},
		{
			name:         "unknown plan value with null state",
			planValue:    types.Int32Unknown(),
			stateValue:   types.Int32Null(),
			expectedPlan: types.Int32Unknown(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := planmodifier.Int32Request{
				Path:       path.Root("test"),
				PlanValue:  tt.planValue,
				StateValue: tt.stateValue,
			}
			resp := &planmodifier.Int32Response{
				PlanValue: tt.planValue,
			}

			modifier := UseStateForUnknownInt32()
			modifier.PlanModifyInt32(context.Background(), req, resp)

			assert.Equal(t, tt.expectedPlan, resp.PlanValue)
		})
	}
}

func TestInt32DefaultValue(t *testing.T) {
	tests := []struct {
		name         string
		planValue    types.Int32
		defaultValue int32
		expectedPlan types.Int32
	}{
		{
			name:         "null plan value",
			planValue:    types.Int32Null(),
			defaultValue: 100,
			expectedPlan: types.Int32Value(100),
		},
		{
			name:         "known plan value",
			planValue:    types.Int32Value(50),
			defaultValue: 100,
			expectedPlan: types.Int32Value(50),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := planmodifier.Int32Request{
				Path:       path.Root("test"),
				PlanValue:  tt.planValue,
				StateValue: types.Int32Null(),
			}
			resp := &planmodifier.Int32Response{
				PlanValue: tt.planValue,
			}

			modifier := Int32DefaultValue(tt.defaultValue)
			modifier.PlanModifyInt32(context.Background(), req, resp)

			assert.Equal(t, tt.expectedPlan, resp.PlanValue)
		})
	}
}
