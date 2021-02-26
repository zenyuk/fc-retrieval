#!/bin/bash
set -e
echo "*******************************************************************************"
echo "*** Update go.mod and go.sum to point to the latest packages on github.com  ***"
echo "*** for client, gateway, gateway-admin, and provider-admin.                 ***"
echo "*******************************************************************************"

cd ..

# Remove any local references
sed '/replace .*/d' go.mod > go.mod.temp1
rm go.mod
mv go.mod.temp1 go.mod

ITEST_DIR="../fc-retrieval-itest"
ITEST_BRANCH=`git rev-parse --abbrev-ref HEAD`
echo "itest repo branch: $ITEST_BRANCH"


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
    BRANCH_EXISTS=`git branch -r --list origin/$ITEST_BRANCH`

    if [ $ITEST_BRANCH != $OTHER_REPO_BRANCH ]; then 
        echo "itest and $1 branch do not match"
        if [ -n "$BRANCH_EXISTS" ]; then
            echo "ERROR: Branch $ITEST_BRANCH exists on the $1 repo, but the $1 repo is currently using branch $OTHER_REPO_BRANCH"
            exit 1
        fi

        echo "Calling go get to use main on $1"
        cd $ITEST_DIR
        go get -u -t github.com/ConsenSys/fc-retrieval-$1@main
    else
        echo itest and $1 branch match
        GITHASH=`git rev-parse $ITEST_BRANCH`
        echo "Calling go get to use $OTHER_REPO_BRANCH on fc-retrieval-$1 ($GITHASH)"
        cd $ITEST_DIR
        go get -u -t -v github.com/ConsenSys/fc-retrieval-$1@$GITHASH
    fi
}


check_repo client ../fc-retrieval-client
check_repo common ../fc-retrieval-common
check_repo gateway-admin ../fc-retrieval-gateway-admin
check_repo provider-admin ../fc-retrieval-provider-admin
check_repo provider ../fc-retrieval-provider
check_repo register ../fc-retrieval-register
go mod tidy