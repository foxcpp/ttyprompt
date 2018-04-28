#!/bin/sh
exec ttyprompt -t 23 --title "SSH password request" -d "" --prompt "$1"
