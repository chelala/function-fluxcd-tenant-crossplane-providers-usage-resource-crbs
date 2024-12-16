package main

import (
	// Standard library imports
	"context"
	"testing"
	"time"

	// Default imports (third-party packages not matching other prefixes)
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/durationpb"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	// Imports with the prefix github.com/crossplane
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	fnv1 "github.com/crossplane/function-sdk-go/proto/v1"
	"github.com/crossplane/function-sdk-go/resource"
)

// Define a TestFunction that embeds Function and overrides fetchProviderRevisions
type TestFunction struct {
	Function
}

func TestRunFunction(t *testing.T) {

	mockProviderRevisions := &unstructured.UnstructuredList{
		Items: []unstructured.Unstructured{
			{
				Object: map[string]interface{}{
					"apiVersion": "pkg.crossplane.io/v1",
					"kind":       "ProviderRevision",
					"metadata": map[string]interface{}{
						"annotations": map[string]interface{}{
							"friendly-name.meta.crossplane.io": "Provider Kubernetes",
							"meta.crossplane.io/description":   "The Crossplane Kubernetes provider enables management of Kubernetes Objects.\n",
							"meta.crossplane.io/license":       "Apache-2.0",
							"meta.crossplane.io/maintainer":    "Crossplane Maintainers \u003cinfo@crossplane.io\u003e",
							"meta.crossplane.io/readme":        "`provider-kubernetes` is a Crossplane Provider that enables deployment and management\nof arbitrary Kubernetes objects on clusters typically provisioned by Crossplane:\n\n- A `Provider` resource type that only points to a credentials `Secret`.\n- An `Object` resource type that is to manage Kubernetes Objects.\n- A managed resource controller that reconciles `Object` typed resources and manages\narbitrary Kubernetes Objects.\n",
							"meta.crossplane.io/source":        "github.com/crossplane-contrib/provider-kubernetes",
						},
						"creationTimestamp": "2024-12-12T18:44:36Z",
						"finalizers": []interface{}{
							"revision.pkg.crossplane.io",
						},
						"generation": 2,
						"labels": map[string]interface{}{
							"pkg.crossplane.io/package": "provider-kubernetes",
						},
						"name": "provider-kubernetes-71953a1e5c15",
						"ownerReferences": []interface{}{
							map[string]interface{}{
								"apiVersion":         "pkg.crossplane.io/v1",
								"blockOwnerDeletion": true,
								"controller":         true,
								"kind":               "Provider",
								"name":               "provider-kubernetes",
								"uid":                "08f4c121-2953-43e6-a45f-088d311169bb",
							},
						},
						"resourceVersion": "3386349",
						"uid":             "5833647c-50d2-4dd5-b79e-23fe1f8143a3",
					},
					"spec": map[string]interface{}{
						"desiredState":                "Active",
						"ignoreCrossplaneConstraints": false,
						"image":                       "xpkg.upbound.io/upbound/provider-kubernetes:v0.16.0",
						"packagePullPolicy":           "IfNotPresent",
						"revision":                    1,
						"runtimeConfigRef": map[string]interface{}{
							"apiVersion": "pkg.crossplane.io/v1beta1",
							"kind":       "DeploymentRuntimeConfig",
							"name":       "provider-kubernetes-config",
						},
						"skipDependencyResolution": false,
						"tlsClientSecretName":      "provider-kubernetes-tls-client",
						"tlsServerSecretName":      "provider-kubernetes-tls-server",
					},
					"status": map[string]interface{}{
						"conditions": []interface{}{
							map[string]interface{}{
								"lastTransitionTime": "2024-12-12T19:03:42Z",
								"reason":             "HealthyPackageRevision",
								"status":             "True",
								"type":               "Healthy",
							},
						},
						"objectRefs": []interface{}{
							map[string]interface{}{
								"apiVersion": "apiextensions.k8s.io/v1",
								"kind":       "CustomResourceDefinition",
								"name":       "providerconfigusages.kubernetes.crossplane.io",
								"uid":        "d34bccec-2cb7-4a7b-9da1-bfb8b22fede9",
							},
							map[string]interface{}{
								"apiVersion": "apiextensions.k8s.io/v1",
								"kind":       "CustomResourceDefinition",
								"name":       "providerconfigs.kubernetes.crossplane.io",
								"uid":        "e5f5664e-e8ad-42a7-8df1-632dfcb292f5",
							},
							map[string]interface{}{
								"apiVersion": "apiextensions.k8s.io/v1",
								"kind":       "CustomResourceDefinition",
								"name":       "observedobjectcollections.kubernetes.crossplane.io",
								"uid":        "3f2155f7-3a1f-48b7-bf0b-eb5f09ea9e38",
							},
							map[string]interface{}{
								"apiVersion": "apiextensions.k8s.io/v1",
								"kind":       "CustomResourceDefinition",
								"name":       "objects.kubernetes.crossplane.io",
								"uid":        "ad2eb470-5b76-4a59-8c00-570edfe65a19",
							},
						},
					},
				},
			},
			{
				Object: map[string]interface{}{
					"apiVersion": "pkg.crossplane.io/v1",
					"kind":       "ProviderRevision",
					"metadata": map[string]interface{}{
						"annotations": map[string]interface{}{
							"auth.upbound.io/group":            "azure.upbound.io",
							"friendly-name.meta.crossplane.io": "Provider Azure (config)",
							"meta.crossplane.io/description":   "Upbound's official Crossplane provider to manage Microsoft Azure\nconfig services in Kubernetes.\n",
							"meta.crossplane.io/maintainer":    "Upbound \u003csupport@upbound.io\u003e",
							"meta.crossplane.io/readme":        "\nProvider Azure is a Crossplane provider for [Microsoft Azure](https://azure.microsoft.com)\ndeveloped and supported by Upbound.\nAvailable resources and their fields can be found in the [Upbound\nMarketplace](https://marketplace.upbound.io/providers/upbound/provider-azure).\nIf you encounter an issue please reach out on support@upbound.io email\naddress. This is a subpackage for the config API group.\n",
							"meta.crossplane.io/source":        "github.com/crossplane-contrib/provider-upjet-azure",
						},
						"creationTimestamp": "2024-12-12T18:44:36Z",
						"finalizers": []interface{}{
							"revision.pkg.crossplane.io",
						},
						"generation": 1,
						"labels": map[string]interface{}{
							"pkg.crossplane.io/package":         "provider-family-azure",
							"pkg.crossplane.io/provider-family": "provider-family-azure",
						},
						"name": "provider-family-azure-7e0a66cff496",
						"ownerReferences": []interface{}{
							map[string]interface{}{
								"apiVersion":         "pkg.crossplane.io/v1",
								"blockOwnerDeletion": true,
								"controller":         true,
								"kind":               "Provider",
								"name":               "provider-family-azure",
								"uid":                "6cf4624f-cf56-4d05-88ad-a8338f5faea6",
							},
						},
						"resourceVersion": "2993347",
						"uid":             "296f7d6c-39a6-4fce-a8a0-b195d5726b13",
					},
					"spec": map[string]interface{}{
						"desiredState":                "Active",
						"ignoreCrossplaneConstraints": false,
						"image":                       "xpkg.upbound.io/upbound/provider-family-azure:v1.10.0",
						"packagePullPolicy":           "IfNotPresent",
						"revision":                    1,
						"runtimeConfigRef": map[string]interface{}{
							"apiVersion": "pkg.crossplane.io/v1beta1",
							"kind":       "DeploymentRuntimeConfig",
							"name":       "provider-azure-family-config",
						},
						"skipDependencyResolution": false,
						"tlsClientSecretName":      "provider-family-azure-tls-client",
						"tlsServerSecretName":      "provider-family-azure-tls-server",
					},
					"status": map[string]interface{}{
						"conditions": []interface{}{
							map[string]interface{}{
								"lastTransitionTime": "2024-12-12T18:44:50Z",
								"reason":             "HealthyPackageRevision",
								"status":             "True",
								"type":               "Healthy",
							},
						},
						"objectRefs": []interface{}{
							map[string]interface{}{
								"apiVersion": "apiextensions.k8s.io/v1",
								"kind":       "CustomResourceDefinition",
								"name":       "subscriptions.azure.upbound.io",
								"uid":        "8b1ff63c-3822-4cd1-b9fb-54aa3c280cc3",
							},
							map[string]interface{}{
								"apiVersion": "apiextensions.k8s.io/v1",
								"kind":       "CustomResourceDefinition",
								"name":       "storeconfigs.azure.upbound.io",
								"uid":        "098f0df7-9b06-4941-889d-d93368fac8fe",
							},
							map[string]interface{}{
								"apiVersion": "apiextensions.k8s.io/v1",
								"kind":       "CustomResourceDefinition",
								"name":       "resourceproviderregistrations.azure.upbound.io",
								"uid":        "2547646d-0473-4673-9d72-3c068c5d83b0",
							},
							map[string]interface{}{
								"apiVersion": "apiextensions.k8s.io/v1",
								"kind":       "CustomResourceDefinition",
								"name":       "resourcegroups.azure.upbound.io",
								"uid":        "e9fc1a11-2734-4745-85be-08c959118a2e",
							},
							map[string]interface{}{
								"apiVersion": "apiextensions.k8s.io/v1",
								"kind":       "CustomResourceDefinition",
								"name":       "providerconfigusages.azure.upbound.io",
								"uid":        "17840102-8815-4be2-8e4b-32ca6c585003",
							},
							map[string]interface{}{
								"apiVersion": "apiextensions.k8s.io/v1",
								"kind":       "CustomResourceDefinition",
								"name":       "providerconfigs.azure.upbound.io",
								"uid":        "384490b0-efce-416a-81c7-711f22b8ddd8",
							},
						},
					},
				},
			},
		},
	}

	expectedDesiredComposed := map[string]*fnv1.Resource{
		"demo000-provider-kubernetes-edit": {
			Resource: resource.MustStructJSON(`{
				"apiVersion": "kubernetes.crossplane.io/v1alpha2",
				"kind": "Object",
				"metadata": {
					"annotations": {
						"crossplane.io/external-name": "demo000-provider-kubernetes-edit"
					}
				},
				"spec": {
					"forProvider": {
						"manifest": {
							"apiVersion": "rbac.authorization.k8s.io/v1",
							"kind": "ClusterRoleBinding",
							"metadata": {
								"labels": {
									"kustomize.toolkit.fluxcd.io/name": "tenants",
									"kustomize.toolkit.fluxcd.io/namespace": "flux-system"
								},
								"name": "demo000-provider-kubernetes-edit"
							},
							"roleRef": {
								"apiGroup": "rbac.authorization.k8s.io",
								"kind": "ClusterRole",
								"name": "crossplane:provider:provider-kubernetes-71953a1e5c15:aggregate-to-edit"
							},
							"subjects": [
								{
									"kind": "ServiceAccount",
									"name": "demo000",
									"namespace": "demo000"
								}
							]
						}
					},
					"watch": false
				},
				"status": {
					"observedGeneration": 0
				}
			}`),
		},
		"demo000-provider-family-azure-edit": {
			Resource: resource.MustStructJSON(`{
				"apiVersion": "kubernetes.crossplane.io/v1alpha2",
				"kind": "Object",
				"metadata": {
					"annotations": {
						"crossplane.io/external-name": "demo000-provider-family-azure-edit"
					}
				},
				"spec": {
					"forProvider": {
						"manifest": {
							"apiVersion": "rbac.authorization.k8s.io/v1",
							"kind": "ClusterRoleBinding",
							"metadata": {
								"labels": {
									"kustomize.toolkit.fluxcd.io/name": "tenants",
									"kustomize.toolkit.fluxcd.io/namespace": "flux-system"
								},
								"name": "demo000-provider-family-azure-edit"
							},
							"roleRef": {
								"apiGroup": "rbac.authorization.k8s.io",
								"kind": "ClusterRole",
								"name": "crossplane:provider:provider-family-azure-7e0a66cff496:aggregate-to-edit"
							},
							"subjects": [
								{
									"kind": "ServiceAccount",
									"name": "demo000",
									"namespace": "demo000"
								}
							]
						}
					},
					"watch": false
				},
				"status": {
					"observedGeneration": 0
				}
			}`),
		},
	}

	type args struct {
		ctx context.Context
		req *fnv1.RunFunctionRequest
	}
	type want struct {
		rsp *fnv1.RunFunctionResponse
		err error
	}

	cases := map[string]struct {
		reason string
		args   args
		want   want
	}{
		"TwoCRBsforTwoProviderRevisions": {
			reason: "The Function should return two new ClusterRoleBindings for the observed XFluxcdTenant. Given there are two ProviderRevisions in the mock response.",
			args: args{
				req: &fnv1.RunFunctionRequest{
					Observed: &fnv1.State{
						Composite: &fnv1.Resource{
							Resource: resource.MustStructJSON(`{
							    "apiVersion": "gitops.idp.someorg.com/v1alpha1",
							    "kind": "XFluxcdTenant",
							    "spec": {
							        "gitAuthProvider": "azure",
							        "gitBranch": "main",
							        "gitPath": "/demo000",
							        "gitUrl": "https://dev.azure.com/Someorg/prj-idp2/_git/repo-idp2",
							        "tenantName": "demo000"
							    }
							}`),
						},
					},
				},
			},
			want: want{
				rsp: &fnv1.RunFunctionResponse{
					Meta: &fnv1.ResponseMeta{Ttl: durationpb.New(60 * time.Second)},
					Conditions: []*fnv1.Condition{
						{
							Type:   "FunctionSuccess",
							Status: fnv1.Status_STATUS_CONDITION_TRUE,
							Reason: "Success",
							Target: fnv1.Target_TARGET_COMPOSITE_AND_CLAIM.Enum(),
						},
					},
					Desired: &fnv1.State{
						Resources: expectedDesiredComposed,
					},
				},
				err: nil,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			testFunc := &TestFunction{
				Function: Function{
					log: logging.NewNopLogger(),
					fetchProviderRevisionsFunc: func(_ context.Context, _ logging.Logger) (*unstructured.UnstructuredList, error) {
						return mockProviderRevisions, nil
					},
				},
			}
			rsp, err := testFunc.RunFunction(tc.args.ctx, tc.args.req)

			if diff := cmp.Diff(tc.want.err, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("%s\nf.RunFunction(...): -want err, +got err:\n%s", tc.reason, diff)
			}

			// Compare the response meta and conditions
			if diff := cmp.Diff(tc.want.rsp.GetMeta(), rsp.GetMeta(), protocmp.Transform()); diff != "" {
				t.Errorf("%s\nf.RunFunction(...): -want meta, +got meta:\n%s", tc.reason, diff)
			}
			if diff := cmp.Diff(tc.want.rsp.GetConditions(), rsp.GetConditions(), protocmp.Transform()); diff != "" {
				t.Errorf("%s\nf.RunFunction(...): -want conditions, +got conditions:\n%s", tc.reason, diff)
			}

			// Compare the desired composed resources
			if diff := cmp.Diff(tc.want.rsp.GetDesired().GetResources(), rsp.GetDesired().GetResources(), protocmp.Transform()); diff != "" {
				t.Errorf("%s\nf.RunFunction(...): -want desired composed, +got desired composed:\n%s", tc.reason, diff)
			}
		})
	}
}
