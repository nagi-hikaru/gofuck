#!/bin/sh
set -eu

rc="${ZDOTDIR:-$HOME}/.zshrc"

if [ -f "$rc" ] && grep -Fqx '# gofuck shell integration' "$rc"; then
    tmp=$(mktemp)
    trap 'rm -f "$tmp"' 0

    awk '
        $0 == "# gofuck shell integration"     { skip = 1; next }
        $0 == "# end gofuck shell integration" && skip { skip = 0; next }
        !skip
        END { if (skip) exit 1 }
    ' "$rc" > "$tmp"

    cat "$tmp" > "$rc"
fi

rm -f -- "$HOME/.local/bin/fuck"
printf '%s\n' '卸载完成，请重新打开终端。'