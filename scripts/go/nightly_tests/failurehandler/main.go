package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/deploymenttheory/terraform-provider-microsoft365/scripts/go/nightly_tests/failurehandler/handler"
	"github.com/google/go-github/v66/github"
	"golang.org/x/oauth2"
)

func main() {
	cfg := parseFlags()

	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Configuration error: %v\n", err)
		flag.Usage()
		os.Exit(1)
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		fmt.Fprintln(os.Stderr, "Error: GITHUB_TOKEN environment variable is required")
		os.Exit(1)
	}

	ctx := context.Background()
	client := createGitHubClient(ctx, token)

	h := handler.New(client, cfg)

	if err := h.Handle(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error handling failure: %v\n", err)
		os.Exit(1)
	}
}

func parseFlags() *handler.Config {
	cfg := &handler.Config{}

	flag.StringVar(&cfg.Owner, "owner", "", "Repository owner")
	flag.StringVar(&cfg.Repo, "repo", "", "Repository name")
	flag.StringVar(&cfg.RunID, "run-id", "", "Workflow run ID")
	flag.StringVar(&cfg.FailedJobs, "failed-jobs", "", "Comma-separated list of failed jobs")
	flag.Parse()

	return cfg
}

func createGitHubClient(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}
