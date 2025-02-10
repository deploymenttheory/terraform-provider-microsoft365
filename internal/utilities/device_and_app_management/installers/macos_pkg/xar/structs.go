package xar

type InstallerMetadata struct {
	// Primary app info
	Name             string
	Version          string
	BundleIdentifier string
	PackageIDs       []string
	SHASum           []byte

	// Additional bundles found in pkg
	IncludedBundles []BundleInfo

	// Installation paths
	InstallLocation string
	AppPaths        []string

	// Optional metadata
	MinOSVersion string
	Description  string
	Size         int64
}

type BundleInfo struct {
	BundleID        string
	Version         string
	Path            string
	CFBundleVersion string
}
