// custom config profiles

Request URL
https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations
Request Method
POST

request body

{
  "id":"00000000-0000-0000-0000-000000000000",
  "displayName":"test",
  "description":"test",
  "roleScopeTagIds":["0"],
  "@odata.type":"#microsoft.graph.macOSCustomConfiguration",
  "deploymentChannel":"deviceChannel",
  "payloadFileName":"dt-mcp-accessibility_hearing_base-prod-v0.0.1.mobileconfig","payload":"PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPCFET0NUWVBFIHBsaXN0IFBVQkxJQyAiLS8vQXBwbGUvL0RURCBQTElTVCAxLjAvL0VOIiAiaHR0cDovL3d3dy5hcHBsZS5jb20vRFREcy9Qcm9wZXJ0eUxpc3QtMS4wLmR0ZCI+CjxwbGlzdCB2ZXJzaW9uPSIxLjAiPgoJPGRpY3Q+CgkJPGtleT5QYXlsb2FkQ29udGVudDwva2V5PgoJCTxhcnJheT4KCQkJPGRpY3Q+CgkJCQk8a2V5PlBheWxvYWREZXNjcmlwdGlvbjwva2V5PgoJCQkJPHN0cmluZy8+CgkJCQk8a2V5PlBheWxvYWREaXNwbGF5TmFtZTwva2V5PgoJCQkJPHN0cmluZz5BY2Nlc3NpYmlsaXR5PC9zdHJpbmc+CgkJCQk8a2V5PlBheWxvYWRFbmFibGVkPC9rZXk+CgkJCQk8dHJ1ZS8+CgkJCQk8a2V5PlBheWxvYWRJZGVudGlmaWVyPC9rZXk+CgkJCQk8c3RyaW5nPjVGNDlBRTVGLTU4ODQtNEUxQi04Nzc5LTFBNDdDMzg3MDdBQzwvc3RyaW5nPgoJCQkJPGtleT5QYXlsb2FkT3JnYW5pemF0aW9uPC9rZXk+CgkJCQk8c3RyaW5nPkRlcGxveW1lbnQgVGhlb3J5PC9zdHJpbmc+CgkJCQk8a2V5PlBheWxvYWRUeXBlPC9rZXk+CgkJCQk8c3RyaW5nPmNvbS5hcHBsZS51bml2ZXJzYWxhY2Nlc3M8L3N0cmluZz4KCQkJCTxrZXk+UGF5bG9hZFVVSUQ8L2tleT4KCQkJCTxzdHJpbmc+NUY0OUFFNUYtNTg4NC00RTFCLTg3NzktMUE0N0MzODcwN0FDPC9zdHJpbmc+CgkJCQk8a2V5PlBheWxvYWRWZXJzaW9uPC9rZXk+CgkJCQk8aW50ZWdlcj4xPC9pbnRlZ2VyPgoJCQkJPGtleT5jbG9zZVZpZXdGYXJQb2ludDwva2V5PgoJCQkJPGludGVnZXI+MTwvaW50ZWdlcj4KCQkJCTxrZXk+Y2xvc2VWaWV3SG90a2V5c0VuYWJsZWQ8L2tleT4KCQkJCTxmYWxzZS8+CgkJCQk8a2V5PmNsb3NlVmlld05lYXJQb2ludDwva2V5PgoJCQkJPGludGVnZXI+MTA8L2ludGVnZXI+CgkJCQk8a2V5PmNsb3NlVmlld1Njcm9sbFdoZWVsVG9nZ2xlPC9rZXk+CgkJCQk8ZmFsc2UvPgoJCQkJPGtleT5jbG9zZVZpZXdTaG93UHJldmlldzwva2V5PgoJCQkJPGZhbHNlLz4KCQkJCTxrZXk+Y2xvc2VWaWV3U21vb3RoSW1hZ2VzPC9rZXk+CgkJCQk8ZmFsc2UvPgoJCQkJPGtleT5jb250cmFzdDwva2V5PgoJCQkJPGludGVnZXI+MDwvaW50ZWdlcj4KCQkJCTxrZXk+Zmxhc2hTY3JlZW48L2tleT4KCQkJCTx0cnVlLz4KCQkJCTxrZXk+Z3JheXNjYWxlPC9rZXk+CgkJCQk8ZmFsc2UvPgoJCQkJPGtleT5tb3VzZURyaXZlcjwva2V5PgoJCQkJPGZhbHNlLz4KCQkJCTxrZXk+bW91c2VEcml2ZXJDdXJzb3JTaXplPC9rZXk+CgkJCQk8aW50ZWdlcj4xPC9pbnRlZ2VyPgoJCQkJPGtleT5tb3VzZURyaXZlcklnbm9yZVRyYWNrcGFkPC9rZXk+CgkJCQk8ZmFsc2UvPgoJCQkJPGtleT5tb3VzZURyaXZlckluaXRpYWxEZWxheTwva2V5PgoJCQkJPHJlYWw+MTwvcmVhbD4KCQkJCTxrZXk+bW91c2VEcml2ZXJNYXhTcGVlZDwva2V5PgoJCQkJPGludGVnZXI+MzwvaW50ZWdlcj4KCQkJCTxrZXk+c2xvd0tleTwva2V5PgoJCQkJPGZhbHNlLz4KCQkJCTxrZXk+c2xvd0tleUJlZXBPbjwva2V5PgoJCQkJPGZhbHNlLz4KCQkJCTxrZXk+c2xvd0tleURlbGF5PC9rZXk+CgkJCQk8aW50ZWdlcj4wPC9pbnRlZ2VyPgoJCQkJPGtleT5zdGVyZW9Bc01vbm88L2tleT4KCQkJCTxmYWxzZS8+CgkJCQk8a2V5PnN0aWNreUtleTwva2V5PgoJCQkJPGZhbHNlLz4KCQkJCTxrZXk+c3RpY2t5S2V5QmVlcE9uTW9kaWZpZXI8L2tleT4KCQkJCTxmYWxzZS8+CgkJCQk8a2V5PnN0aWNreUtleVNob3dXaW5kb3c8L2tleT4KCQkJCTxmYWxzZS8+CgkJCQk8a2V5PnZvaWNlT3Zlck9uT2ZmS2V5PC9rZXk+CgkJCQk8ZmFsc2UvPgoJCQkJPGtleT53aGl0ZU9uQmxhY2s8L2tleT4KCQkJCTxmYWxzZS8+CgkJCTwvZGljdD4KCQk8L2FycmF5PgoJCTxrZXk+UGF5bG9hZERlc2NyaXB0aW9uPC9rZXk+CgkJPHN0cmluZy8+CgkJPGtleT5QYXlsb2FkRGlzcGxheU5hbWU8L2tleT4KCQk8c3RyaW5nPmR0LW1jcC1hY2Nlc3NpYmlsaXR5X2hlYXJpbmdfYmFzZS0wLjAuMS1wcm9kLWV1LTAtMDwvc3RyaW5nPgoJCTxrZXk+UGF5bG9hZEVuYWJsZWQ8L2tleT4KCQk8dHJ1ZS8+CgkJPGtleT5QYXlsb2FkSWRlbnRpZmllcjwva2V5PgoJCTxzdHJpbmc+YTVkY2Y3YTQtNzJiMi00NmYyLWFhMmMtNWNkMjQ3Zjg4ZWE5PC9zdHJpbmc+CgkJPGtleT5QYXlsb2FkT3JnYW5pemF0aW9uPC9rZXk+CgkJPHN0cmluZz5EZXBsb3ltZW50IFRoZW9yeTwvc3RyaW5nPgoJCTxrZXk+UGF5bG9hZFJlbW92YWxEaXNhbGxvd2VkPC9rZXk+CgkJPHRydWUvPgoJCTxrZXk+UGF5bG9hZFNjb3BlPC9rZXk+CgkJPHN0cmluZz5TeXN0ZW08L3N0cmluZz4KCQk8a2V5PlBheWxvYWRUeXBlPC9rZXk+CgkJPHN0cmluZz5Db25maWd1cmF0aW9uPC9zdHJpbmc+CgkJPGtleT5QYXlsb2FkVVVJRDwva2V5PgoJCTxzdHJpbmc+YTVkY2Y3YTQtNzJiMi00NmYyLWFhMmMtNWNkMjQ3Zjg4ZWE5PC9zdHJpbmc+CgkJPGtleT5QYXlsb2FkVmVyc2lvbjwva2V5PgoJCTxpbnRlZ2VyPjE8L2ludGVnZXI+Cgk8L2RpY3Q+CjwvcGxpc3Q+",
  "payloadName":"custom-config-profile-name"
  }

