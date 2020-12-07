#!/bin/bash
git add -A .
git commit -m "."


git push gitlab --all #master
git push github --all #master
git push origin --all #master


