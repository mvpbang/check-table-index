#!/bin/bash
set -e

# env
gosrc=$1
out=$2
ver=$(date +%Y%m-%d-%H24)
outrelase=/tmp/gaga
test ! -d /tmp/gaga && mkdir /tmp/gaga

#add files
files=(
  "config.yml"
  "README.md"
  "CHANGELOG.md"
)

# judge args
if [ $# -ne 2 ];then
  echo "Usege: bash $0 xxx.go  xxx"
  exit 9
fi

# window crose compiler
CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -o "${out}_$ver".exe "$gosrc"

# tar release
# shellcheck disable=SC2068

tar -zcvf "${2}_$ver".tar.gz  "${2}_$ver".exe "${files[@]}"
mv "${2}_$ver".tar.gz "${outrelase}"/

#move clean
mv "${2}_$ver".exe  /tmp/