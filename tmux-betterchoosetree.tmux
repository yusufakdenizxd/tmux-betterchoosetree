#!/usr/bin/env bash

tmux_get() {
    local value="$(tmux show -gqv "$1")"
    [ -n "$value" ] && echo "$value" || echo "$2"
}

tmux_set() {
    tmux set-option -gq "$1" "$2"
}

tmux bind-key w popup -B -w100% -h99% -E "./go/main"
