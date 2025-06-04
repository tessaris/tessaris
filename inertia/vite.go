package inertia

import (
	"encoding/json"
	"fmt"
	"maps"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/tessaris/tessaris/utils"
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

func (v *Vite) getChunkFromManifest(manifest Manifest, entrypoint string) string {
	chunk, exists := manifest[entrypoint]

	if !exists {
		return ""
	}

	tags := ""

	if chunk.Css != nil {
		tags = strings.Join(utils.Map(chunk.Css, func(css string) string {
			assetUrl := path.Join(v.buildDirectory, css)

			return v.makeTagForChunk(entrypoint, assetUrl, nil, nil)
		}), "\n") + "\n"
	}

	assetUrl := path.Join(v.buildDirectory, chunk.File)

	return tags + v.makeTagForChunk(entrypoint, assetUrl, nil, nil)
}

func (v *Vite) ViteTags(entrypoints []string) string {
	if !v.prod {
		eps := append([]string{"@vite/client"}, entrypoints...)

		tags := utils.Map(eps, func(entrypoint string) string {
			return v.makeTagForChunk(entrypoint, v.hotAsset(entrypoint), nil, nil)
		})

		return strings.Join(tags, "\n")
	}

	manifestPath := path.Join(v.publicPath, v.buildDirectory, "manifest.json")

	manifestFile, err := os.ReadFile(manifestPath)

	if err != nil {
		fmt.Println(err)

		return ""
	}

	if manifestFile == nil {
		return ""
	}

	var manifest Manifest

	err = json.Unmarshal(manifestFile, &manifest)

	if err != nil {
		fmt.Println(err)

		return ""
	}

	tags := utils.Map(entrypoints, func(entrypoint string) string {
		return v.getChunkFromManifest(manifest, entrypoint)
	})

	return strings.Join(tags, "\n")
}

func (v *Vite) hotFilePath() string {
	return path.Join(v.publicPath, "hot")
}

func (v *Vite) hotAsset(asset string) string {
	fileContents, err := os.ReadFile(v.hotFilePath())

	if err != nil {
		fmt.Println(err)
		return ""
	}

	assetPath, err := url.JoinPath(strings.Trim(string(fileContents), " "), asset)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	return assetPath
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
