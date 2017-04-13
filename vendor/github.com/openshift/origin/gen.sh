#!/usr/bin/bash

make
oc new-build -D $'FROM scratch\nUSER 1001' --dry-run -o yaml &> scratch.yaml
oc new-build -D $'FROM scratch\nUSER 1001' --to asdf --dry-run -o yaml &> asdf.yaml
oc new-build -D $'FROM centos:7\nUSER 1001' --dry-run -o yaml &> centos.yaml
