# Android Setup

## Experimental Support

Oak experimentally supports android. As of right now, it should be possible with some workarounds to build programs with Oak that run as expected on Android, with some limitations; and these programs will likely need to have some custom settings apart from a build for other platforms. We'll be looking to improve the amount of work needed to adapt Oak to Android going forward.

## Steps

How can you build a program with Oak for android? This unfortunately requires a lot of setup:

1. Download Android Studio
1. Via Android Studio's SDK manager (Setup a Project -> Tools -> SDK Manager) (we won't be using this project):
1. Specify and note the Android SDK location at the top of this window. If on windows, ensure this path does not contain spaces.
1. Download SDK Tools -> NDK (Side by side), Platform-Tools, and Build-Tools
1. Setup your environment (e.g. ~/.bash_profile) to include the following, substituting ANDROID_HOME with the noted android SDK location from before, and noting that the version strings are not likely the same as they will be on your machine:

    ```bash
    export ANDROID_HOME=/d/android
    export ANDROID_NDK_HOME=$ANDROID_HOME/ndk/23.1.7779620/
    export PATH=$PATH:$ANDROID_HOME/platform-tools
    export PATH=$PATH:$ANDROID_HOME/build-tools/32.0.0
    ```

1. Install [gomobile](https://pkg.go.dev/golang.org/x/mobile/cmd/gomobile)
1. (`gomobile init` may require a C compiler)
1. Run `gomobile build --target=android/arm64` from that package (other architectures may work as well, but we lack devices to test them on.) A C compiler is no longer needed at this stage.
1. You may optionally provide an `AndroidManifest.xml` as described in the gomobile docs. (I have not tested this.)
1. From here you can use `adb` to connect to a test device and install the built apk. In my testing, the only commands needed where `adb usb`, `adb install` and `adb logcat` for debugging.
