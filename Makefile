include bin/build/make/help.mak
include bin/build/make/go.mak
include bin/build/make/git.mak

# Diagrams generated from https://github.com/loov/goda.
diagrams:
	@make package=. create-diagram

# Log into pulumi.
pulumi-login:
	@pulumi login --cloud-url https://api.pulumi.com

# Preview pulumi changes.
pulumi-preview:
	@pulumi preview --stack alexfalkowski/$(area)/prod --cwd area/$(area) --diff

# Update pulumi changes.
pulumi-update:
	@pulumi update --yes --stack alexfalkowski/$(area)/prod --cwd area/$(area)

# Cancel pulumi changes.
pulumi-cancel:
	@pulumi cancel --yes --stack alexfalkowski/$(area)/prod --cwd area/$(area)

# Delete pulumi stack.
pulumi-delete:
	@pulumi stack rm --yes --force --stack alexfalkowski/$(area)/prod --cwd area/$(area)

# Lint the API.
api-lint:
	@make -C api lint

# Check the API for breaking changes.
api-breaking:
	@make -C api breaking

# Generate the API.
api-generate:
	@make -C api generate

# Build bump.
build-bump: dep
	@go build -ldflags="-s -w" -mod vendor -tags netgo -a cmd/bump/bump.go

# Build format.
build-format: dep
	@go build -ldflags="-s -w" -mod vendor -tags netgo -a cmd/format/format.go
