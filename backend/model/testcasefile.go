package model

type TestCaseFile struct {
	Name     string `json:"name"`
	Bytes    uint64 `json:"bytes"`
	ModTime  string `json:"mod_time"`
	FileType string `json:"file_type"`
}
