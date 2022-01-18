FROM golang:1.17.6-alpine as builder

COPY ./ /deploy-now-configuration-template-action
RUN cd /deploy-now-configuration-template-action \
 && go get \
 && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o template-action main.go

FROM scratch

COPY --from=builder /deploy-now-configuration-template-action/template-action /template-action

ENTRYPOINT ["/template-action"]