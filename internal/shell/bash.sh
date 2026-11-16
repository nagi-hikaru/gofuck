_gofuck_shell=bash

_gofuck_previous() {
    local latest index offset line
    latest=$(HISTTIMEFORMAT= builtin history 1) || return 1
    if [[ $latest =~ ^[[:blank:]]*([0-9]+) ]]; then
        index=${BASH_REMATCH[1]}
    else
        return 1
    fi

    for ((offset=0; offset<10 && index>0; offset++, index--)); do
        line=$(builtin fc -ln "$index" "$index" 2>/dev/null) || break
        line=${line#$'\t'}
        line=${line# }
        case ${line#${line%%[![:space:]]*}} in
            ''|fuck|fuck[[:space:]]*) continue ;;
        esac
        printf '%s\n' "$line"
        return 0
    done
    return 1
}

source "$(dirname -- "${BASH_SOURCE[0]}")/common.sh"