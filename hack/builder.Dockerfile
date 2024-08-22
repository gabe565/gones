FROM golang:1.23.0
WORKDIR /app

RUN set -x \
    && apt-get update \
    && apt-get install -y --no-install-recommends \
      libasound2-dev \
      libc6-dev \
      libgl1-mesa-dev \
      libglu1-mesa-dev \
      libxcursor-dev \
      libxi-dev \
      libxinerama-dev \
      libxrandr-dev \
      libxxf86vm-dev \
      pkg-config \
  && git config --global --add safe.directory "$PWD" \
  && rm -rf /var/lib/apt/lists/*
