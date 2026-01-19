package utilityDeploymentScheduler

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestIsWithinTimeWindow(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name        string
		checkTime   time.Time
		window      TimeWindowModel
		expected    bool
		expectError bool
	}{
		{
			name:      "Day of week match - Monday",
			checkTime: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC), // Monday
			window: TimeWindowModel{
				DaysOfWeek: types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue("monday"),
					types.StringValue("wednesday"),
					types.StringValue("friday"),
				}),
				TimeOfDayStart: types.StringNull(),
				TimeOfDayEnd:   types.StringNull(),
				DateStart:      types.StringNull(),
				DateEnd:        types.StringNull(),
			},
			expected:    true,
			expectError: false,
		},
		{
			name:      "Day of week no match - Tuesday",
			checkTime: time.Date(2024, 1, 2, 12, 0, 0, 0, time.UTC), // Tuesday
			window: TimeWindowModel{
				DaysOfWeek: types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue("monday"),
					types.StringValue("wednesday"),
					types.StringValue("friday"),
				}),
				TimeOfDayStart: types.StringNull(),
				TimeOfDayEnd:   types.StringNull(),
				DateStart:      types.StringNull(),
				DateEnd:        types.StringNull(),
			},
			expected:    false,
			expectError: false,
		},
		{
			name:      "Time of day within range",
			checkTime: time.Date(2024, 1, 1, 14, 30, 0, 0, time.UTC),
			window: TimeWindowModel{
				DaysOfWeek:     types.ListNull(types.StringType),
				TimeOfDayStart: types.StringValue("09:00:00"),
				TimeOfDayEnd:   types.StringValue("17:00:00"),
				DateStart:      types.StringNull(),
				DateEnd:        types.StringNull(),
			},
			expected:    true,
			expectError: false,
		},
		{
			name:      "Time of day before range",
			checkTime: time.Date(2024, 1, 1, 8, 30, 0, 0, time.UTC),
			window: TimeWindowModel{
				DaysOfWeek:     types.ListNull(types.StringType),
				TimeOfDayStart: types.StringValue("09:00:00"),
				TimeOfDayEnd:   types.StringValue("17:00:00"),
				DateStart:      types.StringNull(),
				DateEnd:        types.StringNull(),
			},
			expected:    false,
			expectError: false,
		},
		{
			name:      "Time of day after range",
			checkTime: time.Date(2024, 1, 1, 18, 30, 0, 0, time.UTC),
			window: TimeWindowModel{
				DaysOfWeek:     types.ListNull(types.StringType),
				TimeOfDayStart: types.StringValue("09:00:00"),
				TimeOfDayEnd:   types.StringValue("17:00:00"),
				DateStart:      types.StringNull(),
				DateEnd:        types.StringNull(),
			},
			expected:    false,
			expectError: false,
		},
		{
			name:      "Date range within range",
			checkTime: time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC),
			window: TimeWindowModel{
				DaysOfWeek:     types.ListNull(types.StringType),
				TimeOfDayStart: types.StringNull(),
				TimeOfDayEnd:   types.StringNull(),
				DateStart:      types.StringValue("2024-01-01T00:00:00Z"),
				DateEnd:        types.StringValue("2024-12-31T23:59:59Z"),
			},
			expected:    true,
			expectError: false,
		},
		{
			name:      "Date range before range",
			checkTime: time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC),
			window: TimeWindowModel{
				DaysOfWeek:     types.ListNull(types.StringType),
				TimeOfDayStart: types.StringNull(),
				TimeOfDayEnd:   types.StringNull(),
				DateStart:      types.StringValue("2024-01-01T00:00:00Z"),
				DateEnd:        types.StringValue("2024-12-31T23:59:59Z"),
			},
			expected:    false,
			expectError: false,
		},
		{
			name:      "Date range after range",
			checkTime: time.Date(2025, 1, 1, 0, 0, 1, 0, time.UTC),
			window: TimeWindowModel{
				DaysOfWeek:     types.ListNull(types.StringType),
				TimeOfDayStart: types.StringNull(),
				TimeOfDayEnd:   types.StringNull(),
				DateStart:      types.StringValue("2024-01-01T00:00:00Z"),
				DateEnd:        types.StringValue("2024-12-31T23:59:59Z"),
			},
			expected:    false,
			expectError: false,
		},
		{
			name:      "Combined match - day, time, and date all match",
			checkTime: time.Date(2024, 6, 17, 14, 30, 0, 0, time.UTC), // Monday
			window: TimeWindowModel{
				DaysOfWeek: types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue("monday"),
					types.StringValue("wednesday"),
					types.StringValue("friday"),
				}),
				TimeOfDayStart: types.StringValue("09:00:00"),
				TimeOfDayEnd:   types.StringValue("17:00:00"),
				DateStart:      types.StringValue("2024-01-01T00:00:00Z"),
				DateEnd:        types.StringValue("2024-12-31T23:59:59Z"),
			},
			expected:    true,
			expectError: false,
		},
		{
			name:      "Combined no match - day matches but time outside range",
			checkTime: time.Date(2024, 6, 17, 18, 30, 0, 0, time.UTC), // Monday
			window: TimeWindowModel{
				DaysOfWeek: types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue("monday"),
					types.StringValue("wednesday"),
					types.StringValue("friday"),
				}),
				TimeOfDayStart: types.StringValue("09:00:00"),
				TimeOfDayEnd:   types.StringValue("17:00:00"),
				DateStart:      types.StringValue("2024-01-01T00:00:00Z"),
				DateEnd:        types.StringValue("2024-12-31T23:59:59Z"),
			},
			expected:    false,
			expectError: false,
		},
		{
			name:      "No restrictions - all null",
			checkTime: time.Date(2024, 6, 17, 18, 30, 0, 0, time.UTC),
			window: TimeWindowModel{
				DaysOfWeek:     types.ListNull(types.StringType),
				TimeOfDayStart: types.StringNull(),
				TimeOfDayEnd:   types.StringNull(),
				DateStart:      types.StringNull(),
				DateEnd:        types.StringNull(),
			},
			expected:    true,
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := isWithinTimeWindow(ctx, tc.checkTime, tc.window)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tc.expected {
					t.Errorf("Expected %v, got %v", tc.expected, result)
				}
			}
		})
	}
}

