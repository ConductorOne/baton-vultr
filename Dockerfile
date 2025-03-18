FROM gcr.io/distroless/static-debian11:nonroot
ENTRYPOINT ["/baton-vultr"]
COPY baton-vultr /