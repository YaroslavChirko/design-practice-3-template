FROM golang:1.14 as build

RUN apt-get update && apt-get install -y ninja-build

RUN go get -u github.com/YaroslavChirko/design-practice-2-template/build/cmd/bood

WORKDIR /go/src/practice-3
COPY . .

RUN CGO_ENABLED=0 bood

# ==== Final image ====
FROM alpine:3.11
WORKDIR /opt/practice-3
COPY entry.sh ./
COPY --from=build /go/src/practice-3/out/bin/* ./
ENTRYPOINT ["/opt/practice-3/entry.sh"]
CMD ["server"]
