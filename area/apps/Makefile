include lean.mk

kube-score: kube-score-lean

# Lint config.
lint: kube-score

# Save kubeconfig.
save-config:
	doctl kubernetes cluster kubeconfig save e18082d3-d874-4e84-8703-c15e6d6fea3a

# Setup all apps.
setup: setup-lean

# Delete all apps.
delete: delete-lean

# Rollout all apps.
rollout: rollout-lean

# Verify all apps.
verify: verify-lean
