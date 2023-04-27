package website

import "embed"

var (
	//go:embed *.html
	Templates embed.FS

	//go:embed css/* js/* images/*
	Static embed.FS
)

const (
	HomeTemplate = "index.html"
)
