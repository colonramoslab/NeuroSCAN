#!/bin/bash
# heck if we are in the root of the neuroscan repository
    # if not, exit with an error message
if [ ! -d ".git" ]; then
    echo "This script must be run from the root of the neuroscan repository."
    exit 1
fi

echo "Building frontend"
cd frontend || exit 1
yarn

cp ./overwrite/Canvas.js ./node_modules/@metacell/geppetto-meta-ui/3d-canvas/
cp ./overwrite/ThreeDEngine.js ./node_modules/@metacell/geppetto-meta-ui/3d-canvas/threeDEngine/

yarn run build

echo "Frontend build complete"

echo "Starting backend build"

cd ../ || exit 1

go build -o ./neuroscan -ldflags="-w -s" cmd/main.go

echo "Build complete"
