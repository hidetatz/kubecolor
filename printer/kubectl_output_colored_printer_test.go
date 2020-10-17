package printer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/dty1er/kubecolor/kubectl"
	"github.com/dty1er/kubecolor/testutil"
)

func Test_KubectlOutputColoredPrinter_Print(t *testing.T) {
	tests := []struct {
		name           string
		darkBackground bool
		subcommandInfo *kubectl.SubcommandInfo
		input          string
		expected       string
	}{
		{
			name:           "kubectl top pod",
			darkBackground: true,
			subcommandInfo: &kubectl.SubcommandInfo{
				Subcommand: kubectl.Top,
			},
			input: testutil.NewHereDoc(`
				NAME        CPU(cores)   MEMORY(bytes)
				app-29twd   779m         221Mi
				app-2hhr6   1036m        220Mi
				app-52mbv   881m         137Mi`),
			expected: testutil.NewHereDoc(`
				[37mNAME        CPU(cores)   MEMORY(bytes)[0m
				[36mapp-29twd[0m   [32m779m[0m         [35m221Mi[0m
				[36mapp-2hhr6[0m   [32m1036m[0m        [35m220Mi[0m
				[36mapp-52mbv[0m   [32m881m[0m         [35m137Mi[0m
			`),
		},
		{
			name:           "kubectl top pod --no-headers",
			darkBackground: true,
			subcommandInfo: &kubectl.SubcommandInfo{
				Subcommand: kubectl.Top,
				NoHeader:   true,
			},
			input: testutil.NewHereDoc(`
				app-29twd   779m         221Mi
				app-2hhr6   1036m        220Mi
				app-52mbv   881m         137Mi`),
			expected: testutil.NewHereDoc(`
				[36mapp-29twd[0m   [32m779m[0m         [35m221Mi[0m
				[36mapp-2hhr6[0m   [32m1036m[0m        [35m220Mi[0m
				[36mapp-52mbv[0m   [32m881m[0m         [35m137Mi[0m
			`),
		},
		{
			name:           "kubectl api-resources",
			darkBackground: true,
			subcommandInfo: &kubectl.SubcommandInfo{
				Subcommand: kubectl.APIResources,
			},
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
		{
			name:           "kubectl api-resources --no-headers",
			darkBackground: true,
			subcommandInfo: &kubectl.SubcommandInfo{
				Subcommand: kubectl.APIResources,
				NoHeader:   true,
			},
			input: testutil.NewHereDoc(`
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
		{
			name:           "kubectl get pod",
			darkBackground: true,
			subcommandInfo: &kubectl.SubcommandInfo{
				Subcommand: kubectl.Get,
			},
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
			name:           "kubectl get pod --no-headers",
			darkBackground: true,
			subcommandInfo: &kubectl.SubcommandInfo{
				Subcommand: kubectl.Get,
				NoHeader:   true,
			},
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
			name:           "kubectl get pod -o wide",
			darkBackground: true,
			subcommandInfo: &kubectl.SubcommandInfo{
				Subcommand:   kubectl.Get,
				FormatOption: kubectl.Wide,
			},
			input: testutil.NewHereDoc(`
				NAME                     READY   STATUS    RESTARTS   AGE     IP           NODE       NOMINATED NODE   READINESS GATES
				nginx-6799fc88d8-dnmv5   1/1     Running   0          7d10h   172.18.0.5   minikube   <none>           <none>
				nginx-6799fc88d8-m8pbc   1/1     Running   0          7d10h   172.18.0.4   minikube   <none>           <none>
				nginx-6799fc88d8-qdf9b   1/1     Running   0          7d10h   172.18.0.3   minikube   <none>           <none>`),
			expected: testutil.NewHereDoc(`
				[37mNAME                     READY   STATUS    RESTARTS   AGE     IP           NODE       NOMINATED NODE   READINESS GATES[0m
				[36mnginx-6799fc88d8-dnmv5[0m   [32m1/1[0m     [35mRunning[0m   [37m0[0m          [33m7d10h[0m   [36m172.18.0.5[0m   [32mminikube[0m   [35m<none>[0m           [37m<none>[0m
				[36mnginx-6799fc88d8-m8pbc[0m   [32m1/1[0m     [35mRunning[0m   [37m0[0m          [33m7d10h[0m   [36m172.18.0.4[0m   [32mminikube[0m   [35m<none>[0m           [37m<none>[0m
				[36mnginx-6799fc88d8-qdf9b[0m   [32m1/1[0m     [35mRunning[0m   [37m0[0m          [33m7d10h[0m   [36m172.18.0.3[0m   [32mminikube[0m   [35m<none>[0m           [37m<none>[0m
			`),
		},
		{
			name:           "kubectl get pod -o json",
			darkBackground: true,
			subcommandInfo: &kubectl.SubcommandInfo{
				Subcommand:   kubectl.Get,
				FormatOption: kubectl.Json,
			},
			input: testutil.NewHereDoc(`
				{
				    "apiVersion": "v1",
				    "kind": "Pod",
				    "num": 598,
				    "bool": true,
				    "null": null
				}`),
			expected: testutil.NewHereDoc(`
				{
				    "[37mapiVersion[0m": "[36mv1[0m",
				    "[37mkind[0m": "[36mPod[0m",
				    "[37mnum[0m": [35m598[0m,
				    "[37mbool[0m": [32mtrue[0m,
				    "[37mnull[0m": [33mnull[0m
				}
			`),
		},
		{
			name:           "kubectl get pod -o yaml",
			darkBackground: true,
			subcommandInfo: &kubectl.SubcommandInfo{
				Subcommand:   kubectl.Get,
				FormatOption: kubectl.Yaml,
			},
			input: testutil.NewHereDoc(`
				apiVersion: v1
				kind: "Pod"
				num: 415
				unknown: <unknown>
				none: <none>
				bool: true`),
			expected: testutil.NewHereDoc(`
				[33mapiVersion[0m: [36mv1[0m
				[33mkind[0m: "[36mPod[0m"
				[33mnum[0m: [35m415[0m
				[33munknown[0m: [33m<unknown>[0m
				[33mnone[0m: [33m<none>[0m
				[33mbool[0m: [32mtrue[0m
			`),
		},
		{
			name:           "kubectl describe pod",
			darkBackground: true,
			subcommandInfo: &kubectl.SubcommandInfo{
				Subcommand: kubectl.Describe,
			},
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
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := strings.NewReader(tt.input)
			var w bytes.Buffer
			printer := KubectlOutputColoredPrinter{
				SubcommandInfo: tt.subcommandInfo,
				DarkBackground: tt.darkBackground,
			}
			printer.Print(r, &w)
			testutil.MustEqual(t, tt.expected, w.String())
		})
	}
}