response

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceConfigurations/$entity",
    "@odata.type": "#microsoft.graph.macOSCustomConfiguration",
    "id": "f1bde6e1-5b28-4abf-9f82-07bda322131d",
    "lastModifiedDateTime": "2025-08-19T07:46:30.8741816Z",
    "roleScopeTagIds": [
        "0"
    ],
    "supportsScopeTags": true,
    "deviceManagementApplicabilityRuleOsEdition": null,
    "deviceManagementApplicabilityRuleOsVersion": null,
    "deviceManagementApplicabilityRuleDeviceMode": null,
    "createdDateTime": "2025-08-19T07:46:30.8741816Z",
    "description": "test",
    "displayName": "test",
    "version": 1,
    "payloadName": "custom-config-profile-name",
    "payloadFileName": "dt-mcp-accessibility_hearing_base-prod-v0.0.1.mobileconfig",
    "payload": "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPCFET0NUWVBFIHBsaXN0IFBVQkxJQyAiLS8vQXBwbGUvL0RURCBQTElTVCAxLjAvL0VOIiAiaHR0cDovL3d3dy5hcHBsZS5jb20vRFREcy9Qcm9wZXJ0eUxpc3QtMS4wLmR0ZCI+CjxwbGlzdCB2ZXJzaW9uPSIxLjAiPgoJPGRpY3Q+CgkJPGtleT5QYXlsb2FkQ29udGVudDwva2V5PgoJCTxhcnJheT4KCQkJPGRpY3Q+CgkJCQk8a2V5PlBheWxvYWREZXNjcmlwdGlvbjwva2V5PgoJCQkJPHN0cmluZy8+CgkJCQk8a2V5PlBheWxvYWREaXNwbGF5TmFtZTwva2V5PgoJCQkJPHN0cmluZz5BY2Nlc3NpYmlsaXR5PC9zdHJpbmc+CgkJCQk8a2V5PlBheWxvYWRFbmFibGVkPC9rZXk+CgkJCQk8dHJ1ZS8+CgkJCQk8a2V5PlBheWxvYWRJZGVudGlmaWVyPC9rZXk+CgkJCQk8c3RyaW5nPjVGNDlBRTVGLTU4ODQtNEUxQi04Nzc5LTFBNDdDMzg3MDdBQzwvc3RyaW5nPgoJCQkJPGtleT5QYXlsb2FkT3JnYW5pemF0aW9uPC9rZXk+CgkJCQk8c3RyaW5nPkRlcGxveW1lbnQgVGhlb3J5PC9zdHJpbmc+CgkJCQk8a2V5PlBheWxvYWRUeXBlPC9rZXk+CgkJCQk8c3RyaW5nPmNvbS5hcHBsZS51bml2ZXJzYWxhY2Nlc3M8L3N0cmluZz4KCQkJCTxrZXk+UGF5bG9hZFVVSUQ8L2tleT4KCQkJCTxzdHJpbmc+NUY0OUFFNUYtNTg4NC00RTFCLTg3NzktMUE0N0MzODcwN0FDPC9zdHJpbmc+CgkJCQk8a2V5PlBheWxvYWRWZXJzaW9uPC9rZXk+CgkJCQk8aW50ZWdlcj4xPC9pbnRlZ2VyPgoJCQkJPGtleT5jbG9zZVZpZXdGYXJQb2ludDwva2V5PgoJCQkJPGludGVnZXI+MTwvaW50ZWdlcj4KCQkJCTxrZXk+Y2xvc2VWaWV3SG90a2V5c0VuYWJsZWQ8L2tleT4KCQkJCTxmYWxzZS8+CgkJCQk8a2V5PmNsb3NlVmlld05lYXJQb2ludDwva2V5PgoJCQkJPGludGVnZXI+MTA8L2ludGVnZXI+CgkJCQk8a2V5PmNsb3NlVmlld1Njcm9sbFdoZWVsVG9nZ2xlPC9rZXk+CgkJCQk8ZmFsc2UvPgoJCQkJPGtleT5jbG9zZVZpZXdTaG93UHJldmlldzwva2V5PgoJCQkJPGZhbHNlLz4KCQkJCTxrZXk+Y2xvc2VWaWV3U21vb3RoSW1hZ2VzPC9rZXk+CgkJCQk8ZmFsc2UvPgoJCQkJPGtleT5jb250cmFzdDwva2V5PgoJCQkJPGludGVnZXI+MDwvaW50ZWdlcj4KCQkJCTxrZXk+Zmxhc2hTY3JlZW48L2tleT4KCQkJCTx0cnVlLz4KCQkJCTxrZXk+Z3JheXNjYWxlPC9rZXk+CgkJCQk8ZmFsc2UvPgoJCQkJPGtleT5tb3VzZURyaXZlcjwva2V5PgoJCQkJPGZhbHNlLz4KCQkJCTxrZXk+bW91c2VEcml2ZXJDdXJzb3JTaXplPC9rZXk+CgkJCQk8aW50ZWdlcj4xPC9pbnRlZ2VyPgoJCQkJPGtleT5tb3VzZURyaXZlcklnbm9yZVRyYWNrcGFkPC9rZXk+CgkJCQk8ZmFsc2UvPgoJCQkJPGtleT5tb3VzZURyaXZlckluaXRpYWxEZWxheTwva2V5PgoJCQkJPHJlYWw+MTwvcmVhbD4KCQkJCTxrZXk+bW91c2VEcml2ZXJNYXhTcGVlZDwva2V5PgoJCQkJPGludGVnZXI+MzwvaW50ZWdlcj4KCQkJCTxrZXk+c2xvd0tleTwva2V5PgoJCQkJPGZhbHNlLz4KCQkJCTxrZXk+c2xvd0tleUJlZXBPbjwva2V5PgoJCQkJPGZhbHNlLz4KCQkJCTxrZXk+c2xvd0tleURlbGF5PC9rZXk+CgkJCQk8aW50ZWdlcj4wPC9pbnRlZ2VyPgoJCQkJPGtleT5zdGVyZW9Bc01vbm88L2tleT4KCQkJCTxmYWxzZS8+CgkJCQk8a2V5PnN0aWNreUtleTwva2V5PgoJCQkJPGZhbHNlLz4KCQkJCTxrZXk+c3RpY2t5S2V5QmVlcE9uTW9kaWZpZXI8L2tleT4KCQkJCTxmYWxzZS8+CgkJCQk8a2V5PnN0aWNreUtleVNob3dXaW5kb3c8L2tleT4KCQkJCTxmYWxzZS8+CgkJCQk8a2V5PnZvaWNlT3Zlck9uT2ZmS2V5PC9rZXk+CgkJCQk8ZmFsc2UvPgoJCQkJPGtleT53aGl0ZU9uQmxhY2s8L2tleT4KCQkJCTxmYWxzZS8+CgkJCTwvZGljdD4KCQk8L2FycmF5PgoJCTxrZXk+UGF5bG9hZERlc2NyaXB0aW9uPC9rZXk+CgkJPHN0cmluZy8+CgkJPGtleT5QYXlsb2FkRGlzcGxheU5hbWU8L2tleT4KCQk8c3RyaW5nPmR0LW1jcC1hY2Nlc3NpYmlsaXR5X2hlYXJpbmdfYmFzZS0wLjAuMS1wcm9kLWV1LTAtMDwvc3RyaW5nPgoJCTxrZXk+UGF5bG9hZEVuYWJsZWQ8L2tleT4KCQk8dHJ1ZS8+CgkJPGtleT5QYXlsb2FkSWRlbnRpZmllcjwva2V5PgoJCTxzdHJpbmc+YTVkY2Y3YTQtNzJiMi00NmYyLWFhMmMtNWNkMjQ3Zjg4ZWE5PC9zdHJpbmc+CgkJPGtleT5QYXlsb2FkT3JnYW5pemF0aW9uPC9rZXk+CgkJPHN0cmluZz5EZXBsb3ltZW50IFRoZW9yeTwvc3RyaW5nPgoJCTxrZXk+UGF5bG9hZFJlbW92YWxEaXNhbGxvd2VkPC9rZXk+CgkJPHRydWUvPgoJCTxrZXk+UGF5bG9hZFNjb3BlPC9rZXk+CgkJPHN0cmluZz5TeXN0ZW08L3N0cmluZz4KCQk8a2V5PlBheWxvYWRUeXBlPC9rZXk+CgkJPHN0cmluZz5Db25maWd1cmF0aW9uPC9zdHJpbmc+CgkJPGtleT5QYXlsb2FkVVVJRDwva2V5PgoJCTxzdHJpbmc+YTVkY2Y3YTQtNzJiMi00NmYyLWFhMmMtNWNkMjQ3Zjg4ZWE5PC9zdHJpbmc+CgkJPGtleT5QYXlsb2FkVmVyc2lvbjwva2V5PgoJCTxpbnRlZ2VyPjE8L2ludGVnZXI+Cgk8L2RpY3Q+CjwvcGxpc3Q+",
    "deploymentChannel": "deviceChannel"
}

