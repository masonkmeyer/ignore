package repository

import "io"

// Provider is the contract for searching repos
type Provider interface {
	// Get returns a reader for the matching search term
	Get(term string) (io.Reader, error)

	// Search gets a list of available gitignore files
	Search(term string) ([]string, error)
}
