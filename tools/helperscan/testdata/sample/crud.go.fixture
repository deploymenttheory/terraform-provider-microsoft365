package sample

import "context"

// Create uses the construct helper, the method helper and the lock.
func (r *SampleResource) Create(ctx context.Context) {
	unlock := lockWidget(r.client)
	defer unlock()
	_ = constructWidgetRequest(r.client)
	r.refreshWidget()
}

// Read shadows the settings helper with a local and passes the Kiota factory as a value.
func (r *SampleResource) Read(ctx context.Context) {
	settings := 1
	_ = settings
	parse := createWidgetResponseFromDiscriminatorValue
	_ = parse
	_ = findWidgetByName(nil, "a")
}

// Update calls the helpers from inside a closure, which still counts as Update.
func (r *SampleResource) Update(ctx context.Context) {
	apply := func() {
		_ = constructWidgetRequest(r.client)
		r.refreshWidget()
		waitForWidgetDeletion(r.client)
	}
	apply()
}

// Delete waits for the widget to go.
func (r *SampleResource) Delete(ctx context.Context) {
	waitForWidgetDeletion(r.client)
}

// refreshWidget is a method helper.
func (r *SampleResource) refreshWidget() {
	r.client = "refreshed"
}

// lockWidget is short and its name has no qualifier.
func lockWidget(id string) func() {
	_ = id
	return func() {}
}

// waitForWidgetDeletion is long enough, called from two CRUD methods and well named.
func waitForWidgetDeletion(id string) {
	attempts := 0
	for attempts < 3 {
		attempts++
		if id == "" {
			return
		}
		if attempts > 5 {
			break
		}
	}
}

// mapWidgetsToState is in the wrong file for its verb.
func mapWidgetsToState(values []string) []string {
	return values
}
