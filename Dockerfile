FROM golang:1.15 as build
WORKDIR /app
COPY . /app
RUN CGO_ENABLED=0 go build -o semver main.go

FROM alpine as final
COPY --from=build /app/semver /
WORKDIR /workdir
CMD /semver /workdir