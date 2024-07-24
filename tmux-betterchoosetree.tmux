#!/usr/bin/env bash

CURRENT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

tmux_get() {
    local value="$(tmux show -gqv "$1")"
    [ -n "$value" ] && echo "$value" || echo "$2"
}

tmux_set() {
    tmux set-option -gq "$1" "$2"
}
tmux bind-key W popup -B -w100% -h99% -E "$CURRENT_DIR/go/main"
