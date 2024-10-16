package job

import (
	"context"
	"fmt"

	"github.com/dictyBase/modware-import/internal/collection"
	"github.com/dictyBase/modware-import/internal/k8s/cli/parameters"
	"github.com/dictyBase/modware-import/internal/k8s/manifest"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/spf13/cobra"

	"github.com/cockroachdb/errors"
	batch "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const nameLen = 50

var backOffLimit int32 = 0

type Job struct {
	args *JobParams
}

func Run(
	cli *cobra.Command,
	labels map[string]string,
	cmd []string,
) (*batch.Job, error) {
	manifest, err := NewJob(&JobParams{
		Cli:        cli,
		Labels:     labels,
		Command:    cmd,
		Fragment:   parameters.Fragment,
		NameLength: parameters.NameLen,
	}).MakeSpec()
	if err != nil {
		return &batch.Job{}, errors.Errorf(
			"error in making job manifest %s",
			err,
		)
	}
	namespace, _ := cli.Flags().GetString("namespace")
	job, err := registry.GetKubeClient(registry.KubeClientKey).BatchV1().
		Jobs(namespace).
		Create(context.Background(), manifest, metav1.CreateOptions{})
	if err != nil {
		return job, errors.Errorf("error in deploying job %s", err)
	}

	return job, nil
}

type JobParams struct {
	Cli        *cobra.Command
	Labels     map[string]string
	Command    []string
	NameLength int
	Fragment   string
}

func NewJob(args *JobParams) *Job {
	return &Job{args: args}
}

func (jobk *Job) MakeSpec() (*batch.Job, error) {
	batchJobSpec, err := jobk.batchJobSpec()
	if err != nil {
		return &batch.Job{}, errors.Errorf("error in creating a job %s", err)
	}
	objMeta, err := jobk.objectMeta()
	if err != nil {
		return &batch.Job{}, fmt.Errorf(
			"error in getting a meta object %s",
			err,
		)
	}
	return &batch.Job{
		ObjectMeta: objMeta,
		Spec:       batchJobSpec,
	}, nil
}

func (jobk *Job) objectMeta() (metav1.ObjectMeta, error) {
	namespace, _ := jobk.args.Cli.Flags().GetString("namespace")
	name, _ := jobk.args.Cli.Flags().GetString("job")
	fullName, err := manifest.FullName(name, jobk.args.Fragment, nameLen)
	if err != nil {
		return metav1.ObjectMeta{}, fmt.Errorf(
			"error in manifesting for making the full name %s",
			err,
		)
	}
	return metav1.ObjectMeta{
		Namespace: namespace,
		Name:      fullName,
		Labels:    jobk.args.Labels,
	}, nil
}

func (jobk *Job) batchJobSpec() (batch.JobSpec, error) {
	podTemplSpec, err := jobk.podTemplateSpec()
	if err != nil {
		return batch.JobSpec{}, errors.Errorf(
			"error in getting pod template spec %s",
			err,
		)
	}

	return batch.JobSpec{
		Template:     podTemplSpec,
		BackoffLimit: &backOffLimit,
	}, nil
}

func (jobk *Job) podTemplateSpec() (apiv1.PodTemplateSpec, error) {
	podSpec, err := jobk.podSpec()
	if err != nil {
		return apiv1.PodTemplateSpec{}, errors.Errorf(
			"error in getting pod spec %s",
			err,
		)
	}

	return apiv1.PodTemplateSpec{Spec: podSpec}, nil
}

func (jobk *Job) podSpec() (apiv1.PodSpec, error) {
	contSpec, err := jobk.containersSpec()
	if err != nil {
		return apiv1.PodSpec{}, errors.Errorf(
			"error in getting container spec %s",
			err,
		)
	}

	return apiv1.PodSpec{
		Containers:    contSpec,
		RestartPolicy: apiv1.RestartPolicyNever,
	}, nil
}

func (jobk *Job) containersSpec() ([]apiv1.Container, error) {
	spec := make([]apiv1.Container, 0)
	name, _ := jobk.args.Cli.Flags().GetString("job")
	contName, err := manifest.RandContainerName(
		name,
		jobk.args.Fragment,
		jobk.args.NameLength,
	)
	if err != nil {
		return spec, errors.Errorf(
			"error in generating random container name %s",
			err,
		)
	}

	return append(spec, apiv1.Container{
		Name:    contName,
		Image:   jobk.imageName(),
		Command: jobk.args.Command,
		Env:     jobk.containerEnvSpec(),
	}), nil
}

func (jobk *Job) imageName() string {
	repo, _ := jobk.args.Cli.Flags().GetString("repo")
	tag, _ := jobk.args.Cli.Flags().GetString("tag")
	return fmt.Sprintf("%s:%s", repo, tag)
}

func (jobk *Job) containerEnvSpec() []apiv1.EnvVar {
	namespace, _ := jobk.args.Cli.Flags().GetString("namespace")
	level, _ := jobk.args.Cli.Flags().GetString("log-level")
	return collection.Extend(
		manifest.MinioEnv(namespace),
		manifest.ArangoConfigManifest(namespace),
		manifest.ArangoSecManifest(namespace),
		manifest.LogEnv(level),
	)
}

func MetaLabel() map[string]string {
	return map[string]string{
		"command": "job",
		"runner":  "dictybot",
	}
}
