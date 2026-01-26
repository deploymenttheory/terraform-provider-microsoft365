package utilityGuidListSharder

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

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

	// Generate deterministic ID based on configuration
	// This ensures the datasource ID remains stable across refreshes
	idString := fmt.Sprintf("%s-%s-%d-%s-%s-%s",
		state.ResourceType.ValueString(),
		state.ODataFilter.ValueString(),
		shardCount,
		strategy,
		state.Seed.ValueString(),
		state.GroupId.ValueString(),
	)
	hash := sha256.Sum256([]byte(idString))
	state.Id = types.StringValue(hex.EncodeToString(hash[:]))
	state.Shards = shardsMapValue

	return nil
}
