#!/bin/bash
# Confirm that we are in the neuroscan directory and that the neuroscan binary exists
if [ ! -d ".git" ]; then
    echo "This script must be run from the root of the neuroscan repository."
    exit 1
fi

if [ ! -f "./neuroscan" ]; then
    echo "The neuroscan binary does not exist. Please run the build script first."
    exit 1
fi

echo "Cleaning up stale video files..."
./neuroscan cleanup
if [ $? -ne 0 ]; then
    echo "Cleanup failed. Please check the neuroscan logs for more details."
    exit 1
fi

echo "Cleanup complete. Stale video files have been removed."
