// Package static holds the files served to the browser as-is, such as
// JavaScript, CSS, and images.
package static

import "embed"

// Files holds every file under js/, embedded into the binary.
//
//go:embed js
var Files embed.FS
