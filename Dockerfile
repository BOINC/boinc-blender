# This file is part of BOINC.
# https://boinc.berkeley.edu
# Copyright (C) 2025 University of California
#
# BOINC is free software; you can redistribute it and/or modify it
# under the terms of the GNU Lesser General Public License
# as published by the Free Software Foundation,
# either version 3 of the License, or (at your option) any later version.
#
# BOINC is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
# See the GNU Lesser General Public License for more details.
#
# You should have received a copy of the GNU Lesser General Public License
# along with BOINC.  If not, see <http://www.gnu.org/licenses/>.

FROM debian
WORKDIR /app
ENV ARGS ""

RUN apt update && apt install -y \
    curl \
    xz-utils \
    unzip \
    zip \
    libx11-6 \
    libxrender1 \
    libxxf86vm-dev \
    libxfixes3 \
    libxi6 \
    libxkbcommon0 \
    libsm6 \
    libgl1 \
    libegl1

RUN curl -L https://ftp.halifax.rwth-aachen.de/blender/release/Blender4.4/blender-4.4.3-linux-x64.tar.xz -o blender.tar.xz
RUN tar -xJvf blender.tar.xz --strip-components=1 -C /bin
RUN rm blender.tar.xz

CMD unzip input.zip && ./boincblender ${ARGS} && zip output.zip *.png
