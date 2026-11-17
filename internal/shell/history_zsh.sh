# 加载 zsh 的 parameter 模块，获取 history 数组
zmodload zsh/parameter

_gofuck_previous() {
    emulate -L zsh
    local index line    # 声明局部变量
    local -i scanned=0  # 声明从 0 开始的整数计数器
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
