FROM golang:1.18 as build-stage

WORKDIR /tmp/search-accomodation
COPY ./ .

RUN make build || (cd - && exit 1)

FROM debian:buster-slim

COPY --from=build-stage "/tmp/search-accomodation/bin/iploc-server" "/opt/iploc-server"
COPY --from=build-stage "/tmp/search-accomodation/storage/" "/opt/storage/"

CMD cd /opt && /opt/iploc-server