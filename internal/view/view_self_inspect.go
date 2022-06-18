package view

// language=gohtml
const SelfInspect = `
	Usage:
	    {{ print "$" | colorize "gray" }} {{ print "go-arch-lint" | colorize "blue" }} self-inspect {{ print "--json" | colorize "magenta" }}
	{{ " " }}
	Note:
	    this command created for integration with dev-tools and IDE's
	    and not have any ascii output support for inspection results.
	    please use this command with flag {{ print "--json" | colorize "magenta" }}

	{{- /*gotype: github.com/fe3dback/go-arch-lint/internal/models.SelfInspect*/ -}}
`
