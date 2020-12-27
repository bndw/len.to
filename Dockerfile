###
# Dependency stage
FROM debian:buster-slim as tools

RUN apt-get update \
  && apt-get install -y curl

# Install hugo
ARG HUGO_VERSION=0.79.1

WORKDIR /tmp
RUN curl -L "https://github.com/gohugoio/hugo/releases/download/v${HUGO_VERSION}/hugo_${HUGO_VERSION}_Linux-64bit.tar.gz" \
	| tar xz \
	&& mv ./hugo /bin/hugo \
	&& hugo version

###
# Build stage
FROM tools as build

WORKDIR /build
COPY . .
RUN hugo 

### 
# Execution stage
FROM nginx:stable-alpine
COPY --from=build /build/public /usr/share/nginx/html
