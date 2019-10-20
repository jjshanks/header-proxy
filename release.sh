# adopted from https://medium.com/travis-on-docker/how-to-version-your-docker-images-1d5c577ebf54

if [[ -n $(git status -s) ]]
then
  echo "must commit changes before running"
  exit
fi

set -ex

docker run --rm -v "$PWD":/app treeder/bump patch
version=`cat VERSION`
echo "version: $version"

git add VERSION
git commit -m "version $version"
git tag -a "$version" -m "version $version"
git push
git push --tags

docker build -t jjshanks/header-proxy:latest .
docker tag jjshanks/header-proxy:latest jjshanks/header-proxy:$version
docker push jjshanks/header-proxy:$version
