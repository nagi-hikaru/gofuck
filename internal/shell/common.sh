fuck() {
    local previous corrected output_file capture_mask

    case "$#:${1-}" in
        0:|1:-y|1:--yes|1:-y=true|1:--yes=true|1:-y=false|1:--yes=false) ;;
        *) command fuck "$@"; return $? ;;
    esac

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
