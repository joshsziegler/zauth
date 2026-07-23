// Package zauth embeds static assets so they ship inside the binary.
package zauth

import "embed"

// Templates holds the Go/HTML templates under templates/.
//
//go:embed all:templates
var Templates embed.FS

// Public holds static web assets (CSS, fonts) under public/.
//
//go:embed all:public
var Public embed.FS
