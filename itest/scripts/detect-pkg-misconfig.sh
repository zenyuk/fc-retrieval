#!/bin/bash
echo "*************************************************"
echo "*** Check that go.mod is configured correctly ***"
echo "*************************************************"

cp ../go.mod ../go.mod.temp
bash ./use-remote-repos.sh
go mod tidy

if cmp -s ../go.mod ../go.mod.temp; then
    echo go.mod file correctly configured.
else 
    echo ERROR: go.mod file not configured correctly. 
    echo Please execute:
    echo   bash use-remote-repos.sh
    echo   go mod tidy
    exit 1
fi

