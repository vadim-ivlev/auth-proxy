#!/bin/bash

branch=$(git branch --show-current)

echo "Deploying current pushed branch '$branch' to dev"
echo "Are you sure?"
echo "Press Ctrl-C to cancel"
read -p "Press any key to continue "


git push origin $branch:dev

echo 'Push button 'Deploy LK and Deploy to Develop' at https://git.rgwork.ru/masterback/auth-proxy/-/pipelines'


