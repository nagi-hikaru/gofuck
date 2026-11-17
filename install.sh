#!/bin/sh
set -eu

src=${1:-./bin/fuck}
bin="$HOME/.local/bin"
rc="${ZDOTDIR:-$HOME}/.zshrc"

if [ ! -f "$src" ]; then
    printf '找不到二进制：%s\n' "$src" >&2
    exit 1
fi

mkdir -p "$bin"
install -m 755 -- "$src" "$bin/fuck"

if ! grep -Fqx '# gofuck shell integration' "$rc" 2>/dev/null; then
    cat >> "$rc" <<'EOF'

# gofuck shell integration
case ":$PATH:" in
    *":$HOME/.local/bin:"*) ;;
    *) export PATH="$HOME/.local/bin:$PATH" ;;
esac
eval "$(command fuck init zsh)"
# end gofuck shell integration
EOF
fi

printf '%s\n' '安装完成，请重新打开终端。'