// plist files

Request URL
https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations
Request Method
POST

request body

{
  "id":"00000000-0000-0000-0000-000000000000",
  "displayName":"test-list-file",
  "description":"test",
  "roleScopeTagIds":["0"],
  "@odata.type":"#microsoft.graph.macOSCustomAppConfiguration",
  "fileName":"dt-mcp-accessibility_hearing_base-prod-v0.0.1.plist","configurationXml":"PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPCFET0NUWVBFIHBsaXN0IFBVQkxJQyAiLS8vQXBwbGUvL0RURCBQTElTVCAxLjAvL0VOIiAiaHR0cDovL3d3dy5hcHBsZS5jb20vRFREcy9Qcm9wZXJ0eUxpc3QtMS4wLmR0ZCI+CjxwbGlzdCB2ZXJzaW9uPSIxLjAiPgoJPGRpY3Q+CgkJPGtleT5QYXlsb2FkQ29udGVudDwva2V5PgoJCTxhcnJheT4KCQkJPGRpY3Q+CgkJCQk8a2V5PlBheWxvYWREZXNjcmlwdGlvbjwva2V5PgoJCQkJPHN0cmluZy8+CgkJCQk8a2V5PlBheWxvYWREaXNwbGF5TmFtZTwva2V5PgoJCQkJPHN0cmluZz5BY2Nlc3NpYmlsaXR5PC9zdHJpbmc+CgkJCQk8a2V5PlBheWxvYWRFbmFibGVkPC9rZXk+CgkJCQk8dHJ1ZS8+CgkJCQk8a2V5PlBheWxvYWRJZGVudGlmaWVyPC9rZXk+CgkJCQk8c3RyaW5nPjVGNDlBRTVGLTU4ODQtNEUxQi04Nzc5LTFBNDdDMzg3MDdBQzwvc3RyaW5nPgoJCQkJPGtleT5QYXlsb2FkT3JnYW5pemF0aW9uPC9rZXk+CgkJCQk8c3RyaW5nPkRlcGxveW1lbnQgVGhlb3J5PC9zdHJpbmc+CgkJCQk8a2V5PlBheWxvYWRUeXBlPC9rZXk+CgkJCQk8c3RyaW5nPmNvbS5hcHBsZS51bml2ZXJzYWxhY2Nlc3M8L3N0cmluZz4KCQkJCTxrZXk+UGF5bG9hZFVVSUQ8L2tleT4KCQkJCTxzdHJpbmc+NUY0OUFFNUYtNTg4NC00RTFCLTg3NzktMUE0N0MzODcwN0FDPC9zdHJpbmc+CgkJCQk8a2V5PlBheWxvYWRWZXJzaW9uPC9rZXk+CgkJCQk8aW50ZWdlcj4xPC9pbnRlZ2VyPgoJCQkJPGtleT5jbG9zZVZpZXdGYXJQb2ludDwva2V5PgoJCQkJPGludGVnZXI+MTwvaW50ZWdlcj4KCQkJCTxrZXk+Y2xvc2VWaWV3SG90a2V5c0VuYWJsZWQ8L2tleT4KCQkJCTxmYWxzZS8+CgkJCQk8a2V5PmNsb3NlVmlld05lYXJQb2ludDwva2V5PgoJCQkJPGludGVnZXI+MTA8L2ludGVnZXI+CgkJCQk8a2V5PmNsb3NlVmlld1Njcm9sbFdoZWVsVG9nZ2xlPC9rZXk+CgkJCQk8ZmFsc2UvPgoJCQkJPGtleT5jbG9zZVZpZXdTaG93UHJldmlldzwva2V5PgoJCQkJPGZhbHNlLz4KCQkJCTxrZXk+Y2xvc2VWaWV3U21vb3RoSW1hZ2VzPC9rZXk+CgkJCQk8ZmFsc2UvPgoJCQkJPGtleT5jb250cmFzdDwva2V5PgoJCQkJPGludGVnZXI+MDwvaW50ZWdlcj4KCQkJCTxrZXk+Zmxhc2hTY3JlZW48L2tleT4KCQkJCTx0cnVlLz4KCQkJCTxrZXk+Z3JheXNjYWxlPC9rZXk+CgkJCQk8ZmFsc2UvPgoJCQkJPGtleT5tb3VzZURyaXZlcjwva2V5PgoJCQkJPGZhbHNlLz4KCQkJCTxrZXk+bW91c2VEcml2ZXJDdXJzb3JTaXplPC9rZXk+CgkJCQk8aW50ZWdlcj4xPC9pbnRlZ2VyPgoJCQkJPGtleT5tb3VzZURyaXZlcklnbm9yZVRyYWNrcGFkPC9rZXk+CgkJCQk8ZmFsc2UvPgoJCQkJPGtleT5tb3VzZURyaXZlckluaXRpYWxEZWxheTwva2V5PgoJCQkJPHJlYWw+MTwvcmVhbD4KCQkJCTxrZXk+bW91c2VEcml2ZXJNYXhTcGVlZDwva2V5PgoJCQkJPGludGVnZXI+MzwvaW50ZWdlcj4KCQkJCTxrZXk+c2xvd0tleTwva2V5PgoJCQkJPGZhbHNlLz4KCQkJCTxrZXk+c2xvd0tleUJlZXBPbjwva2V5PgoJCQkJPGZhbHNlLz4KCQkJCTxrZXk+c2xvd0tleURlbGF5PC9rZXk+CgkJCQk8aW50ZWdlcj4wPC9pbnRlZ2VyPgoJCQkJPGtleT5zdGVyZW9Bc01vbm88L2tleT4KCQkJCTxmYWxzZS8+CgkJCQk8a2V5PnN0aWNreUtleTwva2V5PgoJCQkJPGZhbHNlLz4KCQkJCTxrZXk+c3RpY2t5S2V5QmVlcE9uTW9kaWZpZXI8L2tleT4KCQkJCTxmYWxzZS8+CgkJCQk8a2V5PnN0aWNreUtleVNob3dXaW5kb3c8L2tleT4KCQkJCTxmYWxzZS8+CgkJCQk8a2V5PnZvaWNlT3Zlck9uT2ZmS2V5PC9rZXk+CgkJCQk8ZmFsc2UvPgoJCQkJPGtleT53aGl0ZU9uQmxhY2s8L2tleT4KCQkJCTxmYWxzZS8+CgkJCTwvZGljdD4KCQk8L2FycmF5PgoJCTxrZXk+UGF5bG9hZERlc2NyaXB0aW9uPC9rZXk+CgkJPHN0cmluZy8+CgkJPGtleT5QYXlsb2FkRGlzcGxheU5hbWU8L2tleT4KCQk8c3RyaW5nPmR0LW1jcC1hY2Nlc3NpYmlsaXR5X2hlYXJpbmdfYmFzZS0wLjAuMS1wcm9kLWV1LTAtMDwvc3RyaW5nPgoJCTxrZXk+UGF5bG9hZEVuYWJsZWQ8L2tleT4KCQk8dHJ1ZS8+CgkJPGtleT5QYXlsb2FkSWRlbnRpZmllcjwva2V5PgoJCTxzdHJpbmc+YTVkY2Y3YTQtNzJiMi00NmYyLWFhMmMtNWNkMjQ3Zjg4ZWE5PC9zdHJpbmc+CgkJPGtleT5QYXlsb2FkT3JnYW5pemF0aW9uPC9rZXk+CgkJPHN0cmluZz5EZXBsb3ltZW50IFRoZW9yeTwvc3RyaW5nPgoJCTxrZXk+UGF5bG9hZFJlbW92YWxEaXNhbGxvd2VkPC9rZXk+CgkJPHRydWUvPgoJCTxrZXk+UGF5bG9hZFNjb3BlPC9rZXk+CgkJPHN0cmluZz5TeXN0ZW08L3N0cmluZz4KCQk8a2V5PlBheWxvYWRUeXBlPC9rZXk+CgkJPHN0cmluZz5Db25maWd1cmF0aW9uPC9zdHJpbmc+CgkJPGtleT5QYXlsb2FkVVVJRDwva2V5PgoJCTxzdHJpbmc+YTVkY2Y3YTQtNzJiMi00NmYyLWFhMmMtNWNkMjQ3Zjg4ZWE5PC9zdHJpbmc+CgkJPGtleT5QYXlsb2FkVmVyc2lvbjwva2V5PgoJCTxpbnRlZ2VyPjE8L2ludGVnZXI+Cgk8L2RpY3Q+CjwvcGxpc3Q+","bundleId":"com.domain.thing"}

