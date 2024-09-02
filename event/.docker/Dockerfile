ARG DISTRO_VERSION=latest
FROM ubuntu:${DISTRO_VERSION} AS devtools

CMD mkdir cmd
ADD bin/app /cmd/app
ADD bin/migrate /cmd/migrate
ADD config.yml /config.yml
ADD db/migrations /migrations

CMD ["/cmd/app"]