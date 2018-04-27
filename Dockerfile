FROM golang:1.10-alpine AS builder

WORKDIR /go/src/github.com/ninedraft/shibatest
COPY . .
RUN CGO_ENABLED=0 go install -v -ldflags='-w -s -extldflags="-static"'
RUN apk --no-cache add tzdata zip
WORKDIR /usr/share/zoneinfo
RUN zip -r -0 /zoneinfo.zip .

FROM scratch
COPY --from=builder /go/bin/shibatest /shibatest
COPY --from=builder /zoneinfo.zip /
ENV ZONEINFO /zoneinfo.zip
EXPOSE 8080
ENTRYPOINT ["/shibatest"]