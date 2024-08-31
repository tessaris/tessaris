#!/bin/sh

CURRENT_BRANCH="$(git rev-parse --abbrev-ref HEAD)"
if [ $CURRENT_BRANCH != "main" ]; then
  echo "Error: You are not on the main branch. A new release can only be make from main.";
  exit 1;
fi

# Update the local repository just in case
git pull

VERSION=`cat .version`
echo Adding git tag with version v${VERSION};
git tag v${VERSION};
git push origin v${VERSION};
