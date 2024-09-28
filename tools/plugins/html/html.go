package htmlasset

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/izquiratops/dobunezumi/tools/utils/directory"
	"github.com/izquiratops/dobunezumi/tools/utils/hash"

	"github.com/tdewolff/minify/v2"
	tdewolff "github.com/tdewolff/minify/v2/html"

	"golang.org/x/net/html"
)

type HtmlAssetPlugin struct {
	distLocalPath string
	minifier      *minify.M
}

func NewHtmlAssetPlugin(distLocalPath string) *HtmlAssetPlugin {
	m := minify.New()
	m.AddFunc("text/html", tdewolff.Minify)

	return &HtmlAssetPlugin{
		distLocalPath: distLocalPath,
		minifier:      m,
	}
}

func (p *HtmlAssetPlugin) findRelativeDistPath(htmlAbsPath, assetLocalPath string) (string, error) {
	htmlAbsDir := filepath.Dir(htmlAbsPath)
	sourceAbsPath, err := filepath.Abs(filepath.Join(htmlAbsDir, assetLocalPath))
	if err != nil {
		return "", fmt.Errorf("error getting absolute path for source: %v", err)
	}

	distAbsPath, err := filepath.Abs(filepath.Join(p.distLocalPath, filepath.Base(assetLocalPath)))
	if err != nil {
		return "", fmt.Errorf("error getting absolute path for destination: %v", err)
	}

	if _, err := os.Stat(p.distLocalPath); os.IsNotExist(err) {
		message := fmt.Sprintf("directory %s does not exist.\n", p.distLocalPath)
		log.Fatal(message)
	}

	if err := directory.MoveTo(sourceAbsPath, distAbsPath); err != nil {
		return "", err
	}

	// Setting the new src attribute as ./whatever.ext because everything is moved at the same path level
	return filepath.Base(assetLocalPath), nil
}

func (p *HtmlAssetPlugin) processHtmlNode(n *html.Node, htmlAbsPath string) error {
	if n.Type == html.ElementNode && n.Data == "img" {
		for i, attr := range n.Attr {
			if attr.Key == "src" && !filepath.IsAbs(attr.Val) {
				var err error
				n.Attr[i].Val, err = p.findRelativeDistPath(htmlAbsPath, attr.Val)
				if err != nil {
					return err
				}
				break
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if err := p.processHtmlNode(c, htmlAbsPath); err != nil {
			return err
		}
	}

	return nil
}

func (p *HtmlAssetPlugin) handleHtmlContents(htmlBytes []byte, htmlAbsPath string) (*html.Node, error) {
	doc, err := html.Parse(bytes.NewReader(htmlBytes))
	if err != nil {
		return nil, fmt.Errorf("error parsing HTML: %w", err)
	}

	if err := p.processHtmlNode(doc, htmlAbsPath); err != nil {
		return nil, err
	}

	return doc, nil
}

func (p *HtmlAssetPlugin) minifyHtml(doc *html.Node) ([]byte, error) {
	var buf bytes.Buffer
	if err := html.Render(&buf, doc); err != nil {
		return nil, fmt.Errorf("error rendering HTML: %w", err)
	}

	minified, err := p.minifier.String("text/html", buf.String())
	if err != nil {
		return nil, fmt.Errorf("error minifying HTML: %w", err)
	}

	return []byte(minified), nil
}

func (p *HtmlAssetPlugin) onResolve(args api.OnResolveArgs) (api.OnResolveResult, error) {
	assetPath := filepath.Join(args.ResolveDir, args.Path)

	return api.OnResolveResult{
		Path:      assetPath,
		Namespace: "html-asset",
	}, nil
}

func (p *HtmlAssetPlugin) onLoad(args api.OnLoadArgs) (api.OnLoadResult, error) {
	htmlBytes, err := os.ReadFile(args.Path)
	if err != nil {
		return api.OnLoadResult{}, fmt.Errorf("error reading HTML file: %w", err)
	}

	htmlNode, err := p.handleHtmlContents(htmlBytes, args.Path)
	if err != nil {
		return api.OnLoadResult{}, err
	}

	minifiedContents, err := p.minifyHtml(htmlNode)
	if err != nil {
		return api.OnLoadResult{}, err
	}

	outputFileName := hash.GenerateHash(args.Path, minifiedContents)

	// Write the output
	distFilePath := filepath.Join(p.distLocalPath, outputFileName)
	if err := os.WriteFile(distFilePath, minifiedContents, 0644); err != nil {
		return api.OnLoadResult{}, err
	}

	// Create a JavaScript module to import the output file
	outputModule := fmt.Sprintf("export default %q;", outputFileName)

	return api.OnLoadResult{
		Contents: &outputModule,
		Loader:   api.LoaderJS,
	}, nil
}

func Plugin(distLocalPath string) api.Plugin {
	plugin := NewHtmlAssetPlugin(distLocalPath)

	return api.Plugin{
		Name: "html-asset",
		Setup: func(build api.PluginBuild) {
			build.OnResolve(
				api.OnResolveOptions{Filter: `\.html$`},
				plugin.onResolve,
			)
			build.OnLoad(
				api.OnLoadOptions{Filter: `\.html$`, Namespace: "html-asset"},
				plugin.onLoad,
			)
		},
	}
}
