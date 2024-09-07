#!/bin/bash

# Memeriksa apakah argumen direktori proyek diberikan
if [ -z "$1" ]; then
  echo "Usage: $0 <project_directory>"
  exit 1
fi

project_directory="$1"

# Menemukan dan memformat semua file SQL dalam proyek
echo ''
echo 'Running format all sql file .....'
echo ''
echo '10%'
echo '20%'
echo '30%'
echo '40%'
echo '50%'
find "$project_directory" -type f -name '*.sql' -print0 | xargs -0 -I {} pg_format --s 2 --g -o {} {}
echo '60%'
echo '70%'
echo '80%'
echo '90%'
echo '100%'

echo ''
echo 'Running format all file .....'
echo ''
echo '10%'
echo '20%'
echo '30%'
echo '40%'
echo '50%'
gofmt -w .
echo '60%'
echo '70%'
echo '80%'
echo '90%'
echo '100%'
echo ''
echo 'Format berhasil .....'
echo ''
