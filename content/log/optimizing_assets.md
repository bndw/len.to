---
title: Optimizing Asset Delivery
date: 2020-12-26T20:47:24-07:00
draft: false

---

I've wanted to improve the CSS and JavaScript portions of this website for a 
while. The main items I want to address are:
- Removing JavaScript from HTML templates
- Find a better way to inject site content into JavaScript
- Bundle all CSS into a single file
- Minify and fingerprint CSS and JavaScript

I initially intended to plug various CSS and JavaScript processing tools into the 
Docker build pipeline, but that quickly spiraled out of control when `npm` got involved.
Then I discovered [Hugo Pipes](https://gohugo.io/hugo-pipes/introduction/), which 
solves these problems natively.


##
## Removing JavaScript from HTML templates

I started by extracting all JavaScript from HTML and appending it to a single `./assets/js/main.js` file. 
Then from HTML I load, process, and render the script:
```
{{ $main := resources.Get "js/main.js" | js.Build | minify | fingerprint }}
<script type="text/javascript" src="{{ $main.RelPermalink }}" defer></script>
```

However, this blew up because the current JavaScript implementation depends on Go templating to build an array of tags:
```
let tags = [{{ range $name, $_ := $.Site.Taxonomies.tags }}{{ $name }},{{ end }}]
```

And processing this with [js.Build](https://gohugo.io/hugo-pipes/js/) fails because it's not valid JavaScript. 

##
## A better way to inject site content into JavaScript

Turns out Hugo Pipes supports passing a map of params to js.Build, and those params can be accessed in JavaScript with the
`import` keyword. So all together we get:
```
{{ $tags := slice }}
{{ range $name, $_ := $.Site.Taxonomies.tags }}
  {{ $tags = $tags | append $name }} 
{{ end }}

{{ $main := resources.Get "js/main.js" | js.Build (dict "params" (dict "tags" $tags)) | minify | fingerprint }}
<script type="text/javascript" src="{{ $main.RelPermalink }}" defer></script>
```

And in JavaScript we just import the tags array:
```
import {tags} from '@params';
```

##
## Bundling CSS

The site uses (2) CSS files:
- Skeleton CSS for a basic foundation, served from a CDN
- Local styles in `main.css`

I moved both stylesheets into `./assets/css/` and then bundled, minified, and fingerprinted
with Hugo Pipes:

```
{{ $skeleton := resources.Get "css/skeleton.css" }}
{{ $style := resources.Get "css/main.css" }}
{{ $css := slice $skeleton $style | resources.Concat "css/bundle.css" | minify | fingerprint }}
<link rel="stylesheet" href="{{ $css.Permalink }}" integrity="{{ $css.Data.Integrity }}">
```


##
## Bonus: Removing Google Fonts

I also took this opportunity to replace the Google-hosted OpenSans font with 
the system font stack, using the same technique as github.com:
```
body {
  font-family: -apple-system,BlinkMacSystemFont,Segoe UI,Helvetica,Arial,sans-serif,Apple Color Emoji,Segoe UI Emoji;
}
```

##
These changes, implemented in 
[0b4e300](https://github.com/bndw/len.to/commit/0b4e3002b03bb49ad08efc69e941f99e6a7c0da2) and 
[0c29a6c](https://github.com/bndw/len.to/commit/0c29a6c13b4f2c20943d31fd073e0c721d17b002) 
produced faster page load times, improved user privacy, and generally better code quality.
##
