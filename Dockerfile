FROM debian:buster-slim as build

RUN apt-get update && apt-get install -y curl

# Install hugo
WORKDIR /tmp
RUN curl -SL https://github.com/gohugoio/hugo/releases/download/v0.55.6/hugo_0.55.6_Linux-64bit.tar.gz \
	| tar xz \
	&& mv hugo /usr/bin \
	&& hugo version

# Build the site
WORKDIR /build
COPY . .
RUN hugo 

### 
# execution stage
FROM nginx:stable-alpine

COPY --from=build /build/public /usr/share/nginx/html
