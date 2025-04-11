#!/bin/bash

MICROSERVICES=("auth-service")

for SERVICE in "${MICROSERVICES[@]}"; do
  kubectl exec -n vault vault-0 -- vault policy write "${SERVICE}-policy" - <<EOF
path "secret/data/${SERVICE}/*" {
  capabilities = ["read", "list"]
}
EOF

  kubectl exec -n vault vault-0 -- vault write auth/kubernetes/role/"${SERVICE}-role" \
    bound_service_account_names="${SERVICE}-sa" \
    bound_service_account_namespaces="default" \
    policies="${SERVICE}-policy" \
    ttl="1h"
done