#compdef stuff
# zsh completion for stuff
# Install: cp stuff.zsh /usr/local/share/zsh/site-functions/_stuff

_stuff() {
    local -a global_opts
    global_opts=(
        '--url[API base URL]:url'
        '--api-key[API key for authentication]:key'
        '--json[Output raw JSON]'
        '--help[Show help]'
    )

    _arguments -C \
        $global_opts \
        '1:command:->command' \
        '*::arg:->args'

    case $state in
        command)
            local -a commands
            commands=(
                'item:Manage items'
                'image:Manage images'
                'tag:Manage tags'
            )
            _describe 'command' commands
            ;;
        args)
            case $words[1] in
                item) _stuff_item ;;
                image) _stuff_image ;;
                tag) _stuff_tag ;;
            esac
            ;;
    esac
}

_stuff_item() {
    _arguments -C \
        '1:subcommand:->subcmd' \
        '*::arg:->args'

    case $state in
        subcmd)
            local -a subcmds
            subcmds=(
                'list:List items'
                'add:Create a new item'
                'show:Show item details'
                'edit:Update an item'
                'delete:Delete an item'
            )
            _describe 'subcommand' subcmds
            ;;
        args)
            case $words[1] in
                list)
                    _arguments \
                        {-q,--query}'[Search query]:query' \
                        '--tag[Filter by tag]:tag' \
                        '--condition[Filter by condition]:condition'
                    ;;
                add)
                    _arguments \
                        '--name[Item name (required)]:name' \
                        '--description[Description]:description' \
                        '--brand[Brand]:brand' \
                        '--model[Model]:model' \
                        '--serial[Serial number]:serial' \
                        '--purchase-date[Purchase date (YYYY-MM-DD)]:date' \
                        '--purchase-price[Purchase price]:price' \
                        '--estimated-value[Estimated value]:value' \
                        '--condition[Condition]:condition' \
                        '--notes[Notes]:notes' \
                        '*--tag[Tags (repeatable)]:tag'
                    ;;
                show)
                    _arguments '1:item-id'
                    ;;
                edit)
                    _arguments \
                        '1:item-id' \
                        '--name[Item name]:name' \
                        '--description[Description]:description' \
                        '--brand[Brand]:brand' \
                        '--model[Model]:model' \
                        '--serial[Serial number]:serial' \
                        '--purchase-date[Purchase date (YYYY-MM-DD)]:date' \
                        '--purchase-price[Purchase price]:price' \
                        '--estimated-value[Estimated value]:value' \
                        '--condition[Condition]:condition' \
                        '--notes[Notes]:notes' \
                        '*--tag[Tags (replaces all)]:tag'
                    ;;
                delete)
                    _arguments \
                        '1:item-id' \
                        {-y,--yes}'[Skip confirmation]'
                    ;;
            esac
            ;;
    esac
}

_stuff_image() {
    _arguments -C \
        '1:subcommand:->subcmd' \
        '*::arg:->args'

    case $state in
        subcmd)
            local -a subcmds
            subcmds=(
                'add:Upload an image'
                'delete:Delete an image'
            )
            _describe 'subcommand' subcmds
            ;;
        args)
            case $words[1] in
                add)
                    _arguments \
                        '1:item-id' \
                        '2:file:_files'
                    ;;
                delete)
                    _arguments '1:image-id'
                    ;;
            esac
            ;;
    esac
}

_stuff_tag() {
    _arguments -C \
        '1:subcommand:->subcmd'

    case $state in
        subcmd)
            local -a subcmds
            subcmds=('list:List all tags')
            _describe 'subcommand' subcmds
            ;;
    esac
}

_stuff "$@"
