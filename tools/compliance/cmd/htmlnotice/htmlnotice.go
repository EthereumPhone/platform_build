// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"html"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"android/soong/tools/compliance"

	"github.com/google/blueprint/deptools"
)

var (
	outputFile  = flag.String("o", "-", "Where to write the NOTICE text file. (default stdout)")
	depsFile    = flag.String("d", "", "Where to write the deps file")
	includeTOC  = flag.Bool("toc", true, "Whether to include a table of contents.")
	stripPrefix = flag.String("strip_prefix", "", "Prefix to remove from paths. i.e. path to root")
	title       = flag.String("title", "", "The title of the notice file.")

	failNoneRequested = fmt.Errorf("\nNo license metadata files requested")
	failNoLicenses    = fmt.Errorf("No licenses found")
)

type context struct {
	stdout      io.Writer
	stderr      io.Writer
	rootFS      fs.FS
	includeTOC  bool
	stripPrefix string
	title       string
	deps        *[]string
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: %s {options} file.meta_lic {file.meta_lic...}

Outputs an html NOTICE.html file.

Options:
`, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	// Must specify at least one root target.
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(2)
	}

	if len(*outputFile) == 0 {
		flag.Usage()
		fmt.Fprintf(os.Stderr, "must specify file for -o; use - for stdout\n")
		os.Exit(2)
	} else {
		dir, err := filepath.Abs(filepath.Dir(*outputFile))
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot determine path to %q: %s\n", *outputFile, err)
			os.Exit(1)
		}
		fi, err := os.Stat(dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot read directory %q of %q: %s\n", dir, *outputFile, err)
			os.Exit(1)
		}
		if !fi.IsDir() {
			fmt.Fprintf(os.Stderr, "parent %q of %q is not a directory\n", dir, *outputFile)
			os.Exit(1)
		}
	}

	var ofile io.Writer
	ofile = os.Stdout
	if *outputFile != "-" {
		ofile = &bytes.Buffer{}
	}

	var deps []string

	ctx := &context{ofile, os.Stderr, os.DirFS("."), *includeTOC, *stripPrefix, *title, &deps}

	err := htmlNotice(ctx, flag.Args()...)
	if err != nil {
		if err == failNoneRequested {
			flag.Usage()
		}
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
	if *outputFile != "-" {
		err := os.WriteFile(*outputFile, ofile.(*bytes.Buffer).Bytes(), 0666)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not write output to %q: %s\n", *outputFile, err)
			os.Exit(1)
		}
	}
	if *depsFile != "" {
		err := deptools.WriteDepFile(*depsFile, *outputFile, deps)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not write deps to %q: %s\n", *depsFile, err)
			os.Exit(1)
		}
	}
	os.Exit(0)
}

// htmlNotice implements the htmlnotice utility.
func htmlNotice(ctx *context, files ...string) error {
	// Must be at least one root file.
	if len(files) < 1 {
		return failNoneRequested
	}

	// Read the license graph from the license metadata files (*.meta_lic).
	licenseGraph, err := compliance.ReadLicenseGraph(ctx.rootFS, ctx.stderr, files)
	if err != nil {
		return fmt.Errorf("Unable to read license metadata file(s) %q: %v\n", files, err)
	}
	if licenseGraph == nil {
		return failNoLicenses
	}

	// rs contains all notice resolutions.
	rs := compliance.ResolveNotices(licenseGraph)

	ni, err := compliance.IndexLicenseTexts(ctx.rootFS, licenseGraph, rs)
	if err != nil {
		return fmt.Errorf("Unable to read license text file(s) for %q: %v\n", files, err)
	}

	fmt.Fprintln(ctx.stdout, "<!DOCTYPE html>")
	fmt.Fprintln(ctx.stdout, "<html><head>")
	fmt.Fprintln(ctx.stdout, "<style type=\"text/css\">")
	fmt.Fprintln(ctx.stdout, "body { padding: 2px; margin: 0; }")
	fmt.Fprintln(ctx.stdout, "ul { list-style-type: none; margin: 0; padding: 0; }")
	fmt.Fprintln(ctx.stdout, "li { padding-left: 1em; }")
	fmt.Fprintln(ctx.stdout, ".file-list { margin-left: 1em; }")
	fmt.Fprintln(ctx.stdout, "</style>")
	if 0 < len(ctx.title) {
		fmt.Fprintf(ctx.stdout, "<title>%s</title>\n", html.EscapeString(ctx.title))
	}
	fmt.Fprintln(ctx.stdout, "</head>")
	fmt.Fprintln(ctx.stdout, "<body>")

	if 0 < len(ctx.title) {
		fmt.Fprintf(ctx.stdout, "  <h1>%s</h1>\n", html.EscapeString(ctx.title))
	}
	ids := make(map[string]string)
	if ctx.includeTOC {
		fmt.Fprintln(ctx.stdout, "  <ul class=\"toc\">")
		i := 0
		for installPath := range ni.InstallPaths() {
			id := fmt.Sprintf("id%d", i)
			i++
			ids[installPath] = id
			var p string
			if 0 < len(ctx.stripPrefix) && strings.HasPrefix(installPath, ctx.stripPrefix) {
				p = installPath[len(ctx.stripPrefix):]
				if 0 == len(p) {
					if 0 < len(ctx.title) {
						p = ctx.title
					} else {
						p = "root"
					}
				}
			} else {
				p = installPath
			}
			fmt.Fprintf(ctx.stdout, "    <li id=\"%s\"><strong>%s</strong>\n      <ul>\n", id, html.EscapeString(p))
			for _, h := range ni.InstallHashes(installPath) {
				libs := ni.InstallHashLibs(installPath, h)
				fmt.Fprintf(ctx.stdout, "        <li><a href=\"#%s\">%s</a>\n", h.String(), html.EscapeString(strings.Join(libs, ", ")))
			}
			fmt.Fprintln(ctx.stdout, "      </ul>")
		}
		fmt.Fprintln(ctx.stdout, "  </ul><!-- toc -->")
	}
	for h := range ni.Hashes() {
		fmt.Fprintln(ctx.stdout, "  <hr>")
		for _, libName := range ni.HashLibs(h) {
			fmt.Fprintf(ctx.stdout, "  <strong>%s</strong> used by:\n    <ul class=\"file-list\">\n", html.EscapeString(libName))
			for _, installPath := range ni.HashLibInstalls(h, libName) {
				if id, ok := ids[installPath]; ok {
					if 0 < len(ctx.stripPrefix) && strings.HasPrefix(installPath, ctx.stripPrefix) {
						fmt.Fprintf(ctx.stdout, "      <li><a href=\"#%s\">%s</a>\n", id, html.EscapeString(installPath[len(ctx.stripPrefix):]))
					} else {
						fmt.Fprintf(ctx.stdout, "      <li><a href=\"#%s\">%s</a>\n", id, html.EscapeString(installPath))
					}
				} else {
					if 0 < len(ctx.stripPrefix) && strings.HasPrefix(installPath, ctx.stripPrefix) {
						fmt.Fprintf(ctx.stdout, "      <li>%s\n", html.EscapeString(installPath[len(ctx.stripPrefix):]))
					} else {
						fmt.Fprintf(ctx.stdout, "      <li>%s\n", html.EscapeString(installPath))
					}
				}
			}
			fmt.Fprintf(ctx.stdout, "    </ul>\n")
		}
		fmt.Fprintf(ctx.stdout, "  </ul>\n  <a id=\"%s\"/><pre class=\"license-text\">", h.String())
		fmt.Fprintln(ctx.stdout, html.EscapeString(string(ni.HashText(h))))
		fmt.Fprintln(ctx.stdout, "  </pre><!-- license-text -->")
	}
	fmt.Fprintln(ctx.stdout, "</body></html>")

	*ctx.deps = ni.InputNoticeFiles()

	return nil
}
