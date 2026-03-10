# Justfile for goldenbraid

name := "goldenbraid"
namespace := "dictybase"
github_user := "sba964"
platform := "linux/amd64"
platform_multi := "linux/amd64,linux/arm64"
image := namespace + "/" + name
ghcr_image := "ghcr.io/" + image

[private]
check-kubeconfig:
    @echo "Using KUBECONFIG={{ env('KUBECONFIG') }}" >&2

# Build the docker image for the target platform
build tag="latest":
    docker buildx build --platform {{ platform }} -f build/package/Dockerfile.goldenbraid -t {{ image }}:{{ tag }} .

# Build and push the docker image
push tag="latest":
    docker buildx build --platform {{ platform }} -f build/package/Dockerfile.goldenbraid -t {{ image }}:{{ tag }} --push .

# Build for GitHub Container Registry
build-ghcr tag="latest":
    docker buildx build --platform {{ platform }} -f build/package/Dockerfile.goldenbraid -t {{ ghcr_image }}:{{ tag }} .

# Push to GitHub Container Registry
push-ghcr tag="latest":
    echo $GITHUB_REGISTRY_TOKEN | docker login ghcr.io -u {{ github_user }} --password-stdin
    docker buildx build --platform {{ platform }} -f build/package/Dockerfile.goldenbraid -t {{ ghcr_image }}:{{ tag }} --push .

# Build and push multi-arch image (amd64 + arm64)
push-multi tag="latest":
    docker buildx build --platform {{ platform_multi }} -f build/package/Dockerfile.goldenbraid -t {{ image }}:{{ tag }} --push .

# Push multi-arch image to GitHub Container Registry
push-ghcr-multi tag="latest":
    echo $GITHUB_REGISTRY_TOKEN | docker login ghcr.io -u {{ github_user }} --password-stdin
    docker buildx build --platform {{ platform_multi }} -f build/package/Dockerfile.goldenbraid -t {{ ghcr_image }}:{{ tag }} --push .

# Show parameters for a Dagu workflow
dagu-params file:
    #!/usr/bin/env bash
    set -euo pipefail
    schema=$(yq -r '.params.schema' {{ file }})
    jq -r '.properties | to_entries | .[] | .key + " = " + (.value.default // "(required)" | tostring) + "  # " + (.value.description // "")' "$schema"

# List images
list:
    docker images | grep {{ image }}

# Run goldenbraid plasmid import job in dev cluster
run-goldenbraid tag email k8s_namespace="dev" debug="false": check-kubeconfig
    #!/usr/bin/env bash
    set -euo pipefail
    ttl="{{ if debug == "true" { "300" } else { "120" } }}"
    gen_name="{{ if debug == "true" { "goldenbraid-plasmid-debug-" } else { "goldenbraid-plasmid-" } }}"
    container_name="{{ if debug == "true" { "goldenbraid-plasmid-debug" } else { "goldenbraid-plasmid" } }}"
    debug_args=""
    if [ "{{ debug }}" == "true" ]; then
        debug_args="- --log-level
                - debug
                - --log-format
                - text"
    fi
    kubectl create -f - -o jsonpath='{.metadata.name}' <<EOF
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
                - plasmid
                - --user-email
                - {{ email }}
                ${debug_args}
    EOF

# Run goldenbraid plasmid-ontology job in dev cluster (assigns ontology term to all plasmids)
run-goldenbraid-plasmid-ontology tag ontology_term="vector" k8s_namespace="dev" debug="false": check-kubeconfig
    #!/usr/bin/env bash
    set -euo pipefail
    ttl="{{ if debug == "true" { "300" } else { "120" } }}"
    gen_name="{{ if debug == "true" { "goldenbraid-plasmid-ontology-debug-" } else { "goldenbraid-plasmid-ontology-" } }}"
    container_name="{{ if debug == "true" { "goldenbraid-plasmid-ontology-debug" } else { "goldenbraid-plasmid-ontology" } }}"
    debug_args=""
    if [ "{{ debug }}" == "true" ]; then
        debug_args="- --log-level
                - debug
                - --log-format
                - text"
    fi
    kubectl create -f - -o jsonpath='{.metadata.name}' <<EOF
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
              args:
                - plasmid-ontology
                - --ontology-term
                - {{ ontology_term }}
                ${debug_args}
    EOF

# Look up a GoldenBraid plasmid by exact name (uses goldenbraid-list image)
lookup-plasmid tag name k8s_namespace="dev": check-kubeconfig
    #!/usr/bin/env bash
    set -euo pipefail
    kubectl create -f - -o jsonpath='{.metadata.name}' <<EOF
    apiVersion: batch/v1
    kind: Job
    metadata:
      generateName: goldenbraid-lookup-
      namespace: {{ k8s_namespace }}
    spec:
      ttlSecondsAfterFinished: 120
      template:
        spec:
          restartPolicy: Never
          containers:
            - name: goldenbraid-lookup
              image: ghcr.io/dictybase/goldenbraid-list:{{ tag }}
              env:
                - name: PLASMID_NAME
                  value: "{{ name }}"
              envFrom:
                - secretRef:
                    name: minio
              args:
                - lookup
    EOF

# Run goldenbraid inventory import job in dev cluster
run-goldenbraid-inventory tag k8s_namespace="dev" debug="false": check-kubeconfig
    #!/usr/bin/env bash
    set -euo pipefail
    ttl="{{ if debug == "true" { "300" } else { "120" } }}"
    gen_name="{{ if debug == "true" { "goldenbraid-inventory-debug-" } else { "goldenbraid-inventory-" } }}"
    container_name="{{ if debug == "true" { "goldenbraid-inventory-debug" } else { "goldenbraid-inventory" } }}"
    debug_args=""
    if [ "{{ debug }}" == "true" ]; then
        debug_args="- --log-level
                - debug
                - --log-format
                - text"
    fi
    kubectl create -f - -o jsonpath='{.metadata.name}' <<EOF
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
                - inventory
                ${debug_args}
    EOF

# Wait for a Kubernetes job to complete, fail, or detect stuck pods.

# Delegates to the k8s wait-job subcommand
wait-job name k8s_namespace="dev" timeout="60s": check-kubeconfig
    #!/usr/bin/env bash
    set -euo pipefail
    go run ./cmd/k8s/ wait-job --name {{ name }} --namespace {{ k8s_namespace }} --timeout {{ timeout }} 

# Get the logs for a specific job
job-logs name k8s_namespace="dev": check-kubeconfig
    #!/usr/bin/env bash
    set -euo pipefail
    go run ./cmd/k8s/ job-logs --name {{ name }} --namespace {{ k8s_namespace }} --follow

# Get failure details for a job
job-debug name k8s_namespace="dev": check-kubeconfig
    #!/usr/bin/env bash
    echo "--- Pod Logs ---"
    go run ./cmd/k8s/ job-logs --name {{ name }} --namespace {{ k8s_namespace }} || true
    echo "--- Job Description ---"
    kubectl describe job/{{ name }} -n {{ k8s_namespace }}
