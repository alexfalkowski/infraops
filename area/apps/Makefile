include lean.mk

# Run kube score.
kube-score: kube-score-lean

# Run kubescape.
kubescape: kubescape-lean

# Lint config.
lint: kube-score kubescape

# Save kubeconfig.
save-config:
	@doctl kubernetes cluster kubeconfig save eb09f796-aab5-4da5-bc79-7b9c1e8a4c2c

# Setup all apps.
setup: setup-lean

# Delete all apps.
delete: delete-lean

# Rollout all apps.
rollout: rollout-lean

# Verify all apps.
verify: verify-lean

# Load all apps.
load: load-lean
