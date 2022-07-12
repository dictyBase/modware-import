package manifest

import (
	"fmt"
)

type ImageSpec struct {
	repo, tag  string
	PullPolicy string
}

func NewImageSpec(repo, tag, policy string) *ImageSpec {
	return &ImageSpec{
		repo:       repo,
		tag:        tag,
		PullPolicy: policy,
	}
}

func (s *ImageSpec) ImageManifest() string {
	return fmt.Sprintf("%s:%s", s.repo, s.tag)
}
