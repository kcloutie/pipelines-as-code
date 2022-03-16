package matcher

import (
	"context"

	apipac "github.com/openshift-pipelines/pipelines-as-code/pkg/apis/pipelinesascode/v1alpha1"
	"github.com/openshift-pipelines/pipelines-as-code/pkg/params"
	"github.com/openshift-pipelines/pipelines-as-code/pkg/params/info"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func MatchEventURLRepo(ctx context.Context, cs *params.Run, event *info.Event, ns string) (*apipac.Repository, error) {
	repositories, err := cs.Clients.PipelineAsCode.PipelinesascodeV1alpha1().Repositories(ns).List(
		ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for i := len(repositories.Items) - 1; i >= 0; i-- {
		repo := repositories.Items[i]
		if repo.Spec.URL == event.URL {
			return &repo, nil
		}
	}

	return nil, nil
}
