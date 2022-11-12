ARG GO_VERSION

FROM golang:$GO_VERSION
WORKDIR /app

RUN set -x \
    && apt-get update \
    && apt-get install -y libgl1-mesa-dev xorg-dev \
    && rm -rf /var/lib/apt/lists/*
