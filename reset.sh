#!/usr/bin/env bash
set -e

case "$1" in
	up)
	cd sql/schema && goose postgres "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable" up
	;;
	down)
	cd sql/schema && goose postgres "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable" down
	;;
	*)
	echo "Usage: $0 {up|down}"
	exit 1
	;;
esac
