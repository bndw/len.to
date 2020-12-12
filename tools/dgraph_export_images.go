package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"time"
)

func main() {
	// Build dgraph RDF
	fmt.Println("{")
	fmt.Println("set {")
	walk(images, imgRDF)
	walk(images, imgTagRDF)
	fmt.Println("}")
	fmt.Println("}")
}

// imgRDF prints an image as rdf triples
func imgRDF(img *Post) {
	node := nodeName(img.Title)

	// img node
	rdfObjString(node, "dgraph.type", "Image")
	rdfObjString(node, "name", img.Title)
	rdfObjString(node, "url", img.ImgURL)
	rdfObjString(node, "location", img.Location)
	rdfObjString(node, "date", img.Date.Format(time.RFC3339))

	// edges to tags
	for _, tag := range img.Tags {
		rdfObjNode(node, "tagged", nodeName(tag))
	}
}

// imgTagRDF prints an image's tags as rdf triples
func imgTagRDF(img *Post) {
	for _, tag := range img.Tags {
		node := nodeName(tag)
		rdfObjString(node, "dgraph.type", "Tag")
		rdfObjString(node, "name", tag)
	}
}

func nodeName(s string) string {
	h := sha1.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func rdfObjString(s, p, o string) {
	fmt.Printf("_:%s <%s> \"%s\" .\n", s, p, o)
}
func rdfObjNode(s, p, o string) {
	fmt.Printf("_:%s <%s> _:%s .\n", s, p, o)
}
