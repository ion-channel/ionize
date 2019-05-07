FROM golang:1.11-alpine

ARG BUILD_PATH

ADD . $BUILD_PATH

WORKDIR $BUILD_PATH

RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build -o /ionize .


FROM scratch

ARG GIT_COMMIT_HASH

LABEL org.metadata.vcs-url="https://github.com/ion-channel/ionize" \
      org.metadata.vcs-commit-id=$GIT_COMMIT_HASH \
      org.metadata.name="Ionize" \
      org.metadata.description="Ion Channel Analysis Management Tool"

COPY --from=0 /etc/ssl /etc/ssl
COPY --from=0 /ionize /ionize

WORKDIR /data

CMD ["/ionize", "analyze"]
