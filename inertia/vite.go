package inertia

import (
	"fmt"
	"maps"
	"os"
	"regexp"
	"strings"

	"github.com/tesseris-go/tesseris/utils"
)

type Record map[string]string

type Resolver func(string, string, Record, Record) Record

type Vite struct {
	prod           bool
	nonce          *string
	buildDirectory string
	publicPath     string
	hotFile        string
	integrityKey   string

	styleTagAttributesResolvers  []Resolver
	scriptTagAttributesResolvers []Resolver
}

func NewVite(prod bool) *Vite {
	return &Vite{
		prod:           prod,
		nonce:          nil,
		buildDirectory: "build",
		publicPath:     "./public",
		hotFile:        "hot",
		integrityKey:   "integrity",

		styleTagAttributesResolvers:  []Resolver{},
		scriptTagAttributesResolvers: []Resolver{},
	}
}

func isCssPath(path string) bool {
	matched, _ := regexp.MatchString(`\.(css|less|sass|scss|styl|stylus|pcss|postcss)$`, path)
	return matched
}

func (v *Vite) parseAttributes(attrs Record) []string {
	var res []string

	for key, value := range attrs {
		if value != "" {
			res = append(res, fmt.Sprintf(`%s="%s"`, key, value))
		}
	}

	return res
}

func (v *Vite) makeStylesheetTagWithAttributes(url string, attributes Record) string {
	attrs := Record{
		"rel":  "stylesheet",
		"href": url,
		"nonce": func() string {
			if v.nonce != nil {
				return *v.nonce
			}
			return ""
		}(),
	}

	// Copy additional attributes to attrs
	maps.Copy(attrs, attributes)

	parsedAttrs := v.parseAttributes(attrs)

	return "<link " + strings.Join(parsedAttrs, " ") + " />"
}

func (v *Vite) resolveStylesheetTagAttributes(src, url string, chunk, manifest Record) Record {
	attributes := Record{}

	if v.integrityKey != "" {
		integrity, exists := chunk[v.integrityKey]

		if exists {
			attributes["integrity"] = integrity
		}
	}

	for _, resolver := range v.styleTagAttributesResolvers {
		resolverAttributes := resolver(src, url, chunk, manifest)
		for key, value := range resolverAttributes {
			attributes[key] = value
		}
	}

	return attributes
}

func (v *Vite) makeScriptTagWithAttributes(url string, attributes Record) string {
	attrs := Record{
		"type": "module",
		"src":  url,
		"nonce": func() string {
			if v.nonce != nil {
				return *v.nonce
			}
			return ""
		}(),
	}

	// Copy additional attributes to attrs
	maps.Copy(attrs, attributes)

	parsedAttrs := v.parseAttributes(attrs)

	return "<script " + strings.Join(parsedAttrs, " ") + "></script>"
}

func (v *Vite) resolveScriptTagAttributes(src, url string, chunk, manifest Record) Record {
	attributes := Record{}

	if v.integrityKey != "" {
		integrity, exists := chunk[v.integrityKey]

		if exists {
			attributes["integrity"] = integrity
		}
	}

	for _, resolver := range v.scriptTagAttributesResolvers {
		resolverAttributes := resolver(src, url, chunk, manifest)
		for key, value := range resolverAttributes {
			attributes[key] = value
		}
	}

	return attributes
}

func (v *Vite) makeTagForChunk(src, url string, chunk, manifest Record) string {
	if isCssPath(url) {
		return v.makeStylesheetTagWithAttributes(
			url,
			v.resolveStylesheetTagAttributes(src, url, chunk, manifest),
		)
	}

	return v.makeScriptTagWithAttributes(
		url,
		v.resolveScriptTagAttributes(src, url, chunk, manifest),
	)
}

func (v *Vite) ViteTags(entrypoints []string) string {
	eps := append([]string{"@vite/client"}, entrypoints...)

	if !v.prod {
		tags := utils.Map(eps, func(entrypoint string) string {
			return v.makeTagForChunk(entrypoint, v.hotAsset(entrypoint), nil, nil)
		})

		return strings.Join(tags, "\n")
	}

	// TODO: Implement production tags or maybe just keep using generated bootstrap ssr

	// Steps:
	// 1. Read manifest.json
	// 2. Read entrypoints from manifest.json
	// 3. Generate tags for each entrypoint
	// 4. Return tags

	// manifestPath := v.buildDirectory + "/manifest.json"

	return ""
}

func (v *Vite) hotFilePath() string {
	return v.publicPath + "/" + "hot"
}

func (v *Vite) hotAsset(asset string) string {
	fileContents, err := os.ReadFile(v.hotFilePath())

	if err != nil {
		fmt.Println(err)
		return ""
	}

	return strings.Trim(string(fileContents), " ") + "/" + asset
}

func (v *Vite) ReactRefresh() string {
	if v.prod {
		return ""
	}

	return fmt.Sprintf(`<script type="module">
			import RefreshRuntime from '%s'
			RefreshRuntime.injectIntoGlobalHook(window)
			window.$RefreshReg$ = () => {}
			window.$RefreshSig$ = () => (type) => type
			window.__vite_plugin_react_preamble_installed__ = true
		</script>`,
		v.hotAsset("@react-refresh"),
	)
}
