package utilityGuidListSharder

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// setStateToTerraform sets all values in the state object
func setStateToTerraform(ctx context.Context, state *GuidListSharderDataSourceModel, shards [][]string, resourceType string, shardCount int, strategy string) error {

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

	state.Id = types.StringValue(uuid.New().String())
	state.Shards = shardsMapValue

	return nil
}
