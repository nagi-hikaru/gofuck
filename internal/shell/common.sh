_gofuck_is_fix_args() {
    while [ "$#" -gt 0 ]; do
        case "$1" in
            -y|--yes|-y=true|--yes=true|-y=false|--yes=false) ;;
            -m|--model)
                [ "$#" -ge 2 ] || return 1
                case "$2" in
                    ''|-*) return 1 ;;
                esac
                shift
                ;;
            -m=?*|--model=?*) ;;
            *) return 1 ;;
        esac
        shift
    done
    return 0
}

fuck() {
    local previous corrected output_file capture_mask

    if ! _gofuck_is_fix_args "$@"; then
        command fuck "$@"
        return $?
    fi

    previous=$(_gofuck_previous) || {
        printf '%s\n' 'No recent command to correct.' >&2
        return 1
    }

    output_file="$HOME/.config/gofuck/output.txt"
    (
        capture_mask=$(umask)
        umask 077
        mkdir -p -- "${output_file%/*}" || exit 1
        {
            # 原命令失败时仍保留 stderr，让 Go 继续生成建议。
            (
                umask "$capture_mask"
                export LC_ALL=C
                builtin eval "$previous"
            ) >/dev/null || :
        } 2>| "$output_file"
    ) || return $?

    corrected=$(
        GOFUCK_SCRIPT="$previous" \
        GOFUCK_OUTPUT_FILE="$output_file" command fuck "$@"
    ) || return $?
    [ -n "$corrected" ] || return 1
    eval "$corrected"
}
