## 编译工具链
[NDK 下载](https://developer.android.google.cn/ndk/downloads?hl=zh-cn)
```sh
~/android/android-ndk-r21e/build/tools/make-standalone-toolchain.sh \
--toolchain=aarch64-linux-android-4.9 \
--platform=android-28 \
--install-dir=~/android/toolchain \
--force
```

参数解释：
- toolchain。表示 Android 的 ARCH，arm32 使用 arm-linux-androideabi-4.9，arm64 使用 aarch64-linux-android-4.9。
- platform。表示 Android API 版本：
    - 25：Android 7.1.1
    - 26：Android 8.0
    - 27：Android 8.1
    - 28：Android 9.0
    - 29：Android 10.0
- install-dir。表示 toolchain 安装路径，toolchain 用于交叉编译 Go

## 编译 libgoc.so

```sh
cd go
source build.sh
```

## 编译 hello.cpp

```sh
cd cpp
bash cmake_android.sh
cd build && make
```