package sample

// constructResource is a standard function and is not reported.
func constructResource() string {
	return settings()
}

// constructWidgetRequest meets every part of the helper rule.
func constructWidgetRequest(id string) map[string]string {
	body := map[string]string{}
	body["id"] = id
	body["kind"] = "widget"
	body["state"] = "enabled"
	body["owner"] = "fixture"
	body["region"] = "uk"
	body["tier"] = "standard"
	body["version"] = "1"
	return body
}

// settings is a banned generic name, only called from a non-CRUD function.
func settings() string {
	return "settings"
}
