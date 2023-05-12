NDK_ROOT="~/android/android-ndk-r21e"
ARCH="arm64-v8a"
# 28 对应 Android 9.0
API=28

cmake -B build \
	-DCMAKE_TOOLCHAIN_FILE=$NDK_ROOT/build/cmake/android.toolchain.cmake -Wno-dev \
	-DANDROID_NDK=$NDK_ROOT -DCMAKE_SYSTEM_NAME=Android -DANDROID_PLATFORM=android-${API} \
	-DANDROID_ABI=${ARCH} -DAndroid=ON -DANDROID_STL=c++_static \
	-DANDROID_NATIVE_API_LEVEL=${API}  -DCMAKE_CXX_STANDARD=17