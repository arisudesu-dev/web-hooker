#!/bin/sh
echo "Status: 200 OK"
echo "Content-type: text/plain"
echo ""
echo "Hello from /bin/sh"
echo ""
echo "Script is:" $0
echo ""
echo "Working dir is:" $(pwd)
echo ""
echo "Environment is:"
env | sort
echo ""
