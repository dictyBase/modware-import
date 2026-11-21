package wait

import (
	"testing"

	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
)

func TestFindStuckReasonNullableWaiting(t *testing.T) {
	cases := []struct {
		name    string
		waiting *corev1.ContainerStateWaiting
		want    O.Option[string]
	}{
		{name: "nil", want: O.None[string]()},
		{
			name:    "healthy",
			waiting: &corev1.ContainerStateWaiting{Reason: "ContainerCreating"},
			want:    O.None[string](),
		},
		{
			name:    "stuck",
			waiting: &corev1.ContainerStateWaiting{Reason: "ImagePullBackOff"},
			want:    O.Some("ImagePullBackOff"),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pods := &corev1.PodList{Items: []corev1.Pod{{Status: corev1.PodStatus{
				ContainerStatuses: []corev1.ContainerStatus{
					{State: corev1.ContainerState{Waiting: tc.waiting}},
				},
			}}}}
			require.Equal(t, tc.want, FindStuckReason(pods))
		})
	}
}
