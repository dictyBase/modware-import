# Justfile for content

namespace      := "dictybase"
github_user    := "sba964"
platform       := "linux/amd64"
platform_multi := "linux/amd64,linux/arm64"
image          := namespace + "/modware-dataloader"
ghcr_image     := "ghcr.io/" + image

set default-list := true

# Build the docker image for the target platform
build tag="latest":
    docker buildx build --platform {{ platform }} -f build/package/Dockerfile.dataloader -t {{ image }}:{{ tag }} .

# Build and push the docker image
push tag="latest":
    docker buildx build --platform {{ platform }} -f build/package/Dockerfile.dataloader -t {{ image }}:{{ tag }} --push .

# Build for GitHub Container Registry
build-ghcr tag="latest":
    docker buildx build --platform {{ platform }} -f build/package/Dockerfile.dataloader -t {{ ghcr_image }}:{{ tag }} .

# Push to GitHub Container Registry
push-ghcr tag="latest":
    echo $GITHUB_REGISTRY_TOKEN | docker login ghcr.io -u {{ github_user }} --password-stdin
    docker buildx build --platform {{ platform }} -f build/package/Dockerfile.dataloader -t {{ ghcr_image }}:{{ tag }} --push .

# Build and push multi-arch image (amd64 + arm64)
push-multi tag="latest":
    docker buildx build --platform {{ platform_multi }} -f build/package/Dockerfile.dataloader -t {{ image }}:{{ tag }} --push .

# Push multi-arch image to GitHub Container Registry
push-ghcr-multi tag="latest":
    echo $GITHUB_REGISTRY_TOKEN | docker login ghcr.io -u {{ github_user }} --password-stdin
    docker buildx build --platform {{ platform_multi }} -f build/package/Dockerfile.dataloader -t {{ ghcr_image }}:{{ tag }} --push .

# List images
list-images:
    docker images | grep {{ image }}

load tag k8s_config s3_bucket_path s3_bucket="dictybase" k8s_namespace="dev" debug="false":
    #!/usr/bin/env bash
    set -euo pipefail
    ttl="{{ if debug == "true" { "300" } else { "120" } }}"
    gen_name="{{ if debug == "true" { "load-content-debug-" } else { "load-content-" } }}"
    container_name="{{ if debug == "true" { "load-content-debug" } else { "load-content" } }}"
    debug_args=""
    if [ "{{ debug }}" == "true" ]; then
        debug_args="- --log-level
                - debug
                - --log-format
                - text"
    fi
    kubectl create -f - --kubeconfig {{ k8s_config }} -o jsonpath='{.metadata.name}' <<EOF
    apiVersion: batch/v1
    kind: Job
    metadata:
      generateName: ${gen_name}
      namespace: {{ k8s_namespace }}
    spec:
      ttlSecondsAfterFinished: ${ttl}
      template:
        spec:
          restartPolicy: Never
          containers:
            - name: ${container_name}
              image: {{ ghcr_image }}:{{ tag }}
              envFrom:
                - secretRef:
                    name: minio
              args:
                - content-data
                ${debug_args}
                - --s3-bucket
                - {{ s3_bucket }}
                - --s3-bucket-path
                - {{ s3_bucket_path }}
    EOF
