package main

import (
	"math/rand"
	"text/template"
	"time"
)

func newTemplate() *template.Template {
	const (
		tmpl = `<!DOCTYPE html>
<head>
  <style>
		{{ .Style }}
  </style>
</head>
<body>
	<a href="/">
		<img src="{{ index .Images 1 }}">
		<img class="two" src="{{ index .Images 2 }}">
		<img class="three" src="{{ index .Images 3 }}">
	</a>
</body>
`
	)

	return template.Must(template.New("index").Parse(tmpl))
}

func randomStyle() string {
	const (
		s1 = `
    img {
      display: inline-block;
      position: fixed;
      width: 100%;
    }
    .two {
      left: 20%; z-index: 1;
    }
    .three {
      left: 40%;
      z-index: 2;
    }
  `
		s2 = `
    img {
      display: inline-block;
      position: fixed;
      width: 100%;
    }
    .two {
      left: 20%;
      z-index: 1;
    }
    .three {
      left: 92%;
      z-index: 2;
    }
	`
		s3 = `
    img {
      display: inline-block;
      position: fixed;
      width: 100%;
    }
    .two {
      left: 80%;
      z-index: 1;
    }
    .three {
      left: 90%;
      z-index: 2;
    }
	`
	)
	var styles = []string{s1, s2, s3}
	r := rand.New(rand.NewSource(time.Now().Unix()))
	return styles[r.Intn(len(styles))]
}
