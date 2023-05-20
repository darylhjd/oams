package frontend

import "embed"

//go:embed build/web/*
var Website embed.FS
