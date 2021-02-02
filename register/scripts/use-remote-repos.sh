#!/bin/bash
set -e
echo "*******************************************************************************"
echo "*** Update go.mod and go.sum to point to the latest packages on github.com  ***"
echo "*** for gateway.                                                            ***"
echo "*******************************************************************************"

cd ..

# Remove any local references
sed '/replace .*/d' go.mod > go.mod.temp1
rm go.mod
mv go.mod.temp1 go.mod

REG_DIR="../fc-retrieval-register"
REG_BRANCH=`git rev-parse --abbrev-ref HEAD`
echo "register repo branch: $REG_BRANCH"


# Check repo:
# $1 is the name of the repo. fc-retrieval-client would be client.
# $2 is the relative directory of the repo
#
check_repo() {
    # Check directory of repo exists.
    if [ ! -d "$2" ]; then
        echo "ERROR: $1 directory $2 does not exist"
        exit 1
    fi
    echo "Found $1 repo: $2"

    cd $2
    OTHER_REPO_BRANCH=`git rev-parse --abbrev-ref HEAD`
    echo "$1 repo branch: $OTHER_REPO_BRANCH"
    BRANCH_EXISTS=`git branch -r --list origin/$REG_BRANCH`

    cd $REG_DIR

    if [ $REG_BRANCH != $OTHER_REPO_BRANCH ]; then 
        echo "register and $1 branch do not match"
        if [ -n "$BRANCH_EXISTS" ]; then
            echo "ERROR: Branch $REG_BRANCH exists on the $1 repo, but the $1 repo is currently using branch $OTHER_REPO_BRANCH"
            exit 1
        fi

        echo "Calling go get to use main on $1"
        go get github.com/ConsenSys/fc-retrieval-$1@main
    else
        echo register and $1 branch match
        echo "Calling go get to use $OTHER_REPO_BRANCH on fc-retrieval-$1"
        go get github.com/ConsenSys/fc-retrieval-$1@$REG_BRANCH
    fi
}


check_repo gateway ../fc-retrieval-gateway
go mod tidy
