---
title: Validating Generated HTML
date: 2021-01-01T12:12:10-07:00
draft: false

---

I ran this website through an HTML validator and
discovered it was missing some closing `</div>` tags. Tracking these errors down 
can be challenging with generated HTML, so I decided to add linting to the build process and catch problems early.

##

The [Nu HTML Checker](https://validator.w3.org/nu/) I used above conveniently ships a 
[Docker image](https://hub.docker.com/r/validator/validator/) containing a CLI called `vnu`. 
I ran into two small challenges integrating `vnu` in the multi-stage build:

##
## 1. `vnu` wants a list of files
According to the [usage](https://validator.github.io/validator/#usage) vnu wants a list of files:
```
vnu-runtime-image/bin/vnu OPTIONS FILES
```

However Hugo produces a `public` directory containing a tree of HTML. 
To work around this I used `find` to create a file containing a list of generated HTML files:
```
RUN find /build/public -type f -name "*.html" > /tmp/htmlfiles.txt
```
Now we can just `cat /tmp/htmlfiles.txt` and pass the output to `vnu` 

##
## 2. `validator/validator` is based on a "distroless" image
The [validator/validator](https://hub.docker.com/r/validator/validator/) image is based on `gcr.io/distroless/base` and doesn't have `cat`, so instead I grabbed it from the previous build stage:
```
COPY --from=build /bin/cat /bin/cat
```

Finally, we can run `vnu` with a list of all generated HTML:
```
RUN vnu-runtime-image/bin/vnu --verbose --errors-only $(cat /tmp/htmlfiles.txt)
```

All said and done, here's what the Dockerfile's validator stage looks like:
```
FROM validator/validator:latest as validator
COPY --from=build /bin/cat /bin/cat
COPY --from=build /build/public /build/public
COPY --from=build /tmp/htmlfiles.txt /tmp/htmlfiles.txt
RUN vnu-runtime-image/bin/vnu --verbose --errors-only $(cat /tmp/htmlfiles.txt)
```

##
##
With validation now in place for every generated HTML file, I discovered and fixed 8 existing errors.
Here's [the commit](https://github.com/bndw/len.to/commit/e5a5187c22567668587e8dbc7819f58dc7e03b62) implementing HTML validation and the associated bug fixes.
##
