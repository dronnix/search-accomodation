FROM golang:1.18 as build-stage

WORKDIR /tmp/search-accomodation
COPY ./ .

RUN make build || (cd - && exit 1)

FROM debian:buster-slim

COPY --from=build-stage "/tmp/search-accomodation/bin/iploc-data-importer" "/opt/iploc-data-importer"
COPY --from=build-stage "/tmp/search-accomodation/storage/" "/opt/storage/"
COPY --from=build-stage "/tmp/search-accomodation/data/" "/opt/data/"

CMD cd /opt && /opt/iploc-data-importer