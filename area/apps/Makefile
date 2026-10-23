include ../../bin/build/make/help.mak
include lean.mk

.PHONY: kube-score

# Score apps changed by the last successful update.
kube-score:
	@./release kube-score

# Score all live resources in the namespace.
kube-score-all:
	@./release kube-score-all

# Run kubescape.
kubescape:
	@./release kubescape

# Run kubescape for all apps.
kubescape-all:
	@./release kubescape-all

# Lint apps changed by the last successful update.
lint: kube-score kubescape

# Lint the namespace and scan every supported app image.
lint-all: kube-score-all kubescape-all

# Save kubeconfig.
save-config:
	@doctl kubernetes cluster kubeconfig save lean-thoughts

# Setup all apps.
setup: setup-lean

# Delete all apps.
delete: delete-lean

# Rollout all apps.
rollout-all:
	@./release rollout-all

# Rollout release app changes.
rollout:
	@./release rollout

# Verify all apps.
verify-all:
	@./release verify-all

# Verify release app changes.
verify:
	@./release verify

# Load all apps.
load-all:
	@./release load-all

# Load release app changes.
load:
	@./release load
