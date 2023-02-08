FROM golang:1.20
WORKDIR /app

RUN apt-get update
RUN apt-get install -y --no-install-recommends \
      libasound2-dev \
      libc6-dev \
      libgl1-mesa-dev \
      libglu1-mesa-dev \
      libxcursor-dev \
      libxi-dev \
      libxinerama-dev \
      libxrandr-dev \
      libxxf86vm-dev \
      pkg-config

RUN git config --global --add safe.directory "$PWD"
