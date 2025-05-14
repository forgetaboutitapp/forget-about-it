#/bin/sh
set -e
pwd
go run ./generate.go
cd frontend
dart run build_runner build -d
flutter build web --release
mkdir ../server/web -p
cp ./build/web ../server/web -rf
cd ../
cd server

sqlc generate
