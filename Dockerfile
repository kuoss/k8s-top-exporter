FROM golang:1.25-alpine AS build
WORKDIR /temp/
COPY . ./
RUN go mod download
RUN go build -o /k8s-top-exporter ./cmd/k8s-top-exporter

FROM gcr.io/distroless/static:nonroot
LABEL org.opencontainers.image.source https://github.com/jmnote/k8s-top-exporter
WORKDIR /
COPY --from=build /k8s-top-exporter /k8s-top-exporter

EXPOSE     9977
USER       nonroot
ENTRYPOINT ["/k8s-top-exporter"]
