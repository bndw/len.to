# Build stage
FROM debian:buster-slim as build
ARG HUGO_VERSION=0.79.1

RUN apt-get update \
  && apt-get install -y curl

RUN cd /usr/local/bin \
  && curl -L "https://github.com/gohugoio/hugo/releases/download/v${HUGO_VERSION}/hugo_${HUGO_VERSION}_Linux-64bit.tar.gz" \
  | tar xz \
  && hugo version

WORKDIR /build
COPY . .
RUN hugo 
RUN find ./public -type f -name "*.html" > /tmp/html

# Validation stage
FROM validator/validator:latest as validator
COPY --from=build /bin/cat /bin/cat
COPY --from=build /build/public /build/public
COPY --from=build /tmp/html /tmp/html
RUN vnu-runtime-image/bin/vnu --verbose --errors-only $(cat /tmp/html)

# Execution stage
FROM nginx:stable-alpine
COPY --from=build /build/public /usr/share/nginx/html
