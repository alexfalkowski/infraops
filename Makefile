include bin/build/make/go.mak
include bin/build/make/git.mak

# Diagrams generated from https://github.com/loov/goda.
diagrams:
	$(MAKE) package=. create-diagram

# Log into pulumi.
pulumi-login:
	pulumi login --cloud-url https://api.pulumi.com

# Preview pulumi changes.
pulumi-preview:
	pulumi preview --stack alexfalkowski/$(area)/prod --cwd area/$(area) --diff

# Update pulumi changes.
pulumi-update:
	pulumi update --yes --stack alexfalkowski/$(area)/prod --cwd area/$(area)
