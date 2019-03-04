# gen-api-docs

Generates API Docs based on Rancher Schema

* `data` - Static data, generic descriptions, base objects...
* `openapi` - Types for openapi v3.
* `build` - Rendered swagger/build doc output.

## Testing/Running Locally

Running the default `make` target will create a RKE k8s cluster with DinD, install Rancher and run the main.go against that instance.

```plain
make
```

## Running the Container

Run the resulting image.  The swagger-ui image listens on 8080/tcp

```plain
docker run -d -p 8080:8080 --name swagger rancher/gen-api-docs:dev
```

Browse to http://localhost:8080