package css

import (
	"os"
	"path/filepath"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
)

func onLoadCss(m *minify.M) func(args api.OnLoadArgs) (api.OnLoadResult, error) {
	return func(args api.OnLoadArgs) (api.OnLoadResult, error) {
		srcFilePath := filepath.Clean(args.Path)

		content, err := os.ReadFile(srcFilePath)
		if err != nil {
			return api.OnLoadResult{}, err
		}

		minified, err := m.String("text/css", string(content))
		if err != nil {
			return api.OnLoadResult{}, err
		}

		return api.OnLoadResult{
			Contents: &minified,
			Loader:   api.LoaderText,
		}, nil
	}
}

func Plugin() api.Plugin {
	m := minify.New()
	m.AddFunc("text/css", html.Minify)

	return api.Plugin{
		Name: "css",
		Setup: func(build api.PluginBuild) {
			build.OnLoad(
				api.OnLoadOptions{Filter: `\.css$`},
				onLoadCss(m),
			)
		},
	}
}