response

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceConfigurations/$entity",
    "@odata.type": "#microsoft.graph.macOSCustomAppConfiguration",
    "id": "ff80b04d-8b99-4fee-a22d-42bb4fa56070",
    "lastModifiedDateTime": "2025-08-19T07:58:31.8527423Z",
    "roleScopeTagIds": [
        "0"
    ],
    "supportsScopeTags": true,
    "deviceManagementApplicabilityRuleOsEdition": null,
    "deviceManagementApplicabilityRuleOsVersion": null,
    "deviceManagementApplicabilityRuleDeviceMode": null,
    "createdDateTime": "2025-08-19T07:58:31.8527423Z",
    "description": "test",
    "displayName": "test-list-file",
    "version": 1,
    "bundleId": "com.domain.thing",
    "fileName": "dt-mcp-accessibility_hearing_base-prod-v0.0.1.plist",
    "configurationXml": "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPCFET0NUWVBFIHBsaXN0IFBVQkxJQyAiLS8vQXBwbGUvL0RURCBQTElTVCAxLjAvL0VOIiAiaHR0cDovL3d3dy5hcHBsZS5jb20vRFREcy9Qcm9wZXJ0eUxpc3QtMS4wLmR0ZCI+CjxwbGlzdCB2ZXJzaW9uPSIxLjAiPgoJPGRpY3Q+CgkJPGtleT5QYXlsb2FkQ29udGVudDwva2V5PgoJCTxhcnJheT4KCQkJPGRpY3Q+CgkJCQk8a2V5PlBheWxvYWREZXNjcmlwdGlvbjwva2V5PgoJCQkJPHN0cmluZy8+CgkJCQk8a2V5PlBheWxvYWREaXNwbGF5TmFtZTwva2V5PgoJCQkJPHN0cmluZz5BY2Nlc3NpYmlsaXR5PC9zdHJpbmc+CgkJCQk8a2V5PlBheWxvYWRFbmFibGVkPC9rZXk+CgkJCQk8dHJ1ZS8+CgkJCQk8a2V5PlBheWxvYWRJZGVudGlmaWVyPC9rZXk+CgkJCQk8c3RyaW5nPjVGNDlBRTVGLTU4ODQtNEUxQi04Nzc5LTFBNDdDMzg3MDdBQzwvc3RyaW5nPgoJCQkJPGtleT5QYXlsb2FkT3JnYW5pemF0aW9uPC9rZXk+CgkJCQk8c3RyaW5nPkRlcGxveW1lbnQgVGhlb3J5PC9zdHJpbmc+CgkJCQk8a2V5PlBheWxvYWRUeXBlPC9rZXk+CgkJCQk8c3RyaW5nPmNvbS5hcHBsZS51bml2ZXJzYWxhY2Nlc3M8L3N0cmluZz4KCQkJCTxrZXk+UGF5bG9hZFVVSUQ8L2tleT4KCQkJCTxzdHJpbmc+NUY0OUFFNUYtNTg4NC00RTFCLTg3NzktMUE0N0MzODcwN0FDPC9zdHJpbmc+CgkJCQk8a2V5PlBheWxvYWRWZXJzaW9uPC9rZXk+CgkJCQk8aW50ZWdlcj4xPC9pbnRlZ2VyPgoJCQkJPGtleT5jbG9zZVZpZXdGYXJQb2ludDwva2V5PgoJCQkJPGludGVnZXI+MTwvaW50ZWdlcj4KCQkJCTxrZXk+Y2xvc2VWaWV3SG90a2V5c0VuYWJsZWQ8L2tleT4KCQkJCTxmYWxzZS8+CgkJCQk8a2V5PmNsb3NlVmlld05lYXJQb2ludDwva2V5PgoJCQkJPGludGVnZXI+MTA8L2ludGVnZXI+CgkJCQk8a2V5PmNsb3NlVmlld1Njcm9sbFdoZWVsVG9nZ2xlPC9rZXk+CgkJCQk8ZmFsc2UvPgoJCQkJPGtleT5jbG9zZVZpZXdTaG93UHJldmlldzwva2V5PgoJCQkJPGZhbHNlLz4KCQkJCTxrZXk+Y2xvc2VWaWV3U21vb3RoSW1hZ2VzPC9rZXk+CgkJCQk8ZmFsc2UvPgoJCQkJPGtleT5jb250cmFzdDwva2V5PgoJCQkJPGludGVnZXI+MDwvaW50ZWdlcj4KCQkJCTxrZXk+Zmxhc2hTY3JlZW48L2tleT4KCQkJCTx0cnVlLz4KCQkJCTxrZXk+Z3JheXNjYWxlPC9rZXk+CgkJCQk8ZmFsc2UvPgoJCQkJPGtleT5tb3VzZURyaXZlcjwva2V5PgoJCQkJPGZhbHNlLz4KCQkJCTxrZXk+bW91c2VEcml2ZXJDdXJzb3JTaXplPC9rZXk+CgkJCQk8aW50ZWdlcj4xPC9pbnRlZ2VyPgoJCQkJPGtleT5tb3VzZURyaXZlcklnbm9yZVRyYWNrcGFkPC9rZXk+CgkJCQk8ZmFsc2UvPgoJCQkJPGtleT5tb3VzZURyaXZlckluaXRpYWxEZWxheTwva2V5PgoJCQkJPHJlYWw+MTwvcmVhbD4KCQkJCTxrZXk+bW91c2VEcml2ZXJNYXhTcGVlZDwva2V5PgoJCQkJPGludGVnZXI+MzwvaW50ZWdlcj4KCQkJCTxrZXk+c2xvd0tleTwva2V5PgoJCQkJPGZhbHNlLz4KCQkJCTxrZXk+c2xvd0tleUJlZXBPbjwva2V5PgoJCQkJPGZhbHNlLz4KCQkJCTxrZXk+c2xvd0tleURlbGF5PC9rZXk+CgkJCQk8aW50ZWdlcj4wPC9pbnRlZ2VyPgoJCQkJPGtleT5zdGVyZW9Bc01vbm88L2tleT4KCQkJCTxmYWxzZS8+CgkJCQk8a2V5PnN0aWNreUtleTwva2V5PgoJCQkJPGZhbHNlLz4KCQkJCTxrZXk+c3RpY2t5S2V5QmVlcE9uTW9kaWZpZXI8L2tleT4KCQkJCTxmYWxzZS8+CgkJCQk8a2V5PnN0aWNreUtleVNob3dXaW5kb3c8L2tleT4KCQkJCTxmYWxzZS8+CgkJCQk8a2V5PnZvaWNlT3Zlck9uT2ZmS2V5PC9rZXk+CgkJCQk8ZmFsc2UvPgoJCQkJPGtleT53aGl0ZU9uQmxhY2s8L2tleT4KCQkJCTxmYWxzZS8+CgkJCTwvZGljdD4KCQk8L2FycmF5PgoJCTxrZXk+UGF5bG9hZERlc2NyaXB0aW9uPC9rZXk+CgkJPHN0cmluZy8+CgkJPGtleT5QYXlsb2FkRGlzcGxheU5hbWU8L2tleT4KCQk8c3RyaW5nPmR0LW1jcC1hY2Nlc3NpYmlsaXR5X2hlYXJpbmdfYmFzZS0wLjAuMS1wcm9kLWV1LTAtMDwvc3RyaW5nPgoJCTxrZXk+UGF5bG9hZEVuYWJsZWQ8L2tleT4KCQk8dHJ1ZS8+CgkJPGtleT5QYXlsb2FkSWRlbnRpZmllcjwva2V5PgoJCTxzdHJpbmc+YTVkY2Y3YTQtNzJiMi00NmYyLWFhMmMtNWNkMjQ3Zjg4ZWE5PC9zdHJpbmc+CgkJPGtleT5QYXlsb2FkT3JnYW5pemF0aW9uPC9rZXk+CgkJPHN0cmluZz5EZXBsb3ltZW50IFRoZW9yeTwvc3RyaW5nPgoJCTxrZXk+UGF5bG9hZFJlbW92YWxEaXNhbGxvd2VkPC9rZXk+CgkJPHRydWUvPgoJCTxrZXk+UGF5bG9hZFNjb3BlPC9rZXk+CgkJPHN0cmluZz5TeXN0ZW08L3N0cmluZz4KCQk8a2V5PlBheWxvYWRUeXBlPC9rZXk+CgkJPHN0cmluZz5Db25maWd1cmF0aW9uPC9zdHJpbmc+CgkJPGtleT5QYXlsb2FkVVVJRDwva2V5PgoJCTxzdHJpbmc+YTVkY2Y3YTQtNzJiMi00NmYyLWFhMmMtNWNkMjQ3Zjg4ZWE5PC9zdHJpbmc+CgkJPGtleT5QYXlsb2FkVmVyc2lvbjwva2V5PgoJCTxpbnRlZ2VyPjE8L2ludGVnZXI+Cgk8L2RpY3Q+CjwvcGxpc3Q+"
}

