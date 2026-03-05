# bash completion for stuff
# Install: source stuff.bash
#   or: cp stuff.bash /usr/local/share/bash-completion/completions/stuff

_stuff() {
    local cur prev words cword
    _init_completion || return

    local commands="item image tag"
    local item_commands="list add show edit delete"
    local image_commands="add delete"
    local tag_commands="list"

    local global_flags="--url --api-key --json --help"

    # Determine command depth
    local cmd="" subcmd=""
    for ((i=1; i < cword; i++)); do
        case "${words[i]}" in
            item|image|tag)
                cmd="${words[i]}"
                ;;
            list|add|show|edit|delete)
                if [[ -n "$cmd" ]]; then
                    subcmd="${words[i]}"
                fi
                ;;
        esac
    done

    # Complete top-level commands
    if [[ -z "$cmd" ]]; then
        COMPREPLY=($(compgen -W "$commands $global_flags" -- "$cur"))
        return
    fi

    # Complete subcommands
    if [[ -z "$subcmd" ]]; then
        case "$cmd" in
            item)   COMPREPLY=($(compgen -W "$item_commands $global_flags" -- "$cur")) ;;
            image)  COMPREPLY=($(compgen -W "$image_commands $global_flags" -- "$cur")) ;;
            tag)    COMPREPLY=($(compgen -W "$tag_commands $global_flags" -- "$cur")) ;;
        esac
        return
    fi

    # Complete flags for subcommands
    case "$cmd/$subcmd" in
        item/list)
            COMPREPLY=($(compgen -W "-q --query --tag --condition $global_flags" -- "$cur"))
            ;;
        item/add)
            COMPREPLY=($(compgen -W "--name --description --brand --model --serial --purchase-date --purchase-price --estimated-value --condition --notes --tag $global_flags" -- "$cur"))
            ;;
        item/show)
            COMPREPLY=($(compgen -W "$global_flags" -- "$cur"))
            ;;
        item/edit)
            COMPREPLY=($(compgen -W "--name --description --brand --model --serial --purchase-date --purchase-price --estimated-value --condition --notes --tag $global_flags" -- "$cur"))
            ;;
        item/delete)
            COMPREPLY=($(compgen -W "-y --yes $global_flags" -- "$cur"))
            ;;
        image/add)
            # Second arg is a file path
            if [[ "$prev" != "add" ]]; then
                _filedir
            fi
            ;;
        image/delete)
            COMPREPLY=($(compgen -W "$global_flags" -- "$cur"))
            ;;
        *)
            COMPREPLY=($(compgen -W "$global_flags" -- "$cur"))
            ;;
    esac
}

complete -F _stuff stuff
