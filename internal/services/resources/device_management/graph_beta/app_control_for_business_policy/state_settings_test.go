package graphBetaAppControlForBusinessPolicy

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TestAppControlPolicyXMLRoundTrip verifies XML survives construction and read-back without alteration.
func TestAppControlPolicyXMLRoundTrip(t *testing.T) {
	tests := map[string]string{
		"LF":       "<?xml version=\"1.0\"?>\n<SiPolicy />\n",
		"CRLF":     "<?xml version=\"1.0\"?>\r\n<SiPolicy />\r\n",
		"BOM":      "\ufeff<?xml version=\"1.0\"?>\n<SiPolicy />\n",
		"entities": "<?xml version=\"1.0\"?>\n<SiPolicy><Name>A &amp; B</Name></SiPolicy>\n",
	}
	for name, xml := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			configured := AppControlForBusinessPolicyResourceModel{
				PolicyXML: types.StringValue(xml),
			}
			child, err := constructXMLContentChild(ctx, &configured)
			if err != nil {
				t.Fatal(err)
			}

			// Graph rejects a leading BOM and preserves the rest of the submitted XML.
			// An empty prior state exercises import rather than retaining a configured value.
			imported := AppControlForBusinessPolicyResourceModel{PolicyXML: types.StringNull()}
			if err := extractPolicyXMLContent(ctx, &imported, child); err != nil {
				t.Fatal(err)
			}
			if want, got := strings.TrimPrefix(
				xml,
				"\ufeff",
			), imported.PolicyXML.ValueString(); got != want {
				t.Fatalf("XML changed during import:\nwant: %q\n got: %q", want, got)
			}
			if err := extractPolicyXMLContent(ctx, &configured, child); err != nil {
				t.Fatal(err)
			}
			if got := configured.PolicyXML.ValueString(); got != xml {
				t.Fatalf("XML changed during refresh:\nwant: %q\n got: %q", xml, got)
			}
		})
	}
}

// TestAppControlPolicyXMLRemoteChange verifies remote XML changes replace the prior state.
func TestAppControlPolicyXMLRemoteChange(t *testing.T) {
	ctx := context.Background()
	remoteXML := "<?xml version=\"1.0\"?>\n<SiPolicy><Name>Changed &amp; Current</Name></SiPolicy>\n"
	remote := AppControlForBusinessPolicyResourceModel{PolicyXML: types.StringValue(remoteXML)}
	child, err := constructXMLContentChild(ctx, &remote)
	if err != nil {
		t.Fatal(err)
	}
	prior := AppControlForBusinessPolicyResourceModel{
		PolicyXML: types.StringValue(
			"<?xml version=\"1.0\"?><SiPolicy><Name>Old</Name></SiPolicy>",
		),
	}
	if err := extractPolicyXMLContent(ctx, &prior, child); err != nil {
		t.Fatal(err)
	}
	if got := prior.PolicyXML.ValueString(); got != remoteXML {
		t.Fatalf("remote XML change was hidden: want %q, got %q", remoteXML, got)
	}
}
