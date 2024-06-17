#
# Copyright (C) 2019 The Android Open Source Project
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

# This makefile contains the product partition contents for
# a generic phone or tablet device. Only add something here if
# it definitely doesn't belong on other types of devices (if it
# does, use base_product.mk).
$(call inherit-product, $(SRC_TARGET_DIR)/product/media_product.mk)

PRODUCT_PACKAGES_EXCLUDE += Messaging

# /product packages
PRODUCT_PACKAGES += \
    Calendar \
    Camera \
    DeskClock \
    ExactCalculator \
    Gallery2 \
    Music \
    PdfViewer \
    preinstalled-packages-platform-handheld-product.xml \
    QuickSearchBox \
    SettingsIntelligence \
    ThemePicker \
    ThemesStub \
    frameworks-base-overlays \
    Apps \
    SwiftKeyboard \
    ContactsCompose \
    NFTMintApp \
    ClaimNFT \
    WalletApp \
    ethOSMessenger \
    #ethOSAppStore \


PRODUCT_PACKAGES_DEBUG += \
    frameworks-base-overlays-debug


PRODUCT_COPY_FILES += \
    system/media/bootanimation.zip:product/media/bootanimation.zip