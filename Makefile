include bin/build/make/go.mak
include bin/build/make/git.mak

# Diagrams generated from https://github.com/loov/goda.
diagrams: gh-diagram

gh-diagram:
	$(MAKE) package=. create-diagram
