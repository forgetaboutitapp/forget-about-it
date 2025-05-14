#/bin/sh
set -e
pwd
go run ./generate.go
cd frontend
dart run build_runner build -d
cd ..
pwd
ls -lha
go run ./generate.go
cd ./frontend
flutter build web --release
flutter build apk --release
echo "finished building"
mkdir ../server/web -p
echo "finished mkdir"
cp ./build/web ../server/web -rf
echo "finished copying"
cd ../server
echo "finished cd"
pwd
ls
sqlc generate
echo "finished sqlc"
pwd
#go test ./...
mkdir /home/runner/work/forget-about-it/forget-about-it/server/web -p
cp /home/runner/work/forget-about-it/forget-about-it/frontend/build/web ../server -rf
cd /home/runner/work/forget-about-it/forget-about-it/server
for os in linux windows freebsd openbsd darwin; do
  for arch in arm64 amd64; do
    mkdir -p  "$1"/tmp-release/${os}-${arch}
    ext=""
    if [ "$os" == "windows" ]; then
      ext=".exe"
    fi
    GOOS=$os GOARCH=$arch go build --ldflags '-extldflags "-static"' -o "$1"/tmp-release/forget-about-it-$os-$arch/provision${ext} ./cmd/provision
  done
done
cd /home/runner/work/forget-about-it/forget-about-it/server

for os in linux windows freebsd openbsd darwin; do
  for arch in arm64 amd64; do
    ext=""
    if [ "$os" == "windows" ]; then
      ext=".exe"
    fi
    GOOS=$os GOARCH=$arch go build --ldflags '-extldflags "-static"' -o "$1"/tmp-release/forget-about-it-$os-$arch/server${ext} ./cmd/server
  done
done
mkdir -p "$1"/releases/
cd "$1"/tmp-release/
for os in linux windows freebsd openbsd darwin; do
  for arch in arm64 amd64; do
    if [ "$os" == "windows" ]; then
        zip -r "$1"/releases/forget-about-it-${os}-${arch}.zip forget-about-it-$os-$arch
    else
      tar -czvf "$1"/releases/forget-about-it-${os}-${arch}.tar.gz forget-about-it-$os-$arch
    fi
  done
done
cd /home/runner/work/forget-about-it/forget-about-it/frontend
mkdir "$1"/releases
ls build
ls build/app
ls build/app/outputs
ls build/app/outputs/flutter-apk/
mv build/app/outputs/flutter-apk/app-release.apk "$1"/releases/forget-about-it.apk