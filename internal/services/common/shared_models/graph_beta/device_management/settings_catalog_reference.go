package sharedmodels

// DeviceConfigurationPolicyV2GraphServiceModelReference is an anomimous struct that represents the JSON structure of
// settings catalog settings fors windows, linux, macOS, and iOS. Used for device configuration, endpoint
// privilege management, templated device configuration.
// This struct is used as a reference for understanding the complex nested nature of settings catalog
// that's not explitly clear when using named go structs.This struct is not actively used in the codebase.
var DeviceConfigurationPolicyV2GraphServiceModelReference struct {
	SettingsDetails []struct {
		ID              string `json:"id"`
		SettingInstance struct {

			// For choice setting collections
			ChoiceSettingCollectionValue []struct {
				Children []struct {
					ODataType                        string `json:"@odata.type"`
					SettingDefinitionId              string `json:"settingDefinitionId"`
					SettingInstanceTemplateReference *struct {
						SettingInstanceTemplateId string `json:"settingInstanceTemplateId"`
					} `json:"settingInstanceTemplateReference,omitempty"`
					// For nested simple setting collection within choice setting collection
					SimpleSettingCollectionValue []struct {
						ODataType                     string `json:"@odata.type"`
						SettingValueTemplateReference *struct {
							SettingValueTemplateId string `json:"settingValueTemplateId"`
							UseTemplateDefault     bool   `json:"useTemplateDefault"`
						} `json:"settingValueTemplateReference,omitempty"`
						Value string `json:"value"`
					} `json:"simpleSettingCollectionValue,omitempty"`
					// For nested simple settings within choice setting collection
					SimpleSettingValue *struct {
						ODataType                     string `json:"@odata.type"`
						SettingValueTemplateReference *struct {
							SettingValueTemplateId string `json:"settingValueTemplateId"`
							UseTemplateDefault     bool   `json:"useTemplateDefault"`
						} `json:"settingValueTemplateReference,omitempty"`
						Value any `json:"value"`
					} `json:"simpleSettingValue,omitempty"`
				} `json:"children"`
				SettingValueTemplateReference *struct {
					SettingValueTemplateId string `json:"settingValueTemplateId"`
					UseTemplateDefault     bool   `json:"useTemplateDefault"`
				} `json:"settingValueTemplateReference,omitempty"`

				Value string `json:"value"`
			} `json:"choiceSettingCollectionValue,omitempty"`

			// For choice settings
			ChoiceSettingValue *struct {
				Children []struct {
					// For nested choice settings within choice children
					ChoiceSettingValue *struct {
						Children []struct {
							ODataType           string `json:"@odata.type"`
							SettingDefinitionId string `json:"settingDefinitionId"`
						} `json:"children"`
						SettingValueTemplateReference *struct {
							SettingValueTemplateId string `json:"settingValueTemplateId"`
							UseTemplateDefault     bool   `json:"useTemplateDefault"`
						} `json:"settingValueTemplateReference,omitempty"`
						Value string `json:"value"`
					} `json:"choiceSettingValue,omitempty"`
					// For GroupSettingCollectionValue within Choice children
					GroupSettingCollectionValue []struct {
						Children []struct {
							ODataType                        string `json:"@odata.type"`
							SettingDefinitionId              string `json:"settingDefinitionId"`
							SettingInstanceTemplateReference *struct {
								SettingInstanceTemplateId string `json:"settingInstanceTemplateId"`
							} `json:"settingInstanceTemplateReference,omitempty"`
							SimpleSettingValue *struct {
								ODataType                     string `json:"@odata.type"`
								SettingValueTemplateReference *struct {
									SettingValueTemplateId string `json:"settingValueTemplateId"`
									UseTemplateDefault     bool   `json:"useTemplateDefault"`
								} `json:"settingValueTemplateReference,omitempty"`
								Value any `json:"value"`
							} `json:"simpleSettingValue,omitempty"`
						} `json:"children"`
						SettingValueTemplateReference *struct {
							SettingValueTemplateId string `json:"settingValueTemplateId"`
							UseTemplateDefault     bool   `json:"useTemplateDefault"`
						} `json:"settingValueTemplateReference,omitempty"`
					} `json:"groupSettingCollectionValue,omitempty"`
					ODataType                        string `json:"@odata.type"`
					SettingDefinitionId              string `json:"settingDefinitionId"`
					SettingInstanceTemplateReference *struct {
						SettingInstanceTemplateId string `json:"settingInstanceTemplateId"`
					} `json:"settingInstanceTemplateReference,omitempty"`
					// For SimpleSettingCollectionValue within Choice children
					SimpleSettingCollectionValue []struct {
						ODataType                     string `json:"@odata.type"`
						SettingValueTemplateReference *struct {
							SettingValueTemplateId string `json:"settingValueTemplateId"`
							UseTemplateDefault     bool   `json:"useTemplateDefault"`
						} `json:"settingValueTemplateReference,omitempty"`
						Value string `json:"value"`
					} `json:"simpleSettingCollectionValue,omitempty"`
					// For simple settings within choice children
					SimpleSettingValue *struct {
						ODataType                     string `json:"@odata.type"`
						SettingValueTemplateReference *struct {
							SettingValueTemplateId string `json:"settingValueTemplateId"`
							UseTemplateDefault     bool   `json:"useTemplateDefault"`
						} `json:"settingValueTemplateReference,omitempty"`
						Value any `json:"value"`
					} `json:"simpleSettingValue,omitempty"`
				} `json:"children"`
				SettingValueTemplateReference *struct {
					SettingValueTemplateId string `json:"settingValueTemplateId"`
					UseTemplateDefault     bool   `json:"useTemplateDefault"`
				} `json:"settingValueTemplateReference,omitempty"`
				Value string `json:"value"`
			} `json:"choiceSettingValue,omitempty"`

			// For group setting collections
			GroupSettingCollectionValue []struct {
				Children []struct {
					// For nested choice settings within group setting collection (Level 2)
					ChoiceSettingValue *struct {
						Children []struct {
							ChoiceSettingValue *struct {
								Children []struct {
									ODataType           string `json:"@odata.type"`
									SettingDefinitionId string `json:"settingDefinitionId"`
								} `json:"children"`
								SettingValueTemplateReference *struct {
									SettingValueTemplateId string `json:"settingValueTemplateId"`
									UseTemplateDefault     bool   `json:"useTemplateDefault"`
								} `json:"settingValueTemplateReference,omitempty"`
								Value string `json:"value"`
							} `json:"choiceSettingValue,omitempty"`
							ODataType                        string `json:"@odata.type"`
							SettingDefinitionId              string `json:"settingDefinitionId"`
							SettingInstanceTemplateReference *struct {
								SettingInstanceTemplateId string `json:"settingInstanceTemplateId"`
							} `json:"settingInstanceTemplateReference,omitempty"`
							SimpleSettingValue *struct {
								ODataType                     string `json:"@odata.type"`
								SettingValueTemplateReference *struct {
									SettingValueTemplateId string `json:"settingValueTemplateId"`
									UseTemplateDefault     bool   `json:"useTemplateDefault"`
								} `json:"settingValueTemplateReference,omitempty"`
								Value any `json:"value"`
							} `json:"simpleSettingValue,omitempty"`
						} `json:"children"`
						SettingValueTemplateReference *struct {
							SettingValueTemplateId string `json:"settingValueTemplateId"`
							UseTemplateDefault     bool   `json:"useTemplateDefault"`
						} `json:"settingValueTemplateReference,omitempty"`
						Value string `json:"value"`
					} `json:"choiceSettingValue,omitempty"`
					// For nested group setting collections within group setting collection (Level 2)
					GroupSettingCollectionValue []struct {
						Children []struct {
							// For nested choice settings within group setting collection within group setting collection (Level 3)
							ChoiceSettingValue *struct {
								Children []struct {
									ODataType           string `json:"@odata.type"`
									SettingDefinitionId string `json:"settingDefinitionId"`
									SimpleSettingValue  *struct {
										ODataType                     string `json:"@odata.type"`
										SettingValueTemplateReference *struct {
											SettingValueTemplateId string `json:"settingValueTemplateId"`
											UseTemplateDefault     bool   `json:"useTemplateDefault"`
										} `json:"settingValueTemplateReference,omitempty"`
										Value      any    `json:"value"`
										ValueState string `json:"valueState,omitempty"`
									} `json:"simpleSettingValue,omitempty"`
								} `json:"children"`
								SettingValueTemplateReference *struct {
									SettingValueTemplateId string `json:"settingValueTemplateId"`
									UseTemplateDefault     bool   `json:"useTemplateDefault"`
								} `json:"settingValueTemplateReference,omitempty"`
								Value string `json:"value"`
							} `json:"choiceSettingValue,omitempty"`
							// For nested group setting collections within group setting collection within group setting collection (Level 3)
							GroupSettingCollectionValue []struct {
								Children []struct {
									// For nested choice settings within group setting collection within group setting collection within group setting collection (Level 4)
									ChoiceSettingValue *struct {
										Children []struct {
											ODataType                        string `json:"@odata.type"`
											SettingDefinitionId              string `json:"settingDefinitionId"`
											SettingInstanceTemplateReference *struct {
												SettingInstanceTemplateId string `json:"settingInstanceTemplateId"`
											} `json:"settingInstanceTemplateReference,omitempty"`
											// For nested simple settings within choice settings within group setting collection within group setting collection within group setting collection (Level 5)
											SimpleSettingValue *struct {
												ODataType                     string `json:"@odata.type"`
												SettingValueTemplateReference *struct {
													SettingValueTemplateId string `json:"settingValueTemplateId"`
													UseTemplateDefault     bool   `json:"useTemplateDefault"`
												} `json:"settingValueTemplateReference,omitempty"`
												Value any `json:"value"`
											} `json:"simpleSettingValue,omitempty"`
										} `json:"children"`
										SettingValueTemplateReference *struct {
											SettingValueTemplateId string `json:"settingValueTemplateId"`
											UseTemplateDefault     bool   `json:"useTemplateDefault"`
										} `json:"settingValueTemplateReference,omitempty"`
										Value string `json:"value"`
									} `json:"choiceSettingValue,omitempty"`
									ODataType                        string `json:"@odata.type"`
									SettingDefinitionId              string `json:"settingDefinitionId"`
									SettingInstanceTemplateReference *struct {
										SettingInstanceTemplateId string `json:"settingInstanceTemplateId"`
									} `json:"settingInstanceTemplateReference,omitempty"`
									// For simple settings collection within group setting collection within group setting collection within group setting collection (Level 4)
									SimpleSettingCollectionValue []struct {
										ODataType                     string `json:"@odata.type"`
										SettingValueTemplateReference *struct {
											SettingValueTemplateId string `json:"settingValueTemplateId"`
											UseTemplateDefault     bool   `json:"useTemplateDefault"`
										} `json:"settingValueTemplateReference,omitempty"`
										Value string `json:"value"`
									} `json:"simpleSettingCollectionValue,omitempty"`
									// For simple settings within group setting collection within group setting collection within group setting collection (Level 4)
									SimpleSettingValue *struct {
										ODataType                     string `json:"@odata.type"`
										SettingValueTemplateReference *struct {
											SettingValueTemplateId string `json:"settingValueTemplateId"`
											UseTemplateDefault     bool   `json:"useTemplateDefault"`
										} `json:"settingValueTemplateReference,omitempty"`
										Value any `json:"value"`
									} `json:"simpleSettingValue,omitempty"`
								} `json:"children"`
								SettingValueTemplateReference *struct {
									SettingValueTemplateId string `json:"settingValueTemplateId"`
									UseTemplateDefault     bool   `json:"useTemplateDefault"`
								} `json:"settingValueTemplateReference,omitempty"`
							} `json:"groupSettingCollectionValue,omitempty"`
							ODataType                        string `json:"@odata.type"`
							SettingDefinitionId              string `json:"settingDefinitionId"`
							SettingInstanceTemplateReference *struct {
								SettingInstanceTemplateId string `json:"settingInstanceTemplateId"`
							} `json:"settingInstanceTemplateReference,omitempty"`
							// For nested simple setting collections within group setting collection within group setting collection (Level 3)
							SimpleSettingCollectionValue []struct {
								ODataType                     string `json:"@odata.type"`
								SettingValueTemplateReference *struct {
									SettingValueTemplateId string `json:"settingValueTemplateId"`
									UseTemplateDefault     bool   `json:"useTemplateDefault"`
								} `json:"settingValueTemplateReference,omitempty"`
								Value string `json:"value"`
							} `json:"simpleSettingCollectionValue,omitempty"`
							// For nested simple settings within group setting collection within group setting collection (Level 3)
							SimpleSettingValue *struct {
								ODataType                     string `json:"@odata.type"`
								SettingValueTemplateReference *struct {
									SettingValueTemplateId string `json:"settingValueTemplateId"`
									UseTemplateDefault     bool   `json:"useTemplateDefault"`
								} `json:"settingValueTemplateReference,omitempty"`
								Value      any    `json:"value"`
								ValueState string `json:"valueState,omitempty"`
							} `json:"simpleSettingValue,omitempty"`
						} `json:"children"`
						SettingValueTemplateReference *struct {
							SettingValueTemplateId string `json:"settingValueTemplateId"`
							UseTemplateDefault     bool   `json:"useTemplateDefault"`
						} `json:"settingValueTemplateReference,omitempty"`
					} `json:"groupSettingCollectionValue,omitempty"`
					ODataType                        string `json:"@odata.type"`
					SettingDefinitionId              string `json:"settingDefinitionId"`
					SettingInstanceTemplateReference *struct {
						SettingInstanceTemplateId string `json:"settingInstanceTemplateId"`
					} `json:"settingInstanceTemplateReference,omitempty"`
					// For nested simple settings (string, integer, secret) within group setting collection  (Level 2)
					SimpleSettingValue *struct {
						ODataType                     string `json:"@odata.type"`
						SettingValueTemplateReference *struct {
							SettingValueTemplateId string `json:"settingValueTemplateId"`
							UseTemplateDefault     bool   `json:"useTemplateDefault"`
						} `json:"settingValueTemplateReference,omitempty"`
						Value      any    `json:"value"`
						ValueState string `json:"valueState,omitempty"`
					} `json:"simpleSettingValue,omitempty"`
					// For nested simple setting collections within group setting collection (Level 2)
					SimpleSettingCollectionValue []struct {
						ODataType                     string `json:"@odata.type"`
						SettingValueTemplateReference *struct {
							SettingValueTemplateId string `json:"settingValueTemplateId"`
							UseTemplateDefault     bool   `json:"useTemplateDefault"`
						} `json:"settingValueTemplateReference,omitempty"`
						Value string `json:"value"`
					} `json:"simpleSettingCollectionValue,omitempty"`
				} `json:"children"`
				SettingValueTemplateReference *struct {
					SettingValueTemplateId string `json:"settingValueTemplateId"`
					UseTemplateDefault     bool   `json:"useTemplateDefault"`
				} `json:"settingValueTemplateReference,omitempty"`
			} `json:"groupSettingCollectionValue,omitempty"`

			// Setting instance Odata and template reference
			ODataType                        string `json:"@odata.type"`
			SettingDefinitionId              string `json:"settingDefinitionId"`
			SettingInstanceTemplateReference *struct {
				SettingInstanceTemplateId string `json:"settingInstanceTemplateId"`
			} `json:"settingInstanceTemplateReference,omitempty"`

			// For simple collection settings
			SimpleSettingCollectionValue []struct {
				ODataType                     string `json:"@odata.type"`
				SettingValueTemplateReference *struct {
					SettingValueTemplateId string `json:"settingValueTemplateId"`
					UseTemplateDefault     bool   `json:"useTemplateDefault"`
				} `json:"settingValueTemplateReference,omitempty"`

				Value string `json:"value"`
			} `json:"simpleSettingCollectionValue,omitempty"`

			// For simple settings
			SimpleSettingValue *struct {
				ODataType                     string `json:"@odata.type"`
				SettingValueTemplateReference *struct {
					SettingValueTemplateId string `json:"settingValueTemplateId"`
					UseTemplateDefault     bool   `json:"useTemplateDefault"`
				} `json:"settingValueTemplateReference,omitempty"`
				Value any `json:"value"`
			} `json:"simpleSettingValue,omitempty"`
		} `json:"settingInstance"`
	} `json:"settings"`
}
