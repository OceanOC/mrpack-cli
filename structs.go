package main

type mrpcli struct {
	modPacks   map[string]ModPack
	files      []string
	automated  bool
	nodownload bool
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
	Env       Environment `json:"env"`
}

type Environment struct {
	Client string `json:"client"`
	Server string `json:"server"`
}
