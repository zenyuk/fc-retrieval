#!/bin/bash
set -e

echo "********************************************************"
echo "*** Update go.mod to point to local repos: gateway   ***"
echo "********************************************************"

REPLACE_TEXT_G="replace github.com/ConsenSys/fc-retrieval-provider => ../fc-retrieval-provider"

sed '/replace .*/d' ../go.mod > ../go.mod.temp
rm ../go.mod
mv ../go.mod.temp ../go.mod

printf "%s\n" "$REPLACE_TEXT_G" >> ../go.mod

echo go.mod file is now:
cat ../go.mod
