# Run kubescore for standort.
kube-score-standort:
	$(MAKE) namespace=standort kube-score-namespace

# Delete standort.
delete-standort:
	$(MAKE) namespace=standort delete-namespace

# Setup standort.
setup-standort:
	$(MAKE) namespace=standort setup-namespace
	$(MAKE) namespace=standort setup-otlp-secret
	$(MAKE) namespace=standort setup-konfig-secret

# Rollout standort.
rollout-standort:
	$(MAKE) namespace=standort rollout-namespace

# Verify standort.
verify-standort:
	curl -svf --header "Content-Type: application/json" --request POST --data {}  https://standort.lean-thoughts.com/v2/location