// trusted certificate (root certs)

Request URL
https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations
Request Method
POST

request

{"id":"00000000-0000-0000-0000-000000000000","displayName":"test","description":"test","roleScopeTagIds":["0"],"@odata.type":"#microsoft.graph.macOSTrustedRootCertificate","deploymentChannel":"deviceChannel","certFileName":"MicrosoftRootCertificateAuthority2011.cer","trustedRootCertificate":"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"}

response

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceConfigurations/$entity",
    "@odata.type": "#microsoft.graph.macOSTrustedRootCertificate",
    "id": "afeb20c3-48a0-4bbd-8191-c9c9f2fd62d2",
    "lastModifiedDateTime": "2025-08-19T08:09:39.5333342Z",
    "roleScopeTagIds": [
        "0"
    ],
    "supportsScopeTags": true,
    "deviceManagementApplicabilityRuleOsEdition": null,
    "deviceManagementApplicabilityRuleOsVersion": null,
    "deviceManagementApplicabilityRuleDeviceMode": null,
    "createdDateTime": "2025-08-19T08:09:39.5333342Z",
    "description": "test",
    "displayName": "test",
    "version": 1,
    "trustedRootCertificate": "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",
    "certFileName": "MicrosoftRootCertificateAuthority2011.cer",
    "deploymentChannel": "deviceChannel"
}

