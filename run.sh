#!/bin/sh

set -e

help() {
    echo '  -g\trun go generate [PATH] default $(glide novendor)'
    echo '  -it\trun integration test'
    echo '  -t\trun unit test [PATH] default $(glide novendor)'
    echo '  -b\tcreate binary'
    echo '  -r\trun docker-compose'
    echo '  -i\trun godep save'
}

if [ -z $1 ]; then
  echo 'Please action'
  help
  exit 0
fi

case "$1" in
  -g)
    if [ -z $2 ]; then
        set -x
        go generate $(glide novendor)
    else
        shift
        set -x
        go generate $@
    fi
  ;;
  -it)
    if [ -z $2 ]; then
        set -x
        go test -v -tags integration $(glide novendor)
    else
        shift
        set -x
        go test -v -tags integration $@
    fi
  ;;
  -b)
    set -x
    go build -o app ./bin
  ;;
  -r)
    set -x
    go run *.go
  ;;
  -t)
    if [ -z $2 ]; then
        set -x
        go test -v $(glide novendor)
    else
        shift
        set -x
        go test -v $@
    fi
  ;;
  -h | --help)
    help
  ;;
  *)
    help
  ;;
esac

exit $?
