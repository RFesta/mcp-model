#!/usr/bin/env bash
set -euo pipefail

if [ $# -lt 1 ]; then
  echo "uso: $0 novo-nome-mcp"
  exit 1
fi

OLD="modelo-mcp"
NEW="$1"

# Renomeia pastas/arquivos e substitui conteúdo
grep -rl "$OLD" . | xargs sed -i "s/$OLD/$NEW/g"
mv "cmd/$OLD" "cmd/$NEW" 2>/dev/null || true
echo "OK: base renomeada para $NEW"
