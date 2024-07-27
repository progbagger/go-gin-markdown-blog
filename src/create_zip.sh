#!/bin/bash

mkdir blog

go build -o blog/myblog-binary cmd/server/main.go
cp -r static templates blog/
cp admin_credentials.txt blog/

zip -r blog.zip blog
rm -rf blog
