// Package socrates is a simple HTML-template helper which allows templates to be embedded in others
// using commented definitions.
package socrates

import (
	"bufio"
	"fmt"
	"html/template"
	"io"
	"os"
	"strings"
)

type paths []string

func (ps *paths) load(path string) error {
	*ps = append(*ps, path)

	f, err := os.Open(path)
	if err != nil {
		return err
	}

	br := bufio.NewReader(f)
	firstLine, err := br.ReadString('\n')
	f.Close()
	if err != nil {
		return fmt.Errorf("error reading first line of file %q: %v", path, err)
	}

	prefix, suffix := "<!-- USE ", " -->\n"
	if strings.HasPrefix(firstLine, prefix) && strings.HasSuffix(firstLine, suffix) {
		return ps.load(strings.TrimSuffix(strings.TrimPrefix(firstLine, prefix), suffix))
	}
	return nil
}

// Template contains the path to the template that is to be rendered.  This template can "extend" another
// by referencing the parent template name <!-- USE parent_template.html --> on the first line (must be the
// only text on the first line).  Any blocks defined in parent_template.html can then be overriden by the
// current template.
type Template string

// Execute loads all the files necessary to render the template tpl and then renders it to w.
func (tpl Template) Execute(w io.Writer, x interface{}) error {
	var ps paths
	if err := ps.load(string(tpl)); err != nil {
		return err
	}

	psRev := make([]string, len(ps))
	for i, l := range ps {
		psRev[len(ps)-1-i] = l
	}

	t, err := template.ParseFiles(psRev...)
	if err != nil {
		return err
	}
	return t.Execute(w, x)
}
