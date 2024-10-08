package http

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"regexp"

	"github.com/evanw/esbuild/pkg/api"
)

type HttpPlugin struct {
	esmImportPattern *regexp.Regexp
}

func NewHttpPlugin() *HttpPlugin {
	return &HttpPlugin{
		esmImportPattern: regexp.MustCompile(`^https?://esm.sh`),
	}
}

func (p *HttpPlugin) onResolve(args api.OnResolveArgs) (api.OnResolveResult, error) {
	base, err := url.Parse(args.Importer)
	if err != nil {
		return api.OnResolveResult{}, err
	}

	relative, err := url.Parse(args.Path)
	if err != nil {
		return api.OnResolveResult{}, err
	}

	resolved := base.ResolveReference(relative)
	return api.OnResolveResult{
		Path:      resolved.String(),
		Namespace: "http-url",
	}, nil
}

func (p *HttpPlugin) onLoad(args api.OnLoadArgs) (api.OnLoadResult, error) {
	fmt.Printf("Fetching %s\n", args.Path)

	resp, err := http.Get(args.Path)
	if err != nil {
		return api.OnLoadResult{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return api.OnLoadResult{}, err
	}

	contents := string(body)

	loader := api.LoaderJS
	if resp.Header.Get("content-type") == "application/json" {
		loader = api.LoaderJSON
	}

	return api.OnLoadResult{
		Contents:   &contents,
		Loader:     loader,
		ResolveDir: path.Dir(args.Path),
	}, nil
}

func Plugin() api.Plugin {
	plugin := NewHttpPlugin()

	return api.Plugin{
		Name: "http-import",
		Setup: func(build api.PluginBuild) {
			build.OnResolve(
				api.OnResolveOptions{Filter: plugin.esmImportPattern.String()},
				plugin.onResolve,
			)
			build.OnLoad(
				api.OnLoadOptions{Filter: ".*", Namespace: "http-url"},
				plugin.onLoad,
			)
		},
	}
}
