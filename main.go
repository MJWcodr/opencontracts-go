package opencontractinggo

import (
	"encoding/json"
)

type OpenContract struct {
	URI string `json:"uri"`
	Version string `json:"version"`
	Extensions []string `json:"extensions"`
	PublishedDate string `json:"publishedDate"`
}

// Load Function
func (oc *OpenContract) Load(data []byte) error {
	// Implementation to load data into the struct
	return json.Unmarshal(data, oc)
}

// Print Function
func (oc *OpenContract) Print() string {
	return string(oc.URI + " - " + oc.Version + " - " + oc.PublishedDate)
  :q
}


