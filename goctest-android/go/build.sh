CGO_ENABLED=1 GOOS=android NDK_TOOLCHAIN=~/android/toolchain \
GOARCH=arm64 CC=~/android/toolchain/bin/aarch64-linux-android-gcc \
go build -buildmode=c-shared -o libgoc.so hello.go