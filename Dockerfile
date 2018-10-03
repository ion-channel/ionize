FROM golang:1.10.3-alpine3.8

WORKDIR /go/src/github.com/ion-channel/ionize/

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .


FROM scratch

ARG APP_NAME
ARG BUILD_DATE
ARG VERSION
ARG GIT_COMMIT_HASH
ARG ENVIRONMENT

LABEL org.metadata.build-date=$BUILD_DATE \
      org.metadata.version=$VERSION \
      org.metadata.vcs-url="https://github.com/ion-channel/ionize" \
      org.metadata.vcs-commit-id=$GIT_COMMIT_HASH \
      org.metadata.name="Ionize" \
      org.metadata.description="Ion Channel Analysis Management Tool"

WORKDIR /root/

COPY --from=0 /etc/ssl /etc/ssl
COPY --from=0 /go/src/github.com/ion-channel/ionize/ionize .

ENTRYPOINT ["./ionize"]
