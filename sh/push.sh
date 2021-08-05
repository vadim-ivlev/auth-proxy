#!/bin/bash

echo "Введите комментарий к коммиту:"
read msg
echo "Комментарий = $msg."
echo 


git add -A .
git commit -m "$msg."

git push gitlab --all #master
git push github --all #master
git push origin --all #master


git push gitlab --tags 
git push github --tags 
git push origin --tags 
