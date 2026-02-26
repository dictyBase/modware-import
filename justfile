# Justfile for goldenbraid

name := "goldenbraid"
namespace := "dictybase"
github_user := "sba964"
platform := "linux/amd64"
platform_multi := "linux/amd64,linux/arm64"

image := namespace + "/" + name
ghcr_image := "ghcr.io/" + image

# Build the docker image for the target platform
build tag="latest":
    docker buildx build --platform {{platform}} -f build/package/Dockerfile.goldenbraid -t {{image}}:{{tag}} .

# Build and push the docker image
push tag="latest":
    docker buildx build --platform {{platform}} -f build/package/Dockerfile.goldenbraid -t {{image}}:{{tag}} --push .

# Build for GitHub Container Registry
build-ghcr tag="latest":
    docker buildx build --platform {{platform}} -f build/package/Dockerfile.goldenbraid -t {{ghcr_image}}:{{tag}} .

# Push to GitHub Container Registry
push-ghcr tag="latest":
    echo $GITHUB_REGISTRY_TOKEN | docker login ghcr.io -u {{github_user}} --password-stdin
    docker buildx build --platform {{platform}} -f build/package/Dockerfile.goldenbraid -t {{ghcr_image}}:{{tag}} --push .

# Build and push multi-arch image (amd64 + arm64)
push-multi tag="latest":
    docker buildx build --platform {{platform_multi}} -f build/package/Dockerfile.goldenbraid -t {{image}}:{{tag}} --push .

# Push multi-arch image to GitHub Container Registry
push-ghcr-multi tag="latest":
    echo $GITHUB_REGISTRY_TOKEN | docker login ghcr.io -u {{github_user}} --password-stdin
    docker buildx build --platform {{platform_multi}} -f build/package/Dockerfile.goldenbraid -t {{ghcr_image}}:{{tag}} --push .

# List images
list:
    docker images | grep {{image}}

# Run goldenbraid plasmid import job in dev cluster
run-goldenbraid tag email:
    #!/usr/bin/env bash
    set -euo pipefail
    export KUBECONFIG=$(k3d kubeconfig write k3d-dev-cluster)
    kubectl apply -f - <<EOF
    apiVersion: batch/v1
    kind: Job
    metadata:
      name: goldenbraid-plasmid
      namespace: dev
    spec:
      ttlSecondsAfterFinished: 120
      template:
        spec:
          restartPolicy: Never
          containers:
            - name: goldenbraid-plasmid
              image: {{ghcr_image}}:{{tag}}
              envFrom:
                - secretRef:
                    name: minio
              args:
                - plasmid
                - --user-email
                - {{email}}
    EOF

# Look up a GoldenBraid plasmid by exact name (uses goldenbraid-list image)
lookup-plasmid tag name:
    #!/usr/bin/env bash
    set -euo pipefail
    export KUBECONFIG=$(k3d kubeconfig write k3d-dev-cluster)
    kubectl apply -f - <<EOF
    apiVersion: batch/v1
    kind: Job
    metadata:
      name: goldenbraid-lookup
      namespace: dev
    spec:
      ttlSecondsAfterFinished: 120
      template:
        spec:
          restartPolicy: Never
          containers:
            - name: goldenbraid-lookup
              image: ghcr.io/dictybase/goldenbraid-list:{{tag}}
              env:
                - name: PLASMID_NAME
                  value: "{{name}}"
              envFrom:
                - secretRef:
                    name: minio
              args:
                - lookup
    EOF
