_gofuck_previous() {
    local latest index offset line

    # 取出最新命令的编号
    latest=$(HISTTIMEFORMAT= builtin history 1) || return 1
    if [[ $latest =~ ^[[:blank:]]*([0-9]+) ]]; then
        index=${BASH_REMATCH[1]}
    else
        return 1
    fi

    for ((offset=0; offset<10 && index>0; offset++, index--)); do
        # 取出历史命令
        line=$(builtin fc -ln "$index" "$index" 2>/dev/null) || break
        line=${line#$'\t'}
        line=${line# }

        # 跳过空白和 fuck 命令本身
        case ${line#${line%%[![:space:]]*}} in
            ''|fuck|fuck[[:space:]]*) continue ;;
        esac
        printf '%s\n' "$line"
        return 0
    done
    return 1
}
