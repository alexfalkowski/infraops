include bin/build/make/help.mak
include bin/build/make/go.mak
include bin/build/make/git.mak

# Log into Pulumi.
pulumi-login:
	@pulumi login --cloud-url https://api.pulumi.com

# Preview Pulumi changes.
pulumi-preview:
	@pulumi preview --stack alexfalkowski/$(area)/prod --cwd area/$(area) --diff

# Update Pulumi changes.
pulumi-update:
	@pulumi update --yes --stack alexfalkowski/$(area)/prod --cwd area/$(area)

# Cancel Pulumi changes.
pulumi-cancel:
	@pulumi cancel --yes --stack alexfalkowski/$(area)/prod --cwd area/$(area)

# Delete Pulumi stack.
pulumi-delete:
	@pulumi stack rm --yes --force --stack alexfalkowski/$(area)/prod --cwd area/$(area)

# Refresh Pulumi stack.
pulumi-refresh:
	@pulumi refresh --yes --stack alexfalkowski/$(area)/prod --cwd area/$(area)

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
