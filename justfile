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

# Run goldenbraid plasmid import with debug logging (diagnostic)
run-goldenbraid-debug tag email:
    #!/usr/bin/env bash
    set -euo pipefail
    export KUBECONFIG=$(k3d kubeconfig write k3d-dev-cluster)
    kubectl apply -f - <<EOF
    apiVersion: batch/v1
    kind: Job
    metadata:
      name: goldenbraid-plasmid-debug
      namespace: dev
    spec:
      ttlSecondsAfterFinished: 300
      template:
        spec:
          restartPolicy: Never
          containers:
            - name: goldenbraid-plasmid-debug
              image: {{ghcr_image}}:{{tag}}
              envFrom:
                - secretRef:
                    name: minio
              args:
                - plasmid
                - --user-email
                - {{email}}
                - --log-level
                - debug
                - --log-format
                - text
    EOF

# Run goldenbraid plasmid-ontology job in dev cluster (assigns ontology term to all plasmids)
run-goldenbraid-plasmid-ontology tag ontology_term="vector":
    #!/usr/bin/env bash
    set -euo pipefail
    export KUBECONFIG=$(k3d kubeconfig write k3d-dev-cluster)
    kubectl create -f - <<EOF
    apiVersion: batch/v1
    kind: Job
    metadata:
      generateName: goldenbraid-plasmid-ontology-
      namespace: dev
    spec:
      ttlSecondsAfterFinished: 120
      template:
        spec:
          restartPolicy: Never
          containers:
            - name: goldenbraid-plasmid-ontology
              image: {{ghcr_image}}:{{tag}}
              args:
                - plasmid-ontology
                - --ontology-term
                - {{ontology_term}}
    EOF

# Run goldenbraid plasmid-ontology with debug logging
run-goldenbraid-plasmid-ontology-debug tag ontology_term="vector":
    #!/usr/bin/env bash
    set -euo pipefail
    export KUBECONFIG=$(k3d kubeconfig write k3d-dev-cluster)
    kubectl apply -f - <<EOF
    apiVersion: batch/v1
    kind: Job
    metadata:
      name: goldenbraid-plasmid-ontology-debug
      namespace: dev
    spec:
      ttlSecondsAfterFinished: 300
      template:
        spec:
          restartPolicy: Never
          containers:
            - name: goldenbraid-plasmid-ontology-debug
              image: {{ghcr_image}}:{{tag}}
              args:
                - plasmid-ontology
                - --ontology-term
                - {{ontology_term}}
                - --log-level
                - debug
                - --log-format
                - text
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

# Run goldenbraid inventory import job in dev cluster
run-goldenbraid-inventory tag:
    #!/usr/bin/env bash
    set -euo pipefail
    export KUBECONFIG=$(k3d kubeconfig write k3d-dev-cluster)
    kubectl apply -f - <<EOF
    apiVersion: batch/v1
    kind: Job
    metadata:
      name: goldenbraid-inventory
      namespace: dev
    spec:
      ttlSecondsAfterFinished: 120
      template:
        spec:
          restartPolicy: Never
          containers:
            - name: goldenbraid-inventory
              image: {{ghcr_image}}:{{tag}}
              envFrom:
                - secretRef:
                    name: minio
              args:
                - inventory
    EOF

# Run goldenbraid inventory import with debug logging
run-goldenbraid-inventory-debug tag:
    #!/usr/bin/env bash
    set -euo pipefail
    export KUBECONFIG=$(k3d kubeconfig write k3d-dev-cluster)
    kubectl apply -f - <<EOF
    apiVersion: batch/v1
    kind: Job
    metadata:
      name: goldenbraid-inventory-debug
      namespace: dev
    spec:
      ttlSecondsAfterFinished: 120
      template:
        spec:
          restartPolicy: Never
          containers:
            - name: goldenbraid-inventory-debug
              image: {{ghcr_image}}:{{tag}}
              envFrom:
                - secretRef:
                    name: minio
              args:
                - inventory
                - --log-level
                - debug
                - --log-format
                - text
    EOF

# Wait for a Kubernetes job to complete, fail, or detect stuck pods.
# Delegates to the goldenbraid wait-job subcommand
wait-job name namespace="dev" timeout="60s":
    #!/usr/bin/env bash
    set -euo pipefail
    kubeconfig=$(k3d kubeconfig write k3d-dev-cluster)
    go run ./cmd/goldenbraid/ wait-job --name {{name}} --namespace {{namespace}} --timeout {{timeout}} --kubeconfig "$kubeconfig"

# Get the logs for a specific job
job-logs name namespace="dev":
    #!/usr/bin/env bash
    export KUBECONFIG=$(k3d kubeconfig write k3d-dev-cluster)
    kubectl logs job/{{name}} -n {{namespace}}

# Get failure details for a job
job-debug name namespace="dev":
    #!/usr/bin/env bash
    export KUBECONFIG=$(k3d kubeconfig write k3d-dev-cluster)
    echo "--- Pod Logs ---"
    kubectl logs job/{{name}} -n {{namespace}} || true
    echo "--- Job Description ---"
    kubectl describe job/{{name}} -n {{namespace}}
