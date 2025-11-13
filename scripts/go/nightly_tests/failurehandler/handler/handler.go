package handler

import (
	"context"
	"fmt"

	"github.com/google/go-github/v66/github"
)

// Label names for issues and PRs.
const (
	LabelBug           = "bug"
	LabelTesting       = "testing"
	LabelAutomated     = "automated"
	LabelNightlyFailure = "nightly-failure"
)

var defaultLabels = []string{LabelBug, LabelTesting, LabelAutomated, LabelNightlyFailure}

// Handler handles the creation of PRs and issues for test failures.
type Handler struct {
	client *github.Client
	config *Config
}

// New creates a new failure handler.
func New(client *github.Client, config *Config) *Handler {
	return &Handler{
		client: client,
		config: config,
	}
}

// Handle performs the complete failure handling workflow.
func (h *Handler) Handle(ctx context.Context) error {
	fmt.Printf("Creating failure report for failed jobs: %s\n", h.config.FailedJobs)

	// Get default branch
	defaultBranch, err := h.getDefaultBranch(ctx)
	if err != nil {
		return fmt.Errorf("failed to get default branch: %w", err)
	}

	// Create new branch
	branchName := h.config.GetBranchName()
	if err := h.createBranch(ctx, branchName, defaultBranch); err != nil {
		return fmt.Errorf("failed to create branch: %w", err)
	}

	fmt.Printf("✅ Created branch: %s\n", branchName)

	// Create failure report file
	if err := h.createFailureReport(ctx, branchName); err != nil {
		return fmt.Errorf("failed to create failure report: %w", err)
	}

	fmt.Println("✅ Created FAILURE_REPORT.md")

	// Create pull request
	pr, err := h.createPullRequest(ctx, branchName, defaultBranch)
	if err != nil {
		return fmt.Errorf("failed to create pull request: %w", err)
	}

	fmt.Printf("✅ Created PR #%d: %s\n", pr.GetNumber(), pr.GetHTMLURL())

	// Create issue
	issue, err := h.createIssue(ctx, pr.GetHTMLURL())
	if err != nil {
		return fmt.Errorf("failed to create issue: %w", err)
	}

	fmt.Printf("✅ Created Issue #%d: %s\n", issue.GetNumber(), issue.GetHTMLURL())

	return nil
}

// getDefaultBranch retrieves the default branch name for the repository.
func (h *Handler) getDefaultBranch(ctx context.Context) (string, error) {
	repo, _, err := h.client.Repositories.Get(ctx, h.config.Owner, h.config.Repo)
	if err != nil {
		return "", err
	}
	return repo.GetDefaultBranch(), nil
}

// createBranch creates a new branch from the default branch.
func (h *Handler) createBranch(ctx context.Context, branchName, baseBranch string) error {
	// Get reference to base branch
	baseRef, _, err := h.client.Git.GetRef(ctx, h.config.Owner, h.config.Repo, "refs/heads/"+baseBranch)
	if err != nil {
		return err
	}

	// Create new branch reference
	newRef := &github.Reference{
		Ref: github.String("refs/heads/" + branchName),
		Object: &github.GitObject{
			SHA: baseRef.Object.SHA,
		},
	}

	_, _, err = h.client.Git.CreateRef(ctx, h.config.Owner, h.config.Repo, newRef)
	return err
}

// createFailureReport creates the FAILURE_REPORT.md file in the repository.
func (h *Handler) createFailureReport(ctx context.Context, branchName string) error {
	content := buildFailureReport(
		h.config.GetDate(),
		h.config.ParseFailedJobs(),
		h.config.GetWorkflowURL(),
		h.config.GetCodecovURL(),
	)

	commitMessage := buildCommitMessage(
		h.config.GetDate(),
		h.config.FailedJobs,
		h.config.GetWorkflowURL(),
	)

	fileOptions := &github.RepositoryContentFileOptions{
		Message: github.String(commitMessage),
		Content: []byte(content),
		Branch:  github.String(branchName),
	}

	_, _, err := h.client.Repositories.CreateFile(
		ctx,
		h.config.Owner,
		h.config.Repo,
		"FAILURE_REPORT.md",
		fileOptions,
	)

	return err
}

// createPullRequest creates a pull request for the failure report.
func (h *Handler) createPullRequest(ctx context.Context, headBranch, baseBranch string) (*github.PullRequest, error) {
	title := buildPRTitle(h.config.FailedJobs, h.config.GetDate())
	body := buildPRBody(
		h.config.GetDate(),
		h.config.FailedJobs,
		h.config.GetWorkflowURL(),
		h.config.GetCodecovURL(),
	)

	newPR := &github.NewPullRequest{
		Title: github.String(title),
		Head:  github.String(headBranch),
		Base:  github.String(baseBranch),
		Body:  github.String(body),
	}

	pr, _, err := h.client.PullRequests.Create(ctx, h.config.Owner, h.config.Repo, newPR)
	if err != nil {
		return nil, err
	}

	// Add labels to PR
	if _, _, err := h.client.Issues.AddLabelsToIssue(
		ctx,
		h.config.Owner,
		h.config.Repo,
		pr.GetNumber(),
		defaultLabels,
	); err != nil {
		fmt.Printf("⚠️  Warning: failed to add labels to PR: %v\n", err)
	}

	return pr, nil
}

// createIssue creates a tracking issue for the failure.
func (h *Handler) createIssue(ctx context.Context, prURL string) (*github.Issue, error) {
	title := buildIssueTitle(h.config.GetDate(), h.config.FailedJobs)
	body := buildIssueBody(
		h.config.GetDate(),
		h.config.FailedJobs,
		h.config.GetWorkflowURL(),
		prURL,
		h.config.GetCodecovURL(),
	)

	newIssue := &github.IssueRequest{
		Title:  github.String(title),
		Body:   github.String(body),
		Labels: &defaultLabels,
	}

	issue, _, err := h.client.Issues.Create(ctx, h.config.Owner, h.config.Repo, newIssue)
	return issue, err
}
