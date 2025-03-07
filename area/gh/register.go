package main

import "github.com/alexfalkowski/infraops/internal/gh"

var repositories []*gh.Repository

// RegisterRepository registers a repository to be created.
func RegisterRepository(repository *gh.Repository) {
	repositories = append(repositories, repository)
}
