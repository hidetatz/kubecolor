package printer

import (
	"bytes"
	"strconv"
	"strings"
	"testing"

	"github.com/dty1er/kubecolor/color"
	"github.com/dty1er/kubecolor/testutil"
)

func Test_TablePrinter_Print(t *testing.T) {
	tests := []struct {
		name           string
		colorDeciderFn func(index int, column string) (color.Color, bool)
		withHeader     bool
		darkBackground bool
		input          string
		expected       string
	}{
		{
			name:           "header is not colored - dark",
			colorDeciderFn: nil,
			withHeader:     true,
			darkBackground: true,
			input: testutil.NewHereDoc(`
				NAME          READY   STATUS    RESTARTS   AGE
				nginx-dnmv5   1/1     Running   0          6d6h
				nginx-m8pbc   1/1     Running   0          6d6h
				nginx-qdf9b   1/1     Running   0          6d6h`),
			expected: testutil.NewHereDoc(`
				[37mNAME          READY   STATUS    RESTARTS   AGE[0m
				[36mnginx-dnmv5[0m   [32m1/1[0m     [35mRunning[0m   [37m0[0m          [33m6d6h[0m
				[36mnginx-m8pbc[0m   [32m1/1[0m     [35mRunning[0m   [37m0[0m          [33m6d6h[0m
				[36mnginx-qdf9b[0m   [32m1/1[0m     [35mRunning[0m   [37m0[0m          [33m6d6h[0m
			`),
		},
		{
			name:           "withheader=false, 1st line is not colored in header color but colored as a content of table",
			colorDeciderFn: nil,
			withHeader:     false,
			darkBackground: true,
			input: testutil.NewHereDoc(`
				nginx-dnmv5   1/1     Running   0          6d6h
				nginx-m8pbc   1/1     Running   0          6d6h
				nginx-qdf9b   1/1     Running   0          6d6h`),
			expected: testutil.NewHereDoc(`
				[36mnginx-dnmv5[0m   [32m1/1[0m     [35mRunning[0m   [37m0[0m          [33m6d6h[0m
				[36mnginx-m8pbc[0m   [32m1/1[0m     [35mRunning[0m   [37m0[0m          [33m6d6h[0m
				[36mnginx-qdf9b[0m   [32m1/1[0m     [35mRunning[0m   [37m0[0m          [33m6d6h[0m
				`),
		},
		{
			name:           "when darkBackground=false, color preset for light is used",
			colorDeciderFn: nil,
			withHeader:     true,
			darkBackground: false,
			input: testutil.NewHereDoc(`
				NAME          READY   STATUS    RESTARTS   AGE
				nginx-dnmv5   1/1     Running   0          6d6h
				nginx-m8pbc   1/1     Running   0          6d6h
				nginx-qdf9b   1/1     Running   0          6d6h`),
			expected: testutil.NewHereDoc(`
				[30mNAME          READY   STATUS    RESTARTS   AGE[0m
				[36mnginx-dnmv5[0m   [32m1/1[0m     [35mRunning[0m   [30m0[0m          [33m6d6h[0m
				[36mnginx-m8pbc[0m   [32m1/1[0m     [35mRunning[0m   [30m0[0m          [33m6d6h[0m
				[36mnginx-qdf9b[0m   [32m1/1[0m     [35mRunning[0m   [30m0[0m          [33m6d6h[0m
			`),
		},
		{
			name: "colorDeciderFn works",
			colorDeciderFn: func(_ int, column string) (color.Color, bool) {
				if column == "CrashLoopBackOff" {
					return color.Red, true
				}

				// When Readiness is "n/m" then yellow
				if strings.Count(column, "/") == 1 {
					if arr := strings.Split(column, "/"); arr[0] != arr[1] {
						_, e1 := strconv.Atoi(arr[0])
						_, e2 := strconv.Atoi(arr[1])
						if e1 == nil && e2 == nil { // check both is number
							return color.Yellow, true
						}
					}

				}

				return 0, false
			},
			withHeader:     true,
			darkBackground: true,
			// "CrashLoopBackOff" will be red, "0/1" will be yellow
			input: testutil.NewHereDoc(`
				NAME          READY   STATUS             RESTARTS   AGE
				nginx-dnmv5   1/1     CrashLoopBackOff   0          6d6h
				nginx-m8pbc   1/1     Running            0          6d6h
				nginx-qdf9b   0/1     Running            0          6d6h`),
			expected: testutil.NewHereDoc(`
				[37mNAME          READY   STATUS             RESTARTS   AGE[0m
				[36mnginx-dnmv5[0m   [32m1/1[0m     [31mCrashLoopBackOff[0m   [37m0[0m          [33m6d6h[0m
				[36mnginx-m8pbc[0m   [32m1/1[0m     [35mRunning[0m            [37m0[0m          [33m6d6h[0m
				[36mnginx-qdf9b[0m   [33m0/1[0m     [35mRunning[0m            [37m0[0m          [33m6d6h[0m
			`),
		},
		{
			name:           "a table whose some parts are missing can be handled",
			colorDeciderFn: nil,
			withHeader:     true,
			darkBackground: true,
			input: testutil.NewHereDoc(`
				NAME                              SHORTNAMES   APIGROUP                       NAMESPACED   KIND
				bindings                                                                      true         Binding
				componentstatuses                 cs                                          false        ComponentStatus
				pods                              po                                          true         Pod
				podtemplates                                                                  true         PodTemplate
				replicationcontrollers            rc                                          true         ReplicationController
				resourcequotas                    quota                                       true         ResourceQuota
				secrets                                                                       true         Secret
				serviceaccounts                   sa                                          true         ServiceAccount
				services                          svc                                         true         Service
				mutatingwebhookconfigurations                  admissionregistration.k8s.io   false        MutatingWebhookConfiguration
				customresourcedefinitions         crd,crds     apiextensions.k8s.io           false        CustomResourceDefinition
				controllerrevisions                            apps                           true         ControllerRevision
				daemonsets                        ds           apps                           true         DaemonSet
				statefulsets                      sts          apps                           true         StatefulSet
				tokenreviews                                   authentication.k8s.io          false        TokenReview
			`),
			expected: testutil.NewHereDoc(`
				[37mNAME                              SHORTNAMES   APIGROUP                       NAMESPACED   KIND[0m
				[36mbindings[0m                                                                      [32mtrue[0m         [35mBinding[0m
				[36mcomponentstatuses[0m                 [37mcs[0m                                          [32mfalse[0m        [35mComponentStatus[0m
				[36mpods[0m                              [37mpo[0m                                          [32mtrue[0m         [35mPod[0m
				[36mpodtemplates[0m                                                                  [32mtrue[0m         [35mPodTemplate[0m
				[36mreplicationcontrollers[0m            [37mrc[0m                                          [32mtrue[0m         [35mReplicationController[0m
				[36mresourcequotas[0m                    [37mquota[0m                                       [32mtrue[0m         [35mResourceQuota[0m
				[36msecrets[0m                                                                       [32mtrue[0m         [35mSecret[0m
				[36mserviceaccounts[0m                   [37msa[0m                                          [32mtrue[0m         [35mServiceAccount[0m
				[36mservices[0m                          [37msvc[0m                                         [32mtrue[0m         [35mService[0m
				[36mmutatingwebhookconfigurations[0m                  [33madmissionregistration.k8s.io[0m   [32mfalse[0m        [35mMutatingWebhookConfiguration[0m
				[36mcustomresourcedefinitions[0m         [37mcrd,crds[0m     [33mapiextensions.k8s.io[0m           [32mfalse[0m        [35mCustomResourceDefinition[0m
				[36mcontrollerrevisions[0m                            [33mapps[0m                           [32mtrue[0m         [35mControllerRevision[0m
				[36mdaemonsets[0m                        [37mds[0m           [33mapps[0m                           [32mtrue[0m         [35mDaemonSet[0m
				[36mstatefulsets[0m                      [37msts[0m          [33mapps[0m                           [32mtrue[0m         [35mStatefulSet[0m
				[36mtokenreviews[0m                                   [33mauthentication.k8s.io[0m          [32mfalse[0m        [35mTokenReview[0m
			`),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := strings.NewReader(tt.input)
			var w bytes.Buffer
			printer := NewTablePrinter(tt.withHeader, tt.darkBackground, tt.colorDeciderFn)
			printer.Print(r, &w)
			testutil.MustEqual(t, tt.expected, w.String())
		})
	}
}
