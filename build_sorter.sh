#!/bin/bash
cd cmd/file_sort || exit
go build -o file_sort
chmod +x file_sort
mv file_sort ../../file_sort