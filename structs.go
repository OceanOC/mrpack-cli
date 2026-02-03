package main

type mrpcli struct {
	files      []string
	automated  bool
	nodownload bool
	nolauncher bool
	outputDir  string
	modpackDir string
}

type ModPack struct {
	Name          string       `json:"name"`
	Game          string       `json:"game"`
	VersionID     string       `json:"versionId"`
	FormatVersion int          `json:"formatVersion"`
	Dependencies  Dependencies `json:"dependencies"`
	Files         []Files      `json:"files"`
}

type Dependencies struct {
	Fabric    string `json:"fabric-loader"`
	Quilt     string `json:"quilt-loader"`
	Minecraft string `json:"minecraft"`
	NeoForge  string `json:"neoforge"`
	Forge     string `json:"forge"`
}

type Files struct {
	Path      string      `json:"path"`
	Downloads []string    `json:"downloads"`
	Hashes    []string    `json:"hashes"`
	FileSize  int         `json:"fileSize"`
	Env       Environment `json:"env"`
}

type Environment struct {
	Client string `json:"client"`
	Server string `json:"server"`
}

type MineLauncher struct {
	Profiles map[string]Profile `json:"profiles"`
}

type Profile struct {
	Created       string `json:"created"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	Icon          string `json:"icon"`
	LastUsed      string `json:"lastUsed"`
	GameDirectory string `json:"gameDir"`
	LastVersionID string `json:"lastVersionId"`
}
