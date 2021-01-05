---
title: Emoji Favicons 🤔
date: 2020-12-30
draft: false

---


Turns out [SVG](https://en.wikipedia.org/wiki/Scalable_Vector_Graphics) can draw images from text and is a valid favicon format [supported by most major browsers](https://en.wikipedia.org/wiki/Favicon#File_format_support). 


##
## Here's how it works

Create a `favicon.svg` file. There are a lot of opinions on what size(s) to make your favicon,
but modern browsers support a 32x32 icon so we'll use that. Just set the `font-size` to 32 and offset the `y` axis by the same amount, so the icon is visible:
```
<svg xmlns="http://www.w3.org/2000/svg">
  <text font-size="32" y="32">🖼</text>
</svg>
```

Then use it in your `index.html`:
```
<head>
  <link rel="icon" href="favicon.svg">
</head>
```

Here's [the commit](https://github.com/bndw/len.to/commit/9c534a6369c72a03c12158a3334570201196649f) implementing the emoji favicon on this site.
