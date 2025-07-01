#!/bin/bash

IMAGE_NAME="smr"

tags=$(docker images --format "{{.Repository}}:{{.Tag}} {{.CreatedAt}}" | grep "^$IMAGE_NAME:" | sort -rk2)

tag_list=($(echo "$tags" | awk '{print $1}'))

keep_tags=("${tag_list[@]:0:2}")

delete_tags=("${tag_list[@]:2}")

for tag in "${delete_tags[@]}"; do
    echo "Deleting tag: $tag"
    docker rmi "$tag"
done

echo "Kept tags:"
for tag in "${keep_tags[@]}"; do
    echo "  $tag"
done