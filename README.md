# gen-api-docs

Generates API Docs based on Rancher Schema

* `data` - Static data, generic descriptions, base objects...
* `openapi` - Types for openapi v3.
* `build` - Rendered swagger/build doc output.
* `templates` - Go templates for Markdown output (not used)

## Testing/Running Locally

Running the default `make` target will create a RKE k8s cluster with DinD, install Rancher and run the main.go against that instance.

```
make
```

Running against an existing Rancher instance. This will require a cluster to be defined.

```
export RANCHER_TOKEN="<bearer token>"
export RANCHER_URL="https://rancher.mydomain.com/v3"
go run main.go
```