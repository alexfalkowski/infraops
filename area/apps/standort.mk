# Run kubescore for standort.
kube-score-standort:
	$(MAKE) namespace=standort kube-score-namespace

# Delete standort.
delete-standort:
	$(MAKE) namespace=standort delete-namespace

# Setup standort.
setup-standort:
	$(MAKE) namespace=standort setup-namespace
	$(MAKE) namespace=standort setup-otlp

# Rollout standort.
rollout-standort:
	$(MAKE) namespace=standort rollout-namespace

# Verify standort.
verify-standort:
	curl -sf https://standort.lean-thoughts.com/v2/location
