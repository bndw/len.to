###
# Dependency stage
FROM debian:buster-slim as tools

RUN apt-get update \
  && apt-get install -y curl

# Install hugo
WORKDIR /tmp
RUN curl -SL https://github.com/gohugoio/hugo/releases/download/v0.55.6/hugo_0.55.6_Linux-64bit.tar.gz \
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
