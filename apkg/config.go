package apkg

// Config represents the ASUSTOR package (apkg) configuration.
type Config struct {
	General  GeneralConfig   `json:"general"`
	Desktop  *DesktopConfig  `json:"desktop,omitempty"`
	Register *RegisterConfig `json:"register,omitempty"`
}

// GeneralConfig represents the general package configuration.
type GeneralConfig struct {
	Package      string   `json:"package,omitempty"`
	Name         string   `json:"name,omitempty"`
	Version      string   `json:"version,omitempty"`
	Depends      []string `json:"depends,omitempty"`
	Conflicts    []string `json:"conflicts,omitempty"`
	Developer    string   `json:"developer,omitempty"`
	Maintainer   string   `json:"maintainer,omitempty"`
	Email        string   `json:"email,omitempty"`
	Website      string   `json:"website,omitempty"`
	Architecture string   `json:"architecture,omitempty"`
	Firmware     string   `json:"firmware,omitempty"`
	Model        string   `json:"model,omitempty"`
	DefaultLang  string   `json:"default-lang,omitempty"`
	MemoryLimit  int      `json:"memory-limit,omitempty"`
	MemoryAdvice int      `json:"memory-advice,omitempty"`
}

// DesktopConfig represents the desktop configuration.
type DesktopConfig struct {
	App *struct {
		Type      string `json:"type,omitempty"`
		SessionID string `json:"session-id,omitempty"`
		Protocol  string `json:"protocol,omitempty"`
		Port      int    `json:"port,omitempty"`
		URL       string `json:"url,omitempty"`
	} `json:"app,omitempty"`
	Privilege *struct {
		Accessible   string `json:"accessible,omitempty"`
		Customizable bool   `json:"customizable,omitempty"`
	} `json:"privilege,omitempty"`
}

// RegisterConfig represents the register configuration.
type RegisterConfig struct {
	SymbolicLink map[string][]string `json:"symbolic-link,omitempty"`
	ShareFolder  []struct {
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`
	} `json:"share-folder,omitempty"`
	Port         []int `json:"port,omitempty"`
	BootPriority *struct {
		StartOrder int `json:"start-order,omitempty"`
		StopOrder  int `json:"stop-order,omitempty"`
	} `json:"boot-priority,omitempty"`
	Prerequisites *struct {
		EnableService  []string `json:"enable-service,omitempty"`
		RestartService []string `json:"restart-service,omitempty"`
	} `json:"prerequisites,omitempty"`
}
