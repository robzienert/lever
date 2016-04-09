#!/usr/bin/env bash -e
echo "" > coverage.txt

for d in $(find . -type d -not -path "./vendor*" -and -not -path "./.*" -maxdepth 10); do
	if ls $d/*.go &> /dev/null; then
		go test -v -race -coverprofile=profile.out -covermode=atomic $d
		if [ -f profile.out ]; then
			cat profile.out >> coverage.txt
			rm profile.out
		fi
	fi
done