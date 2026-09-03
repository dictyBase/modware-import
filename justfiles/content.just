# Justfile for content

namespace      := "dictybase"
github_user    := "sba964"
platform       := "linux/amd64"
platform_multi := "linux/amd64,linux/arm64"
image          := namespace + "/modware-import"
ghcr_image     := "ghcr.io/" + image

set default-list := true

load-content tag k8s_config s3_bucket_path s3_bucket="dictybase" k8s_namespace="dev" debug="false":
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
                - --s3-bucket
                - {{ s3_bucket }}
                - --s3-bucket-path
                - {{ s3_bucket_path }}
                ${debug_args}
    EOF
