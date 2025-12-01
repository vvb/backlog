#!/bin/bash

# Build the application
echo "Building backlog application..."
go build -o bl

if [ $? -ne 0 ]; then
    echo "Build failed!"
    exit 1
fi

echo "Build successful!"

# Ask user if they want to install globally
read -p "Do you want to install 'bl' globally to /usr/local/bin? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]
then
    sudo mv bl /usr/local/bin/
    echo "✓ Installed 'bl' to /usr/local/bin/"
    echo "You can now use 'bl' from anywhere!"
else
    echo "Binary 'bl' is available in the current directory"
    echo "You can move it manually or add it to your PATH"
fi

echo ""
echo "Run 'bl --help' to get started!"

