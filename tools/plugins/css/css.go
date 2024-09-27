package css

import (
	"os"
	"path/filepath"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
)

type CssAssetPlugin struct {
	minifier *minify.M
}

func NewCssAssetPlugin() *CssAssetPlugin {
	m := minify.New()
	m.AddFunc("text/css", html.Minify)

	return &CssAssetPlugin{
		minifier: m,
	}
}

func (p *CssAssetPlugin) onLoadCss(args api.OnLoadArgs) (api.OnLoadResult, error) {
	srcFilePath := filepath.Clean(args.Path)

	content, err := os.ReadFile(srcFilePath)
	if err != nil {
		return api.OnLoadResult{}, err
	}

	minified, err := p.minifier.String("text/css", string(content))
	if err != nil {
		return api.OnLoadResult{}, err
	}

	return api.OnLoadResult{
		Contents: &minified,
		Loader:   api.LoaderCSS,
	}, nil
}

func Plugin() api.Plugin {
	plugin := NewCssAssetPlugin()

	return api.Plugin{
		Name: "css",
		Setup: func(build api.PluginBuild) {
			build.OnLoad(
				api.OnLoadOptions{Filter: `\.css$`},
				plugin.onLoadCss,
			)
		},
	}
}
