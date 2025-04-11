#!/bin/bash

export VAULT_ADDR="$VAULT_ADDR"
export VAULT_TOKEN="$VAULT_ROOT_TOKEN"

MICROSERVICES=("auth-service")

for SERVICE in "${MICROSERVICES[@]}"; do
  vault policy write "${SERVICE}-policy" - <<EOF
  path "secret/data/${SERVICE}/*" {
    capabilities = ["create", "read", "update", "delete", "list"]
  }

  path "secret/metadata/${SERVICE}/*" {
    capabilities = ["list"]
  }
EOF

  vault write auth/kubernetes/role/"${SERVICE}-role" \
    bound_service_account_names="${SERVICE}-sa" \
    bound_service_account_namespaces="default" \
    policies="${SERVICE}-policy" \
    ttl="1h"
done