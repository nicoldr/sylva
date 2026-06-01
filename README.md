# sylvactl

A CLI tool to retrieve FluxCD resources (`HelmRelease`, `Kustomization`) from a Kubernetes cluster.

## How it works

`sylvactl` uses the [controller-runtime](https://github.com/kubernetes-sigs/controller-runtime) client to connect to a Kubernetes cluster via kubeconfig and query FluxCD Custom Resources. FluxCD extends Kubernetes with its own CRDs — `HelmRelease` (managed by helm-controller) and `Kustomization` (managed by kustomize-controller). `sylvactl` registers these CRD types into a runtime scheme and lists them across namespaces, presenting the results in a human-readable table or structured JSON/YAML format. The CLI is built with [Cobra](https://github.com/spf13/cobra), following the same command structure as `kubectl`.

## Project Structure

```
sylva/
├── cmd/
│   └── sylvactl/
│       └── main.go              # binary entrypoint
├── internal/
│   └── cmd/
│       ├── root.go              # root cobra command and global flags
│       ├── get.go               # "get" subcommand
│       ├── helmrelease.go       # "get helmreleases" implementation
│       ├── kustomization.go     # "get kustomizations" implementation
│       ├── client.go            # Kubernetes client setup
│       └── helpers.go           # shared utilities (formatting, conditions)
├── cluster/
│   ├── apps/
│   │   └── podinfo/             # sample HelmRelease for testing
│   └── flux/                    # GitRepository and Kustomization for GitOps wiring
├── go.mod
└── go.sum
```

## Prerequisites

| Tool | Version | Install |
|------|---------|---------|
| Docker | any | https://docs.docker.com/get-docker/ |
| kind | v0.23.0+ | https://kind.sigs.k8s.io/ |
| kubectl | v1.30.0+ | https://kubernetes.io/docs/tasks/tools/ |
| flux CLI | v2.3.0+ | https://fluxcd.io/flux/installation/ |
| Go | v1.23.0+ | https://go.dev/dl/ |

## Cluster Setup

### 1. Create the Kind cluster

```bash
kind create cluster --name sylva --wait 60s
```

### 2. Install FluxCD

```bash
flux install \
  --namespace=flux-system \
  --network-policy=false \
  --components=source-controller,kustomize-controller,helm-controller
```

Verify:

```bash
flux check
```

### 3. Deploy test workloads

```bash
kubectl apply -f cluster/apps/podinfo/namespace.yaml
kubectl apply -f cluster/apps/podinfo/
```

### 4. Wire the GitOps loop (optional)

```bash
kubectl apply -f cluster/flux/
```

This connects Flux to the GitHub repo so it continuously reconciles `cluster/apps/podinfo/`.

## Build

```bash
go build -o bin/sylvactl ./cmd/sylvactl/
```

## Usage

### List HelmReleases

```bash
# All namespaces
./bin/sylvactl get helmreleases

# Specific namespace
./bin/sylvactl get helmreleases -n podinfo
```

### List Kustomizations

```bash
# All namespaces
./bin/sylvactl get kustomizations

# Specific namespace
./bin/sylvactl get kustomizations -n flux-system
```

### Output formats

```bash
# Table (default)
./bin/sylvactl get helmreleases -n podinfo

# JSON
./bin/sylvactl get helmreleases -n podinfo -o json

# YAML
./bin/sylvactl get helmreleases -n podinfo -o yaml
```

### Global flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--namespace` | `-n` | all | Namespace to query |
| `--output` | `-o` | `table` | Output format: `table`, `json`, `yaml` |
| `--kubeconfig` | | `~/.kube/config` | Path to kubeconfig file |

## Teardown

```bash
kind delete cluster --name sylva
```