func TestEvaluateInclusionWindows(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name               string
		currentTime        time.Time
		inclusionWindows   types.Object
		expectedMatch      bool
		expectedMsgPattern string
		expectError        bool
	}{
		{
			name:        "Null inclusion windows - always allowed",
			currentTime: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			inclusionWindows: types.ObjectNull(map[string]attr.Type{
				"window": types.ListType{
					ElemType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"days_of_week":      types.ListType{ElemType: types.StringType},
							"time_of_day_start": types.StringType,
							"time_of_day_end":   types.StringType,
							"date_start":        types.StringType,
							"date_end":          types.StringType,
						},
					},
				},
			}),
			expectedMatch:      true,
			expectedMsgPattern: "",
			expectError:        false,
		},
		{
			name:        "Single window match - Monday within business hours",
			currentTime: time.Date(2024, 1, 1, 14, 0, 0, 0, time.UTC), // Monday
			inclusionWindows: createInclusionWindowsObject([]TimeWindowModel{
				{
					DaysOfWeek: types.ListValueMust(types.StringType, []attr.Value{
						types.StringValue("monday"),
					}),
					TimeOfDayStart: types.StringValue("09:00:00"),
					TimeOfDayEnd:   types.StringValue("17:00:00"),
					DateStart:      types.StringNull(),
					DateEnd:        types.StringNull(),
				},
			}),
			expectedMatch:      true,
			expectedMsgPattern: "within inclusion window 1",
			expectError:        false,
		},
		{
			name:        "Single window no match - Wrong day",
			currentTime: time.Date(2024, 1, 2, 14, 0, 0, 0, time.UTC), // Tuesday
			inclusionWindows: createInclusionWindowsObject([]TimeWindowModel{
				{
					DaysOfWeek: types.ListValueMust(types.StringType, []attr.Value{
						types.StringValue("monday"),
					}),
					TimeOfDayStart: types.StringValue("09:00:00"),
					TimeOfDayEnd:   types.StringValue("17:00:00"),
					DateStart:      types.StringNull(),
					DateEnd:        types.StringNull(),
				},
			}),
			expectedMatch:      false,
			expectedMsgPattern: "outside all inclusion windows",
			expectError:        false,
		},
		{
			name:        "Multiple windows - match second window",
			currentTime: time.Date(2024, 1, 2, 14, 0, 0, 0, time.UTC), // Tuesday
			inclusionWindows: createInclusionWindowsObject([]TimeWindowModel{
				{
					DaysOfWeek: types.ListValueMust(types.StringType, []attr.Value{
						types.StringValue("monday"),
					}),
					TimeOfDayStart: types.StringNull(),
					TimeOfDayEnd:   types.StringNull(),
					DateStart:      types.StringNull(),
					DateEnd:        types.StringNull(),
				},
				{
					DaysOfWeek: types.ListValueMust(types.StringType, []attr.Value{
						types.StringValue("tuesday"),
					}),
					TimeOfDayStart: types.StringNull(),
					TimeOfDayEnd:   types.StringNull(),
					DateStart:      types.StringNull(),
					DateEnd:        types.StringNull(),
				},
			}),
			expectedMatch:      true,
			expectedMsgPattern: "within inclusion window",
			expectError:        false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			match, _, err := evaluateInclusionWindows(ctx, tc.currentTime, tc.inclusionWindows)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if match != tc.expectedMatch {
					t.Errorf("Expected match=%v, got %v", tc.expectedMatch, match)
				}
			}
		})
	}
}

