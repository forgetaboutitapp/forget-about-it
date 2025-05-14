#/bin/sh
set -e
pwd
cd frontend
pwd
dart run build_runner build -d
flutter build web --release
mkdir ../server/web -p
cp ./build/web ../server -rf
cd ../
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
cd server

sqlc generate
for os in windows; do
for arch in amd64; do
mkdir -p  /tmp/tmp-release/${os}-${arch}
ext=""
if [ "$os" == "windows" ]; then
ext=".exe"
fi
GOOS=$os GOARCH=$arch go build --ldflags '-extldflags "-static"' -o /tmp/tmp-release/forget-about-it-$os-$arch/provision${ext} ./cmd/provision
GOOS=$os GOARCH=$arch go build --ldflags '-extldflags "-static"' -o /tmp/tmp-release/forget-about-it-$os-$arch/server${ext} ./cmd/server
done
done