// scep cert

Request URL
https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations
Request Method
POST

request

{
  "id":"00000000-0000-0000-0000-000000000000","displayName":"scep-cert",
  "description":"scep-cert",
  "roleScopeTagIds":["0"],
  "@odata.type":"#microsoft.graph.macOSScepCertificateProfile","renewalThresholdPercentage":20,"deploymentChannel":"deviceChannel","certificateStore":"machine","certificateValidityPeriodScale":"years","certificateValidityPeriodValue":1,"subjectNameFormat":"custom","subjectNameFormatString":"CN={{AAD_Device_ID}}","rootCertificate@odata.bind":"https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('afeb20c3-48a0-4bbd-8191-c9c9f2fd62d2')","keySize":"size4096","keyUsage":"digitalSignature,keyEncipherment","customSubjectAlternativeNames":[{"sanType":"emailAddress","name":"domain.com"},{"sanType":"userPrincipalName","name":"some-upn"},{"sanType":"domainNameService","name":"some-dns-record"},{"sanType":"universalResourceIdentifier","name":"some-uri"}],"extendedKeyUsages":[{"name":"Any Purpose","objectIdentifier":"2.5.29.37.0"},{"name":"Client Authentication","objectIdentifier":"1.3.6.1.5.5.7.3.2"},{"name":"Secure Email","objectIdentifier":"1.3.6.1.5.5.7.3.4"}],"scepServerUrls":["https://scep-url.com","https://scep-url-2.com"],"allowAllAppsAccess":true}

