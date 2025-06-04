package inertia

type ManifestEntry struct {
	File           string   `json:"file"`
	Name           string   `json:"name"`
	Src            string   `json:"src"`
	IsEntry        bool     `json:"isEntry,omitempty"`
	DynamicImports []string `json:"dynamicImports,omitempty"`
	Imports        []string `json:"imports,omitempty"`
	Css            []string `json:"css,omitempty"`
}

type Manifest map[string]ManifestEntry