func TestEvaluateExclusionWindows(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name               string
		currentTime        time.Time
		exclusionWindows   types.Object
		expectedBlocked    bool
		expectedMsgPattern string
		expectError        bool
	}{
		{
			name:        "Null exclusion windows - not blocked",
			currentTime: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			exclusionWindows: types.ObjectNull(map[string]attr.Type{
				"window": types.ListType{
					ElemType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"days_of_week":      types.ListType{ElemType: types.StringType},
							"time_of_day_start": types.StringType,
							"time_of_day_end":   types.StringType,
							"date_start":        types.StringType,
							"date_end":          types.StringType,
						},
					},
				},
			}),
			expectedBlocked:    false,
			expectedMsgPattern: "",
			expectError:        false,
		},
		{
			name:        "Single window blocks - Weekend",
			currentTime: time.Date(2024, 1, 6, 14, 0, 0, 0, time.UTC), // Saturday
			exclusionWindows: createExclusionWindowsObject([]TimeWindowModel{
				{
					DaysOfWeek: types.ListValueMust(types.StringType, []attr.Value{
						types.StringValue("saturday"),
						types.StringValue("sunday"),
					}),
					TimeOfDayStart: types.StringNull(),
					TimeOfDayEnd:   types.StringNull(),
					DateStart:      types.StringNull(),
					DateEnd:        types.StringNull(),
				},
			}),
			expectedBlocked:    true,
			expectedMsgPattern: "within exclusion window 1",
			expectError:        false,
		},
		{
			name:        "Single window does not block - Weekday",
			currentTime: time.Date(2024, 1, 1, 14, 0, 0, 0, time.UTC), // Monday
			exclusionWindows: createExclusionWindowsObject([]TimeWindowModel{
				{
					DaysOfWeek: types.ListValueMust(types.StringType, []attr.Value{
						types.StringValue("saturday"),
						types.StringValue("sunday"),
					}),
					TimeOfDayStart: types.StringNull(),
					TimeOfDayEnd:   types.StringNull(),
					DateStart:      types.StringNull(),
					DateEnd:        types.StringNull(),
				},
			}),
			expectedBlocked:    false,
			expectedMsgPattern: "outside all exclusion windows",
			expectError:        false,
		},
		{
			name:        "Multiple windows - second window blocks",
			currentTime: time.Date(2024, 12, 25, 14, 0, 0, 0, time.UTC), // Holiday
			exclusionWindows: createExclusionWindowsObject([]TimeWindowModel{
				{
					DaysOfWeek: types.ListValueMust(types.StringType, []attr.Value{
						types.StringValue("saturday"),
						types.StringValue("sunday"),
					}),
					TimeOfDayStart: types.StringNull(),
					TimeOfDayEnd:   types.StringNull(),
					DateStart:      types.StringNull(),
					DateEnd:        types.StringNull(),
				},
				{
					DaysOfWeek:     types.ListNull(types.StringType),
					TimeOfDayStart: types.StringNull(),
					TimeOfDayEnd:   types.StringNull(),
					DateStart:      types.StringValue("2024-12-20T00:00:00Z"),
					DateEnd:        types.StringValue("2025-01-05T23:59:59Z"),
				},
			}),
			expectedBlocked:    true,
			expectedMsgPattern: "within exclusion window",
			expectError:        false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			blocked, _, err := evaluateExclusionWindows(ctx, tc.currentTime, tc.exclusionWindows)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if blocked != tc.expectedBlocked {
					t.Errorf("Expected blocked=%v, got %v", tc.expectedBlocked, blocked)
				}
			}
		})
	}
}