response

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceConfigurations/$entity",
    "@odata.type": "#microsoft.graph.macOSScepCertificateProfile",
    "id": "df301cbc-a3e9-4369-8613-291ce7292a85",
    "lastModifiedDateTime": "2025-08-19T08:13:19.0484998Z",
    "roleScopeTagIds": [
        "0"
    ],
    "supportsScopeTags": true,
    "deviceManagementApplicabilityRuleOsEdition": null,
    "deviceManagementApplicabilityRuleOsVersion": null,
    "deviceManagementApplicabilityRuleDeviceMode": null,
    "createdDateTime": "2025-08-19T08:13:19.0484998Z",
    "description": "scep-cert",
    "displayName": "scep-cert",
    "version": 1,
    "renewalThresholdPercentage": 20,
    "subjectNameFormat": "custom",
    "subjectAlternativeNameType": null,
    "certificateValidityPeriodValue": 1,
    "certificateValidityPeriodScale": "years",
    "scepServerUrls": [
        "https://scep-url.com",
        "https://scep-url-2.com"
    ],
    "subjectNameFormatString": "CN={{AAD_Device_ID}}",
    "keyUsage": "keyEncipherment,digitalSignature",
    "keySize": "size4096",
    "hashAlgorithm": null,
    "subjectAlternativeNameFormatString": null,
    "certificateStore": "machine",
    "allowAllAppsAccess": true,
    "deploymentChannel": "deviceChannel",
    "extendedKeyUsages": [
        {
            "name": "Any Purpose",
            "objectIdentifier": "2.5.29.37.0"
        },
        {
            "name": "Client Authentication",
            "objectIdentifier": "1.3.6.1.5.5.7.3.2"
        },
        {
            "name": "Secure Email",
            "objectIdentifier": "1.3.6.1.5.5.7.3.4"
        }
    ],
    "customSubjectAlternativeNames": [
        {
            "sanType": "emailAddress",
            "name": "domain.com"
        },
        {
            "sanType": "userPrincipalName",
            "name": "some-upn"
        },
        {
            "sanType": "domainNameService",
            "name": "some-dns-record"
        },
        {
            "sanType": "universalResourceIdentifier",
            "name": "some-uri"
        }
    ]
}

