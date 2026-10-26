# syntax=docker/dockerfile:1
# Grounded's Xbox-gate stub API (groundedd). Plain HTTP on 23580 behind the TLS
# terminator; /healthz on 23581.
FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS build
ARG TARGETOS TARGETARCH
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH GOWORK=off go build -trimpath -buildvcs=false -ldflags="-s -w" -o /out/groundedd ./cmd/groundedd

FROM gcr.io/distroless/static
COPY --from=build /out/groundedd /groundedd
USER nonroot:nonroot
EXPOSE 23580 23581
ENTRYPOINT ["/groundedd"]
