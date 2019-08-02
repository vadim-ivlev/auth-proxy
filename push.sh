#!/bin/bash
git add -A .
git commit -m "."

git push gitlab primitive
git push github primitive
git push origin primitive

git push gitlab master
git push github master
git push origin master


