FROM node:13-alpine AS node
COPY web/ /web
WORKDIR /web
RUN npm install && npm run build

FROM golang:1.13-alpine
COPY firmware/ /build
COPY --from=node /web/dist/ /web/
RUN apk add make upx git
WORKDIR /build
RUN mkdir /dist && make clean all
