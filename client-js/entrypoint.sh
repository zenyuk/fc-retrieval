#!/usr/bin/env sh
set -exuo pipefail

env

echo 'client-js ready'
npm run test-e2e
