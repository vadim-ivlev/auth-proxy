#!/bin/bash

CurBuild=$1
CurHost=$2
Success=0

# 120 секунд
for attempt in {1..6}; do
    printf "Attempt: $attempt\n"
    Build=$(curl --insecure $CurHost)
    quotes="\""
    # удаляем кавычки
    Build=${Build//$quotes/}
    # проверяем - если полученная версия с сервера совпадает с текущая версией сборки то все хорошо 
    if [ "$Build" != "" ] && [ "$Build" = "$CurBuild" ]; then
        echo "Success! Assembly versions are the same: $Build"
        Success=1
        break
    fi
    echo "Build: $Build, Current build: $CurBuild"
    sleep 20
done

if [ "$Success" = "0" ]; then
    echo "Versions do not match! Ending with an error:  $Build"
    exit 1
fi