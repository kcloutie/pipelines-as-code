package customparams

import (
	"context"
	"fmt"
	"strings"

	"github.com/openshift-pipelines/pipelines-as-code/pkg/formatting"
	"go.uber.org/zap"
)

func (p *CustomParams) getChangedFilesCommaSeparated(ctx context.Context) (string, string, string, string, string) {
	if p.vcx == nil {
		return "", "", "", "", ""
	}
	allChangedFiles, addedFiles, deletedFiles, modifiedFiles, renamedFiles, err := p.vcx.GetFiles(context.Background(), p.event)
	if err != nil {
		p.eventEmitter.EmitMessage(p.repo, zap.ErrorLevel, "ParamsError", fmt.Sprintf("error getting changed files: %s", err.Error()))
	}

	return strings.Join(uniqueStringArray(allChangedFiles), ","), strings.Join(uniqueStringArray(addedFiles), ","), strings.Join(uniqueStringArray(deletedFiles), ","), strings.Join(uniqueStringArray(modifiedFiles), ","), strings.Join(uniqueStringArray(renamedFiles), ",")
}

func uniqueStringArray(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// makeStandardParamsFromEvent will create a map of standard params out of the event
func (p *CustomParams) makeStandardParamsFromEvent(ctx context.Context) map[string]string {
	repoURL := p.event.URL
	// On bitbucket server you are have a special url for checking it out, they
	// seemed to fix it in 2.0 but i guess we have to live with this until then.
	if p.event.CloneURL != "" {
		repoURL = p.event.CloneURL
	}
	allChangedFiles, addedFiles, deletedFiles, modifiedFiles, renamedFiles := p.getChangedFilesCommaSeparated(ctx)
	return map[string]string{
		"revision":          p.event.SHA,
		"repo_url":          repoURL,
		"repo_owner":        strings.ToLower(p.event.Organization),
		"repo_name":         strings.ToLower(p.event.Repository),
		"target_branch":     formatting.SanitizeBranch(p.event.BaseBranch),
		"source_branch":     formatting.SanitizeBranch(p.event.HeadBranch),
		"sender":            strings.ToLower(p.event.Sender),
		"target_namespace":  p.repo.GetNamespace(),
		"event_type":        p.event.EventType,
		"all_changed_files": allChangedFiles,
		"added_files":       addedFiles,
		"deleted_files":     deletedFiles,
		"modified_files":    modifiedFiles,
		"renamed_files":     renamedFiles,
	}
}
