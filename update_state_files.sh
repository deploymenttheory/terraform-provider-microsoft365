#!/bin/bash

# Script to update state.go files to use the convert package instead of the state package

# Base directory
BASE_DIR="/Users/dafyddwatkins/GitHub/deploymenttheory/terraform-provider-microsoft365"

# Find all state.go files
find "$BASE_DIR/internal/services/resources" -name "state.go" | while read -r file; do
  echo "Processing $file..."
  
  # Skip files that have already been updated
  if grep -q "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert" "$file"; then
    echo "  Already updated, skipping."
    continue
  fi
  
  # Update import statements
  sed -i '' 's|"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/state"|"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"|g' "$file"
  
  # Update function calls
  sed -i '' 's|state\.StringPointerValue|convert.GraphToFrameworkString|g' "$file"
  sed -i '' 's|state\.BoolPointerValue|convert.GraphToFrameworkBool|g' "$file"
  sed -i '' 's|state\.TimeToString|convert.GraphToFrameworkTime|g' "$file"
  sed -i '' 's|state\.EnumPtrToTypeString|convert.GraphToFrameworkEnum|g' "$file"
  sed -i '' 's|state\.Int32PointerValue|convert.GraphToFrameworkInt32|g' "$file"
  sed -i '' 's|state\.Int32PointerToInt64Value|convert.GraphToFrameworkInt32AsInt64|g' "$file"
  sed -i '' 's|state\.Int64PointerValue|convert.GraphToFrameworkInt64|g' "$file"
  sed -i '' 's|state\.StringSliceToSet|convert.GraphToFrameworkStringSet|g' "$file"
  sed -i '' 's|state\.StringListToTypeList|convert.GraphToFrameworkStringList|g' "$file"
  sed -i '' 's|state\.BoolPtrToTypeBool|convert.GraphToFrameworkBool|g' "$file"
  sed -i '' 's|state\.BoolPtrToBool|convert.GraphToFrameworkBoolWithDefault|g' "$file"
  sed -i '' 's|types\.StringPointerValue|convert.GraphToFrameworkString|g' "$file"
  sed -i '' 's|types\.StringValue(\*|convert.GraphToFrameworkString(|g' "$file"
  
  echo "  Updated successfully."
done

echo "All files processed." 