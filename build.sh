#!/bin/bash
rm -Rf dist/
mkdir dist/
go build
mv SimpleStaticServer dist/
