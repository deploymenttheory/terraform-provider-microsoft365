package utilityGuidListSharder

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// setComputedState sets all computed attributes on the state model
// This is the single point where state values are assigned
func setComputedState(ctx context.Context, state *GuidListSharderDataSourceModel, shards [][]string, resourceType string, shardCount int, strategy string) error {
	// Convert shards to Terraform Map of Sets format
	shardsMap := make(map[string]types.Set, len(shards))
	for i, shard := range shards {
		shardSet, diags := types.SetValueFrom(ctx, types.StringType, shard)
		if diags.HasError() {
			return fmt.Errorf("failed to convert shard %d to set: %v", i, diags.Errors())
		}
		shardsMap[fmt.Sprintf("shard_%d", i)] = shardSet
	}

	shardsMapValue, diags := types.MapValueFrom(ctx, types.SetType{ElemType: types.StringType}, shardsMap)
	if diags.HasError() {
		return fmt.Errorf("failed to convert shards map to state: %v", diags.Errors())
	}

	// Set all computed attributes
	state.Id = types.StringValue(fmt.Sprintf("%s-%d-%s", resourceType, shardCount, strategy))
	state.Shards = shardsMapValue

	return nil
}
