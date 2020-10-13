package multimc

type Pack struct {
	Components    []Component `json:"components"`
	FormatVersion int         `json:"formatVersion"`
}

type Component struct {
	Important      bool   `json:"important,omitempty"`
	UID            string `json:"uid"`
	Version        string `json:"version"`
	CachedVolatile bool   `json:"cachedVolatile,omitempty"`
	DependencyOnly bool   `json:"dependencyOnly,omitempty"`
}
