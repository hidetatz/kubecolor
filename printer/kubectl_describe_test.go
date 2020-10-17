package printer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/dty1er/kubecolor/testutil"
)

func Test_DescribePrinter_Print(t *testing.T) {
	tests := []struct {
		name           string
		darkBackground bool
		input          string
		expected       string
	}{
		{
			name:           "values can be colored by its type",
			darkBackground: true,
			input: testutil.NewHereDoc(`
				Name:         nginx-lpv5x
				Namespace:    default
				Priority:     0
				Node:         minikube/172.17.0.3
				Ready:        true
				Start Time:   Sat, 10 Oct 2020 14:07:17 +0900
				Labels:       app=nginx
				Annotations:  <none>`),
			expected: testutil.NewHereDoc(`
				[33mName[0m:         [36mnginx-lpv5x[0m
				[33mNamespace[0m:    [36mdefault[0m
				[33mPriority[0m:     [35m0[0m
				[33mNode[0m:         [36mminikube/172.17.0.3[0m
				[33mReady[0m:        [32mtrue[0m
				[33mStart Time[0m:   [36mSat, 10 Oct 2020 14:07:17 +0900[0m
				[33mLabels[0m:       [36mapp=nginx[0m
				[33mAnnotations[0m:  [33m<none>[0m
			`),
		},
		{
			name:           "key color changes based on its indentation",
			darkBackground: true,
			input: testutil.NewHereDoc(`
				IP:           172.18.0.7
				IPs:
				  IP:           172.18.0.7
				Controlled By:  ReplicaSet/nginx
				Containers:
				  nginx:
				    Container ID:   docker://2885230a30908c8a6bda5a5366619c730b25b994eea61c931bba08ef4a8c8593
				      Started:      Sat, 10 Oct 2020 14:07:44 +0900`),
			expected: testutil.NewHereDoc(`
				[33mIP[0m:           [36m172.18.0.7[0m
				[33mIPs[0m:
				  [37mIP[0m:           [36m172.18.0.7[0m
				[33mControlled By[0m:  [36mReplicaSet/nginx[0m
				[33mContainers[0m:
				  [37mnginx[0m:
				    [33mContainer ID[0m:   [36mdocker://2885230a30908c8a6bda5a5366619c730b25b994eea61c931bba08ef4a8c8593[0m
				      [37mStarted[0m:      [36mSat, 10 Oct 2020 14:07:44 +0900[0m
			`),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := strings.NewReader(tt.input)
			var w bytes.Buffer
			printer := DescribePrinter{DarkBackground: tt.darkBackground}
			printer.Print(r, &w)
			testutil.MustEqual(t, tt.expected, w.String())
		})
	}
}
