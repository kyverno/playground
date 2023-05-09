FROM node:20-alpine as frontend

WORKDIR /var/www

RUN apk add --no-cache --virtual .gyp python3 make g++ \
    && npm set progress=false \
    && npm config set depth 0

COPY frontend .

RUN npm install \
    && npm run build

FROM golang:1.20-alpine as builder

ARG LD_FLAGS="-s -w"
ARG TARGETPLATFORM

WORKDIR /app

COPY backend .

COPY --from=frontend /var/www/dist /app/dist

RUN export GOOS="$(echo ${TARGETPLATFORM} | cut -d / -f1)" && \
    export GOARCH="$(echo ${TARGETPLATFORM} | cut -d / -f2)" && \
    apk --no-cache add ca-certificates && \
    update-ca-certificates

RUN go env

RUN go get -d -v \
    && go install -v

RUN go build -ldflags="${LD_FLAGS}" -o /app/build/playground -v

FROM scratch

WORKDIR /app

ENV GIN_MODE=release

USER 1234

COPY --from=builder /app/build/playground /app/playground
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

EXPOSE 8080

ENTRYPOINT ["/app/playground", "-host", "0.0.0.0"]
