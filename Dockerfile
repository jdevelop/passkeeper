FROM node:13-alpine AS node
COPY web/ /web
WORKDIR /web
RUN npm install && npm run build

FROM golang:1.13-alpine AS builder
COPY firmware/ /build
COPY --from=node /web/dist/ /web/
RUN apk add make upx git
WORKDIR /build
RUN mkdir /dist && make clean all

#FROM jdevelop/passkeeper:buildroot-2018.08.2-rpi-zero-w as buildroot
FROM jdevelop/passkeeper:buildroot-2018.08.2-rpi-zero as buildroot
COPY --from=builder /dist/ /build/board/rootfs_overlay/root/
COPY buildroot/.config /build/.config
COPY buildroot/linux-config /build/linux-config
WORKDIR /build
RUN make O=/build PASSKEEPER=/build FORCE_UNSAFE_CONFIGURE=1 -C /buildroot/buildroot-2018.08.2 linux-rebuild
RUN make O=/build PASSKEEPER=/build FORCE_UNSAFE_CONFIGURE=1 -C /buildroot/buildroot-2018.08.2

FROM alpine:3.10
COPY --from=buildroot /build/images/sdcard.img /dist/
WORKDIR /dist