// pkcs cert

Request URL
https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations
Request Method
POST

request 
{"id":"00000000-0000-0000-0000-000000000000","displayName":"pkcs-cert","description":"pkcs-cert","roleScopeTagIds":["0"],"@odata.type":"#microsoft.graph.macOSPkcsCertificateProfile","renewalThresholdPercentage":20,"deploymentChannel":"deviceChannel","certificateStore":"machine","certificateValidityPeriodScale":"years","certificateValidityPeriodValue":1,"subjectNameFormat":"custom","subjectNameFormatString":"CN={{AAD_Device_ID}}","certificationAuthority":"some-auth","certificationAuthorityName":"some-name","certificateTemplateName":"some-template-name","customSubjectAlternativeNames":[{"sanType":"emailAddress","name":"some-email"},{"sanType":"userPrincipalName","name":"some-upn"},{"sanType":"domainNameService","name":"some-dns"},{"sanType":"universalResourceIdentifier","name":"some-uri"},{"sanType":"customAzureADAttribute","name":"some-custom-att"},{"sanType":"emailAddress","name":"some-other-email"}],"allowAllAppsAccess":true}

response

{
    "@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceConfigurations/$entity",
    "@odata.type": "#microsoft.graph.macOSPkcsCertificateProfile",
    "id": "f2ae2670-9b33-4ca4-b616-fb4f38d9a5d0",
    "lastModifiedDateTime": "2025-08-19T08:16:54.4905069Z",
    "roleScopeTagIds": [
        "0"
    ],
    "supportsScopeTags": true,
    "deviceManagementApplicabilityRuleOsEdition": null,
    "deviceManagementApplicabilityRuleOsVersion": null,
    "deviceManagementApplicabilityRuleDeviceMode": null,
    "createdDateTime": "2025-08-19T08:16:54.4905069Z",
    "description": "pkcs-cert",
    "displayName": "pkcs-cert",
    "version": 1,
    "renewalThresholdPercentage": 20,
    "subjectNameFormat": "custom",
    "subjectAlternativeNameType": null,
    "certificateValidityPeriodValue": 1,
    "certificateValidityPeriodScale": "years",
    "certificationAuthority": "some-auth",
    "certificationAuthorityName": "some-name",
    "certificateTemplateName": "some-template-name",
    "subjectAlternativeNameFormatString": null,
    "subjectNameFormatString": "CN={{AAD_Device_ID}}",
    "certificateStore": "machine",
    "allowAllAppsAccess": true,
    "deploymentChannel": "deviceChannel",
    "customSubjectAlternativeNames": [
        {
            "sanType": "emailAddress",
            "name": "some-email"
        },
        {
            "sanType": "userPrincipalName",
            "name": "some-upn"
        },
        {
            "sanType": "domainNameService",
            "name": "some-dns"
        },
        {
            "sanType": "universalResourceIdentifier",
            "name": "some-uri"
        },
        {
            "sanType": "customAzureADAttribute",
            "name": "some-custom-att"
        },
        {
            "sanType": "emailAddress",
            "name": "some-other-email"
        }
    ]
}