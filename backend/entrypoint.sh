#!/bin/sh
# This script ensures all dependencies are tidy before starting the server.

# Exit immediately if a command exits with a non-zero status.
set -e

# Run go mod tidy to ensure go.sum is correct.
echo "Running go mod tidy..."
go mod tidy

# Now, execute the main command passed to the container (air).
echo "Starting air hot-reloader..."
exec "$@"
