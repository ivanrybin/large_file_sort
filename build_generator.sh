#!/bin/bash
cd cmd/file_gen || exit
go build -o file_gen
chmod +x file_gen
mv file_gen ../../file_gen