FROM golang:1.18.3-alpine as builder

COPY ./ /template-renderer
RUN cd /template-renderer \
 && go get \
 && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go test template-renderer/cmd -v \
 && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o templater main.go

FROM scratch

COPY --from=builder /template-renderer/templater /templater

ENTRYPOINT ["/templater"]