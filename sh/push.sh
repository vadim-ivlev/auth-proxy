#!/bin/bash
git add -A .
git commit -m "."


git push gitlab --all
git push github --all
git push origin --all


