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
go test ./...
mkdir /home/runner/work/forget-about-it/forget-about-it/server/web -p
cp /home/runner/work/forget-about-it/forget-about-it/frontend/build/web ../server -rf