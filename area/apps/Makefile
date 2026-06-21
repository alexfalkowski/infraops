include ../../bin/build/make/help.mak
include lean.mk

# Run kube score.
kube-score: kube-score-lean

# Run kubescape.
kubescape: kubescape-lean

# Lint live app resources.
lint: kube-score kubescape

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
