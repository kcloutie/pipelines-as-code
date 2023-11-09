//go:build e2e
// +build e2e

package test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	tgithub "github.com/openshift-pipelines/pipelines-as-code/test/pkg/github"
	"github.com/openshift-pipelines/pipelines-as-code/test/pkg/options"
	"github.com/openshift-pipelines/pipelines-as-code/test/pkg/payload"
	"github.com/openshift-pipelines/pipelines-as-code/test/pkg/wait"
	"github.com/tektoncd/pipeline/pkg/names"
	"gotest.tools/v3/assert"
)

func TestGithubChangedFiles(t *testing.T) {
	ctx := context.Background()
	label := "Github Changed Files"
	maxNumberOfConcurrentPipelineRuns := 1

	targetNS := names.SimpleNameGenerator.RestrictLengthWithRandomSuffix("pac-e2e-ns")
	runcnx, opts, ghcnx, err := tgithub.Setup(ctx, false)
	assert.NilError(t, err)

	logmsg := fmt.Sprintf("Testing %s with Github APPS integration on %s", label, targetNS)
	runcnx.Clients.Log.Info(logmsg)

	repoinfo, resp, err := ghcnx.Client.Repositories.Get(ctx, opts.Organization, opts.Repo)
	assert.NilError(t, err)
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		t.Errorf("Repository %s not found in %s", opts.Organization, opts.Repo)
	}

	// set concurrency
	opts.Concurrency = maxNumberOfConcurrentPipelineRuns

	err = tgithub.CreateCRD(ctx, t, repoinfo, runcnx, opts, targetNS)
	assert.NilError(t, err)

	yamlFiles := map[string]string{
		".tekton/pullrequest.yaml": "testdata/pipelinerun-changed-files-pullrequest.yaml",
		".tekton/push.yaml":        "testdata/pipelinerun-changed-files-push.yaml",
		"deleted.txt":              "testdata/changed_files_deleted",
		"modified.txt":             "testdata/changed_files_modified",
		"renamed.txt":              "testdata/changed_files_renamed",
	}

	entries, err := payload.GetEntries(yamlFiles, targetNS, options.MainBranch, options.PullRequestEvent, map[string]string{})
	assert.NilError(t, err)

	targetRefName := fmt.Sprintf("refs/heads/%s",
		names.SimpleNameGenerator.RestrictLengthWithRandomSuffix("pac-e2e-test"))

	sha, err := tgithub.PushFilesToRef(ctx, ghcnx.Client, logmsg, repoinfo.GetDefaultBranch(), targetRefName,
		opts.Organization, opts.Repo, entries)
	assert.NilError(t, err)
	runcnx.Clients.Log.Infof("Commit %s has been created and pushed to %s", sha, targetRefName)

	prNumber, err := tgithub.PRCreate(ctx, runcnx, ghcnx, opts.Organization,
		opts.Repo, targetRefName, repoinfo.GetDefaultBranch(), logmsg)
	assert.NilError(t, err)
	defer tgithub.TearDown(ctx, t, runcnx, ghcnx, prNumber, targetRefName, targetNS, opts)

	runcnx.Clients.Log.Info("waiting to let controller process the event")
	time.Sleep(5 * time.Second)

	waitOpts := wait.Opts{
		RepoName:        targetNS,
		Namespace:       targetNS,
		MinNumberStatus: 1,
		PollTimeout:     wait.DefaultTimeout,
		TargetSHA:       sha,
	}
	assert.NilError(t, wait.UntilMinPRAppeared(ctx, runcnx.Clients, waitOpts, 1))

}
