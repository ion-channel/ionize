FROM scratch

ARG GIT_COMMIT_HASH

LABEL org.metadata.vcs-url="https://github.com/ion-channel/ionize" \
      org.metadata.vcs-commit-id=$GIT_COMMIT_HASH \
      org.metadata.name="Ionize" \
      org.metadata.description="Ion Channel Analysis Management Tool"

COPY --from=alpine /etc/ssl /etc/ssl
ADD ionize /ionize

WORKDIR /data

CMD ["/ionize", "analyze"]
