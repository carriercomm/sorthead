sorthead
======

## ABOUT

Faster replacement for `sort ... | head`.
Works orders of magnitude faster on large amount of data.
Shows partial results (the current top) while working.

## USAGE

	Usage of sorthead:
	  -I, --interactive=false: interactive mode (it is the default when no -N given)
	  -k, --key=0: sort by field number N, not the whole string
	  -N, --lines=10: print the first N lines instead of the first 10 (in interactive mode default is window size)
	  -n, --numeric-sort=false: compare according to string numerical value
	  -r, --reverse=false: reverse the result of comparisons

## INSTALLATION

* install `golang`
* setup `GOPATH` as described here: http://golang.org/doc/code.html#GOPATH
  * `mkdir $HOME/go`
  * `export GOPATH=$HOME/go`
* install required libraries (they will be installed into `$GOPATH`):
  * `go get github.com/nsf/termbox-go`
  * `go get github.com/ogier/pflag`
* run `go build` in this project's directory
* (optional) `strip sorthead`
* install `sorthead` to `/usr/bin`
* (optional) now you can uninstall `golang` and run `rm -rf $GOPATH`
