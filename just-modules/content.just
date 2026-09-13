# Justfile for content

namespace      := "dictybase"
github_user    := "sba964"
platform_multi := "linux/amd64,linux/arm64"
image          := namespace + "/modware-dataloader"
ghcr_image     := "ghcr.io/" + image

set default-list := true

# ─────────────────────────────────────────────
# Local image builds (macOS dev)
# ─────────────────────────────────────────────

# Build and push the multi-arch content loader image to GitHub Container Registry.
# Local macOS workflow only — CI publishes the Docker Hub image built from the same
# Dockerfile on push to develop.
# Usage: just content::push-ghcr-multi [--tag <tag>]
[arg("tag", long="tag", short="t", help="Image tag (default latest)")]
[group('content-images')]
[no-cd]
push-ghcr-multi tag="latest":
    #!/usr/bin/env bash
    set -euo pipefail

    root="{{ invocation_directory() }}"

    tag="{{ tag }}"
    if [ -z "${tag}" ]; then
        echo "ERROR: no tag. Pass --tag <tag>." >&2
        exit 1
    fi

    if [ -z "${GITHUB_REGISTRY_TOKEN:-}" ]; then
        echo "ERROR: GITHUB_REGISTRY_TOKEN is not set — required to log in to ghcr.io." >&2
        exit 1
    fi

    echo "${GITHUB_REGISTRY_TOKEN}" | docker login ghcr.io -u "{{ github_user }}" --password-stdin
    docker buildx build \
        --platform {{ platform_multi }} \
        -f "${root}/build/package/Dockerfile.dataloader" \
        -t "{{ ghcr_image }}:${tag}" \
        --push \
        "${root}"

# ─────────────────────────────────────────────
# Content data jobs
# ─────────────────────────────────────────────

# Launch a content-data import job in the target cluster.
# S3 credentials come from the minio secret keys user/pass; the content API host
# and port are injected by the content-api Service link in the job namespace.
# Usage: just content::load --tag <tag> --k8s-config <path> --s3-bucket-path <path>
#        [--s3-bucket dictybase] [--k8s-namespace dev] [--debug false]
[arg("tag", long="tag", short="t", help="Image tag to run (required)")]
[arg("k8s_config", long="k8s-config", short="c", help="Path to the kubeconfig file (required)")]
[arg("s3_bucket_path", long="s3-bucket-path", short="p", help="Path inside the S3 bucket for input content files (required)")]
[arg("s3_bucket", long="s3-bucket", short="b", help="S3 bucket name (default dictybase)")]
[arg("k8s_namespace", long="k8s-namespace", short="n", help="Kubernetes namespace (default dev)")]
[arg("debug", long="debug", short="d", pattern="true|false", help="Debug logging and a longer job TTL")]
[group('content-jobs')]
[no-cd]
load tag="" k8s_config="" s3_bucket_path="" s3_bucket="dictybase" k8s_namespace="dev" debug="false":
    #!/usr/bin/env bash
    set -euo pipefail

    tag="{{ tag }}"
    k8s_config="{{ k8s_config }}"
    s3_bucket_path="{{ s3_bucket_path }}"
    if [ -z "${tag}" ] || [ -z "${k8s_config}" ] || [ -z "${s3_bucket_path}" ]; then
        echo "ERROR: tag, kubeconfig and S3 bucket path are all required." >&2
        echo "Pass --tag <tag> --k8s-config <path> --s3-bucket-path <path>." >&2
        exit 1
    fi

    if [ "{{ debug }}" = "true" ]; then
        ttl="300"
        gen_name="load-content-debug-"
        container_name="load-content-debug"
        debug_args="- --log-level
                - debug
                - --log-format
                - text"
    else
        ttl="120"
        gen_name="load-content-"
        container_name="load-content"
        debug_args=""
    fi

    kubectl create -f - --kubeconfig "${k8s_config}" -o jsonpath='{.metadata.name}' <<EOF
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
              image: {{ ghcr_image }}:${tag}
              env:
                - name: ACCESS_KEY
                  valueFrom:
                    secretKeyRef:
                      name: minio
                      key: user
                - name: SECRET_KEY
                  valueFrom:
                    secretKeyRef:
                      name: minio
                      key: pass
              args:
                ${debug_args}
                - content-data
                - --s3-bucket
                - {{ s3_bucket }}
                - --s3-bucket-path
                - ${s3_bucket_path}
    EOF