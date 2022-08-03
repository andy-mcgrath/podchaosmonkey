# PodChaosMonkey

PodChaosMonkey is designed to run and poll for Pods within a single Kubernetes Namespace and will randomly delete a Pod with a name that matches supplied RegEx filter string.

---

## Configuration

Configuration is provided by a set of Environment variables which can be provided as standard with a Pod or from a `.env` file.

The exception to this is the `ENVIRONMENT` environment variable which is used to set the `.env` file as either `development.env` or `production.env`.

The following are the configuration options available

| Configuration Options | Default Values | Description |
| - | - | - |
| ENVIRONMENT | "DEV" | For running under local development defaults to `DEV` and will load the file `dev.env`, setting to anything else will load configuration options from the file `production.dev` |
| KILL_TIME_DELAY | 120 | Time delay in seconds between PodChaosMonkey deleting a random pod |
| LOG_LEVEL | "info" | *Currently not in use* |
| NAMESPACE | "chaos" | Namespace where PodChaosMonkey will delete pods |
| POD_FILTER | ".*" | RegEx string used to filter pod names during selection, only pod names matching the filter will be selected. Default all pods within the set Namespace will be selected |
| DEFAULT_TIMEOUT | 5 | Timeout in seconds used during calls to the Kubernetes API |

NOTE: The file `production.env` is copied into the container image but these can be overwritten by setting Environment Variables

---

## Development Environment

[Tilt.dev](https://tilt.dev/) has been used for local development.

A local running Kubernetes is required and the Tiltfile has been limited to run only on the following contexts:

- rancher-desktop
- minikube
- dockerdesktop
- kind

To using install Tilt, ensure a local Kubernetes is running and your local context is set corectly then run `tilt up` in PodChaosMonkey's root directory.

Alternatively the following commands will build and deploy to a Kubernetes cluster:

Docker build:
```
docker build -t podchaosmonkey -f ./deployment/Dockerfile .
```

Kubernetes Deployment:
```
kubectl apply -f ./deployment/kubernetes.yaml
```

Options Test Nginx deployment:
```
kubectl apply -f ./deployment/test-deployment.yaml
```

---

## Notes

- Language: Go v1.18
  - `pkg/config` and `pkg/kuberclientset` are from another application with minor modifictations
- Testing on Rancher-Desktop, Kubernetes v1.23.8
- Docker service image is using `alpine` to allow Tilt to work in a development environment.
- Currently limited to polling a single Namespace, future feature to expact to support multiple or whole cluster.
- ServiceAccount for PodChaosMonkey assigned clusterwide permissions, future feature using `RoleBinding` per namespace over a `ClusterRoleBinding`
- Pod filter is currently based on the pod name only, future feature to update to allow for wider selection options e.g. labels, annotation, based on selected deployments etc.

---
