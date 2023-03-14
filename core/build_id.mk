#
# Copyright (C) 2008 The Android Open Source Project
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
# BUILD_ID is usually used to specify the branch name
# (like "MAIN") or a branch name and a release candidate
# (like "CRB01").  It must be a single word, and is
# capitalized by convention.

ifneq (,$(filter flame coral,$(TARGET_PRODUCT)))
    BUILD_ID=TP1A.221005.002.B2
else ifneq (,$(filter oriole raven bluejay,$(TARGET_PRODUCT)))
    BUILD_ID=T2B3.230109.009
else ifneq (,$(filter panther,$(TARGET_PRODUCT)))
    BUILD_ID=TQ2A.230305.008
else
    BUILD_ID=TQ2A.230305.008.C1
endif
