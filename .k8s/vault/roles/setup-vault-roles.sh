#!/bin/bash

export VAULT_ADDR="http://vault.vault.svc:8200"
export VAULT_TOKEN="$VAULT_ROOT_TOKEN"

MICROSERVICES=("auth-service")

for SERVICE in "${MICROSERVICES[@]}"; do
  vault policy write "${SERVICE}-policy" - <<EOF
  path "secret/data/${SERVICE}/*" {
    capabilities = ["read", "list"]
  }
EOF

  vault write auth/kubernetes/role/"${SERVICE}-role" \
    bound_service_account_names="${SERVICE}-sa" \
    bound_service_account_namespaces="default" \
    policies="${SERVICE}-policy" \
    ttl="1h"
done