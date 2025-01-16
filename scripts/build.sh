#!/bin/bash
cd /home/inghamemerson/code/NeuroSCAN || exit 1

echo "Pulling latest changes"
git fetch --all
git reset --hard origin/main

echo "Fetch complete"

echo "Building frontend"
cd frontend || exit 1
yarn

cp ./overwrite/Canvas.js ./node_modules/@metacell/geppetto-meta-ui/3d-canvas/
cp ./overwrite/ThreeDEngine.js ./node_modules/@metacell/geppetto-meta-ui/3d-canvas/threeDEngine/

yarn run build

echo "Frontend build complete"

echo "Starting backend build"

cd /home/inghamemerson/code/NeuroSCAN || exit 1

go build -o ./neuroscan -ldflags="-w -s" cmd/neuroscan/main.go

echo "Build complete"

echo "Restarting neuroscan service"
sudo systemctl restart neuroscan

echo "Service restarted"
