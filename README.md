# Function Fluxcd Tenant Crossplane Providers Usage Resource CRBs
This composition Function creates the necessary Cluster Role Bindings so Fluxcd Tenant Service Account can create the ProviderUsage object for every managed resource, when composition is applied as the Fluxcd Tenan Service Account defined.

## Details
If using fluxcd and using the multi tenant recommendation (https://fluxcd.io/flux/cmd/flux_create_tenant/ and https://github.com/fluxcd/flux2-multi-tenancy/tree/dev-team) upon creating XR or XRC an object will be created in for every Crossplane provider every time is used (providerusage).
Given a Crossplane composition is created for a Fluxcd Tenant is created, the composition function will add the CRBs for that FluxCD Tenant, which is created in a separated namespace with SAs limited to that namespace and no permission to create ProviderUsage at the cluster level.

## Usage
Instead of fluxcd cli like:
```bash
# Create a tenant with access to a namespace 
flux create tenant dev-team \
  --with-namespace=frontend \
  --label=environment=dev

# Generate tenant namespaces and role bindings in YAML format
flux create tenant dev-team \
  --with-namespace=frontend \
  --with-namespace=backend \
  --export > dev-team.yaml
```
Here is snippet for an equivalent Crossplane composition and the function:

```yaml
apiVersion: apiextensions.crossplane.io/v1
kind: Composition
metadata:
  name: xfluxcdtenants.gitops.idp.chelala.one
spec:
  compositeTypeRef:
    apiVersion: gitops.idp.chelala.one/v1alpha1
    kind: XFluxcdTenant
  mode: Pipeline
  pipeline:
    - step: patch-and-transform
      functionRef:
        name: function-patch-and-transform
      input:
        apiVersion: pt.fn.crossplane.io/v1beta1
        kind: Resources
        resources:
          - name: ftnamespace
            base:
         ................
          - name: ftsa
            base:
         ........
          - name: ftrb
            base:
          ............
          - name: ftgitrepo
            base:
           ...............
          - name: ftkusto
            base:
            ...............
    - step: function-fluxcd-tenant-crossplane-providers-usage-resource-crbs
      functionRef:
        name: function-fluxcd-tenant-crossplane-providers-usage-resource-crbs
    - step: automatically-detect-ready-composed-resources
      functionRef:
        name: function-auto-ready
```

# General Information for function-template-go
[![CI](https://github.com/chelala/function-fluxcd-tenant-crossplane-providers-usage-resource-crbs/actions/workflows/ci.yml/badge.svg)](https://github.com/chelala/function-fluxcd-tenant-crossplane-providers-usage-resource-crbs/actions/workflows/ci.yml)

A template for writing a [composition function][functions] in [Go][go].

To learn how to use this template:

* [Follow the guide to writing a composition function in Go][function guide]
* [Learn about how composition functions work][functions]
* [Read the function-sdk-go package documentation][package docs]

If you just want to jump in and get started:

1. Replace `function-template-go` with your function in `go.mod`,
   `package/crossplane.yaml`, and any Go imports. (You can also do this
   automatically by running the `./init.sh <function-name>` script.)
1. Update `input/v1beta1/` to reflect your desired input (and run `go generate ./...`)
1. Add your logic to `RunFunction` in `fn.go`
1. Add tests for your logic in `fn_test.go`
1. Update this file, `README.md`, to be about your function!

This template uses [Go][go], [Docker][docker], and the [Crossplane CLI][cli] to
build functions.

```shell
# Run code generation - see input/generate.go
$ go generate ./...

# Run tests - see fn_test.go
$ go test ./...

# Build the function's runtime image - see Dockerfile
$ docker build . --tag=runtime

# Build a function package - see package/crossplane.yaml
$ crossplane xpkg build -f package --embed-runtime-image=runtime
```

[functions]: https://docs.crossplane.io/latest/concepts/composition-functions
[go]: https://go.dev
[function guide]: https://docs.crossplane.io/knowledge-base/guides/write-a-composition-function-in-go
[package docs]: https://pkg.go.dev/github.com/crossplane/function-sdk-go
[docker]: https://www.docker.com
[cli]: https://docs.crossplane.io/latest/cli

Tutorial: https://docs.crossplane.io/v1.18/guides/write-a-composition-function-in-go/

macos Docker Desktop fix:
```bash
sudo ln -s /Users/chelala/.docker/run/docker.sock /var/run/docker.sock
```
