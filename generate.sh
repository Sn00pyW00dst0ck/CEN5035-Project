#!/bin/bash

# Function to check if a command exists
check_command() {
    command -v "$1" >/dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo "$1 is not installed!"
        case $1 in
            "go")
                echo "Please install Go, https://go.dev/doc/install."
                ;;
            "oapi-codegen")
                echo "Run 'go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest' to install oapi-codegen."
                ;;
            "npm")
                echo "Please install NodeJS and npm, https://docs.npmjs.com/downloading-and-installing-node-js-and-npm."
                ;;
            "npx")
                echo "Run 'npm install -g npx' to install npx."
                ;;
            *)
                echo "Unknown command $1. Please install it manually."
                ;;
        esac
        exit 1
    else
        echo "$1 is already installed."
    fi
}

echo "Checking For Required Tools..."
echo -n "   "
check_command go
echo -n "   "
check_command oapi-codegen
echo -n "   "
check_command npm
echo -n "   "
check_command npx

echo "Generating Backend Server..."
oapi-codegen --config=models/config.yaml models/schema.yaml

echo "Generating Frontend Types..."
npx openapi-typescript models/schema.yaml -o frontend/src/types/api.ts