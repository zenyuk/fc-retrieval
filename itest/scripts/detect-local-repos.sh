#!/bin/bash
echo "**********************************************************************************"
echo "*** Check that go.mod does not point to a local client, gateway, or admin repo ***"
echo "**********************************************************************************"

REPLACE_EXISTS=`grep "replace " ../go.mod`
if [ -n "$REPLACE_EXISTS" ]; then
    echo "Using local client, gateway, and or admin repos"
    echo "Replace set for:"
    echo "$REPLACE_EXISTS"
    exit 1
fi
echo "Using remote repos"
