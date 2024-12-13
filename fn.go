package main

import (
	"context"
	"fmt"

	"github.com/crossplane-contrib/provider-kubernetes/apis/object/v1alpha2"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/function-sdk-go/errors"
	fnv1 "github.com/crossplane/function-sdk-go/proto/v1"
	"github.com/crossplane/function-sdk-go/request"
	"github.com/crossplane/function-sdk-go/resource"
	"github.com/crossplane/function-sdk-go/resource/composed"
	"github.com/crossplane/function-sdk-go/response"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

// Function returns whatever response you ask it to.
type Function struct {
	fnv1.UnimplementedFunctionRunnerServiceServer

	log logging.Logger
}

// RunFunction runs the Function.
func (f *Function) RunFunction(_ context.Context, req *fnv1.RunFunctionRequest) (*fnv1.RunFunctionResponse, error) {
	f.log.Info("Running function", "tag", req.GetMeta().GetTag())

	rsp := response.To(req, response.DefaultTTL)

	// 1. Create a Kubernetes client
	config, err := rest.InClusterConfig()
	if err != nil {
		f.log.Info("Failed to get in-cluster config", "error", err)
	}
	// Create a dynamic client
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		f.log.Info("Failed to create dynamic client", "error", err)
	}

	// Define GVR for ProviderRevisions
	gvr := schema.GroupVersionResource{
		Group:    "pkg.crossplane.io",
		Version:  "v1",
		Resource: "providerrevisions",
	}

	// 2. List ProviderRevisions
	providerRevisions, err := dynamicClient.Resource(gvr).Namespace("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		f.log.Info("Failed to list ProviderRevisions", "error", err)
	}

	// Get all desired composed resources from the request. The function will
	// update this map of resources, then save it. This get, update, set pattern
	// ensures the function keeps any resources added by other functions.
	desired, err := request.GetDesiredComposedResources(req)
	if err != nil {
		response.Fatal(rsp, errors.Wrapf(err, "cannot get desired resources from %T", req))
		return rsp, nil
	}

	xr, err := request.GetObservedCompositeResource(req)
	if err != nil {
		response.Fatal(rsp, errors.Wrapf(err, "cannot get the composite resource from %T", req))
		return rsp, nil
	}
	tenanName, err := xr.Resource.GetString("spec.tenantName")
	if err != nil {
		response.Fatal(rsp, errors.Wrapf(err, "cannot get the XR tenant name from %T", req))
		return rsp, nil
	}

	// Add object/v1alpha2 types (including object) to the composed resource scheme.
	// composed. From uses this to automatically set apiVersion and kind.
	_ = v1alpha2.SchemeBuilder.AddToScheme(composed.Scheme)

	// 3. Process the results
	for _, pr := range providerRevisions.Items {
		f.log.Info("XR tenant name", "tenantName", tenanName)
		f.log.Info("Label pkg.crossplane.io/package", "providerName", pr.GetLabels()["pkg.crossplane.io/package"])
		f.log.Info("ProviderRevision Name", "providerRevisionName", pr.GetName())

		manifestFmt := []byte(`{
		    "apiVersion": "rbac.authorization.k8s.io/v1",
		    "kind": "ClusterRoleBinding",
		    "metadata": {
		        "labels": {
		            "kustomize.toolkit.fluxcd.io/name": "tenants",
		            "kustomize.toolkit.fluxcd.io/namespace": "flux-system"
		        },
		        "name": "%s-%s-edit"
		    },
		    "roleRef": {
		        "apiGroup": "rbac.authorization.k8s.io",
		        "kind": "ClusterRole",
		        "name": "crossplane:provider:%s:aggregate-to-edit"
		    },
		    "subjects": [
		        {
		            "kind": "ServiceAccount",
		            "name": "%s",
		            "namespace": "%s"
		        }
		    ]
		}
		`)
		ocrb := &v1alpha2.Object{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{
					"crossplane.io/external-name": fmt.Sprintf("%s-%s-edit", tenanName, pr.GetLabels()["pkg.crossplane.io/package"]),
				},
			},
			Spec: v1alpha2.ObjectSpec{
				ForProvider: v1alpha2.ObjectParameters{
					Manifest: runtime.RawExtension{
						Raw: []byte(fmt.Sprintf(
							string(manifestFmt),
							tenanName,
							pr.GetLabels()["pkg.crossplane.io/package"],
							pr.GetName(),
							tenanName,
							tenanName,
						)),
					},
				},
			},
		}

		// Convert the object to the unstructured resource data format the SDK
		// uses to store desired composed resources.
		unsocrb, err := composed.From(ocrb)
		if err != nil {
			response.Fatal(rsp, errors.Wrapf(err, "cannot convert %T to %T", unsocrb, &composed.Unstructured{}))
			return rsp, nil
		}

		// Add the bucket to the map of desired composed resources. It's
		// important that the function adds the same bucket every time it's
		// called. It's also important that the bucket is added with the same
		// resource.Name every time it's called. The function prefixes the name
		// with "xbuckets-" to avoid collisions with any other composed
		// resources that might be in the desired resources map.
		desired[resource.Name(fmt.Sprintf("%s-%s-edit", tenanName, pr.GetLabels()["pkg.crossplane.io/package"]))] = &resource.DesiredComposed{Resource: unsocrb}
	}

	// Finally, save the updated desired composed resources to the response.
	if err := response.SetDesiredComposedResources(rsp, desired); err != nil {
		response.Fatal(rsp, errors.Wrapf(err, "cannot set desired composed resources in %T", rsp))
		return rsp, nil
	}

	// Log what the function did. This will only appear in the function's pod
	// logs. A function can use response.Normal and response.Warning to emit
	// Kubernetes events associated with the XR it's operating on.
	f.log.Info("Added necessary cluster role bindings so fluxcd tenant sa can edit crossplane providers usage resources", "tenantName", tenanName)

	// You can set a custom status condition on the claim. This allows you to
	// communicate with the user. See the link below for status condition
	// guidance.
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties
	response.ConditionTrue(rsp, "FunctionSuccess", "Success").
		TargetCompositeAndClaim()

	return rsp, nil
}
