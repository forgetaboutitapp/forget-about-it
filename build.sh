#/bin/sh
set -e
pwd
dart pub global activate fastforge 
sudo apt install locate
maindir=$(pwd)
go run ./generate.go
cd frontend
dart run build_runner build -d
cd ..
pwd
ls -lha
go run ./generate.go
cd ./frontend
flutter build web --release --wasm
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
mkdir "$maindir"/server/web -p
cp "$maindir"/frontend/build/web ../server -rf
cd "$maindir"/server
for os in linux windows freebsd openbsd darwin; do
  for arch in arm64 amd64; do
    mkdir -p  "$1"/tmp-release/${os}-${arch}
    ext=""
    if [ "$os" == "windows" ]; then
      ext=".exe"
    fi
    echo building $os $arch provision
    GOOS=$os GOARCH=$arch go build --ldflags '-extldflags "-static"' -o "$1"/tmp-release/forget-about-it-$os-$arch/provision${ext} ./cmd/provision
  done
done
cd "$maindir"/server

for os in linux windows freebsd openbsd darwin; do
  for arch in arm64 amd64; do
    ext=""
    if [ "$os" == "windows" ]; then
      ext=".exe"
    fi
    echo building $os $arch server
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
mv "$maindir"/frontend/build/app/outputs/flutter-apk/app-release.apk "$1"/releases/forget-about-it.apk
