#!/bin/bash
git add -A .
echo "commit message = $1."
git commit -m "$1."

git push gitlab --all #master
git push github --all #master
git push origin --all #master


