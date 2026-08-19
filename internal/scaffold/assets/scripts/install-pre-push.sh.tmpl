#!/bin/sh
set -eu

repo_root=$(git rev-parse --show-toplevel)
hook_path="$repo_root/.git/hooks/pre-push"
temporary_hook="$hook_path.gh-pr-quality-gate"

cat > "$temporary_hook" <<'EOF'
#!/bin/sh
# Managed by gh-pr-quality-gate
exec gh pr-quality-gate validate --repo "$(git rev-parse --show-toplevel)" --run-local
EOF

if [ -e "$hook_path" ] && ! cmp -s "$temporary_hook" "$hook_path"; then
  rm -f "$temporary_hook"
  echo "A different pre-push hook already exists at $hook_path. Integrate it manually." >&2
  exit 1
fi

mv "$temporary_hook" "$hook_path"
chmod +x "$hook_path"
echo "Installed $hook_path"
