#!/bin/bash

set -ex

DGRAPH=localhost:8080
IMG_RDF=dgraph_imgs.rdf
SCHEMA_RDF=dgraph_schema.rdf

go run dgraph_schema.go > $SCHEMA_RDF
go run dgraph_export_images.go > $IMG_RDF

# load schema
curl -H "Content-Type: application/rdf" "$DGRAPH/alter?commitNow=true" -XPOST --data-binary @$SCHEMA_RDF

# load images
curl -H "Content-Type: application/rdf" "$DGRAPH/mutate?commitNow=true" -XPOST --data-binary @$IMG_RDF
