mod goldenbraid "just-modules/goldenbraid.just"
mod content "just-modules/content.just"

set default-list := true

# Wait for a Kubernetes job to complete, fail, or detect stuck pods.

# Delegates to the k8s wait-job subcommand
wait-job name k8s_config k8s_namespace="dev" timeout="60s":
    #!/usr/bin/env bash
    set -euo pipefail
    go run ./cmd/k8s/ wait-job --name {{ name }} --kubeconfig {{  k8s_config }} --namespace {{ k8s_namespace }} --timeout {{ timeout }} 

# Get the logs for a specific job
job-logs name k8s_config k8s_namespace="dev":
    #!/usr/bin/env bash
    set -euo pipefail
    go run ./cmd/k8s/ job-logs --name {{ name }} --kubeconfig {{ k8s_config }} --namespace {{ k8s_namespace }} --follow

# Get failure details for a job
job-debug name k8s_config k8s_namespace="dev": 
    #!/usr/bin/env bash
    echo "--- Pod Logs ---"
    go run ./cmd/k8s/ job-logs --name {{ name }} --kubeconfig {{ k8s_config }} --namespace {{ k8s_namespace }} || true
    echo "--- Job Description ---"
    kubectl describe job/{{ name }} --kubeconfig {{ k8s_config }} -n {{ k8s_namespace }}
