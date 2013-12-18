sorthead
======

## ABOUT

This command is the replacement for `sort ... |head` which works
orders of magnitude faster on large amount of data.

## USAGE

	Usage of sorthead:
	  -N=10: print the first N lines instead of the first 10
	  -k=0: sort by field number N, not the whole string
	  -n=false: compare according to string numerical value
	  -r=false: reverse the result of comparisons

## INSTALLATION

* install `golang`
* run `go build` in this project's directory
* install `sorthead` to `/usr/bin`
