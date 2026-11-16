# 加载 zsh 的 parameter 模块
zmodload zsh/parameter
_gofuck_previous() {
    emulate -L zsh
    local index line
    local -i scanned=0
    # 按历史编号从新到旧检查，最多检查 10 条
    for index in ${(Onk)history}; do
        (( ++scanned <= 10 )) || break
        line=${history[$index]}
        # 跳过空白记录和 fuck 自身
        case ${line#${line%%[![:space:]]*}} in
            ''|fuck|fuck[[:space:]]*) continue ;;
        esac
        print -r -- "$line"
        return 0
    done
    return 1
}

# 加载 common.sh
source "${${(%):-%N}:A:h}/common.sh"
