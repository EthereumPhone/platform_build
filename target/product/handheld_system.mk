#
# Copyright (C) 2018 The Android Open Source Project
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

# This makefile contains the system partition contents for
# a generic phone or tablet device. Only add something here if
# it definitely doesn't belong on other types of devices (if it
# does, use base_vendor.mk).
$(call inherit-product, $(SRC_TARGET_DIR)/product/media_system.mk)
$(call inherit-product-if-exists, frameworks/base/data/fonts/fonts.mk)
$(call inherit-product-if-exists, external/google-fonts/dancing-script/fonts.mk)
$(call inherit-product-if-exists, external/google-fonts/carrois-gothic-sc/fonts.mk)
$(call inherit-product-if-exists, external/google-fonts/coming-soon/fonts.mk)
$(call inherit-product-if-exists, external/google-fonts/cutive-mono/fonts.mk)
$(call inherit-product-if-exists, external/google-fonts/source-sans-pro/fonts.mk)
$(call inherit-product-if-exists, external/noto-fonts/fonts.mk)
$(call inherit-product-if-exists, external/roboto-fonts/fonts.mk)
$(call inherit-product-if-exists, external/hyphenation-patterns/patterns.mk)
$(call inherit-product-if-exists, frameworks/base/data/keyboards/keyboards.mk)
$(call inherit-product-if-exists, frameworks/webview/chromium/chromium.mk)

PRODUCT_PACKAGES += \
    BasicDreams \
    BlockedNumberProvider \
    BluetoothMidiService \
    BookmarkProvider \
    BuiltInPrintService \
    CalendarProvider \
    cameraserver \
    CameraExtensionsProxy \
    CaptivePortalLogin \
    CertInstaller \
    DocumentsUI \
    DownloadProviderUi \
    EasterEgg \
    ExternalStorageProvider \
    FusedLocation \
    GmsCompat \
    InputDevices \
    KeyChain \
    librs_jni \
    ManagedProvisioning \
    MmsService \
    MtpService \
    MusicFX \
    NfcNci \
    PacProcessor \
    PrintRecommendationService \
    PrintSpooler \
    ProxyHandler \
    screenrecord \
    SecureElement \
    SharedStorageBackup \
    SimAppDialog \
    talkback \
    Telecom \
    TelephonyProvider \
    TeleService \
    Traceur \
    UserDictionaryProvider \
    VpnDialogs \
    vr \
    ethOSBrowser \
    LightNodeStats \
    ethOSWallpaper

PRODUCT_SYSTEM_SERVER_APPS += \
    FusedLocation \
    InputDevices \
    KeyChain \
    Telecom \

PRODUCT_COPY_FILES += \
    frameworks/av/media/libeffects/data/audio_effects.conf:system/etc/audio_effects.conf\
    system/media/bootanimation.zip:system/media/bootanimation.zip\
    system/media/nimbus_verified_proxy:system/bin/nimbus_verified_proxy\
    system/media/helios:system/bin/helios\
#    system/core/libllm/llama:system/bin/llama\
#    system/media/ggml-model-q4_0.bin:system/media/llmmodel/ggml-model-q4_0.bin

PRODUCT_VENDOR_PROPERTIES += \
    ro.carrier?=unknown \
    ro.config.notification_sound?=OnTheHunt.ogg \
    ro.config.alarm_alert?=Alarm_Classic.ogg