// Helper functions to create test objects

func createInclusionWindowsObject(windows []TimeWindowModel) types.Object {
	windowList := make([]attr.Value, len(windows))
	for i, w := range windows {
		windowList[i] = types.ObjectValueMust(
			map[string]attr.Type{
				"days_of_week":      types.ListType{ElemType: types.StringType},
				"time_of_day_start": types.StringType,
				"time_of_day_end":   types.StringType,
				"date_start":        types.StringType,
				"date_end":          types.StringType,
			},
			map[string]attr.Value{
				"days_of_week":      w.DaysOfWeek,
				"time_of_day_start": w.TimeOfDayStart,
				"time_of_day_end":   w.TimeOfDayEnd,
				"date_start":        w.DateStart,
				"date_end":          w.DateEnd,
			},
		)
	}

	return types.ObjectValueMust(
		map[string]attr.Type{
			"window": types.ListType{
				ElemType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"days_of_week":      types.ListType{ElemType: types.StringType},
						"time_of_day_start": types.StringType,
						"time_of_day_end":   types.StringType,
						"date_start":        types.StringType,
						"date_end":          types.StringType,
					},
				},
			},
		},
		map[string]attr.Value{
			"window": types.ListValueMust(
				types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"days_of_week":      types.ListType{ElemType: types.StringType},
						"time_of_day_start": types.StringType,
						"time_of_day_end":   types.StringType,
						"date_start":        types.StringType,
						"date_end":          types.StringType,
					},
				},
				windowList,
			),
		},
	)
}

func createExclusionWindowsObject(windows []TimeWindowModel) types.Object {
	windowList := make([]attr.Value, len(windows))
	for i, w := range windows {
		windowList[i] = types.ObjectValueMust(
			map[string]attr.Type{
				"days_of_week":      types.ListType{ElemType: types.StringType},
				"time_of_day_start": types.StringType,
				"time_of_day_end":   types.StringType,
				"date_start":        types.StringType,
				"date_end":          types.StringType,
			},
			map[string]attr.Value{
				"days_of_week":      w.DaysOfWeek,
				"time_of_day_start": w.TimeOfDayStart,
				"time_of_day_end":   w.TimeOfDayEnd,
				"date_start":        w.DateStart,
				"date_end":          w.DateEnd,
			},
		)
	}

	return types.ObjectValueMust(
		map[string]attr.Type{
			"window": types.ListType{
				ElemType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"days_of_week":      types.ListType{ElemType: types.StringType},
						"time_of_day_start": types.StringType,
						"time_of_day_end":   types.StringType,
						"date_start":        types.StringType,
						"date_end":          types.StringType,
					},
				},
			},
		},
		map[string]attr.Value{
			"window": types.ListValueMust(
				types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"days_of_week":      types.ListType{ElemType: types.StringType},
						"time_of_day_start": types.StringType,
						"time_of_day_end":   types.StringType,
						"date_start":        types.StringType,
						"date_end":          types.StringType,
					},
				},
				windowList,
			),
		},
	)
}
