otlp_secret := $(shell echo -n "790760:$(GRAFANA_OTLP_TOKEN)" | base64 -w 0)

# Run kubescore for lean.
kube-score-lean:
	$(MAKE) namespace=lean kube-score-namespace

# Delete lean.
delete-lean:
	$(MAKE) namespace=lean delete-namespace

# Setup lean.
setup-lean:
	$(MAKE) namespace=lean setup-namespace
	kubectl create secret generic otlp-secret --from-literal=token=$(otlp_secret) --namespace lean
	kubectl create secret generic konfig-secret --from-literal=token=$(KONFIG_TOKEN) --namespace lean
	kubectl create secret generic gh-secret --from-literal=token=$(GITHUB_TOKEN) --namespace lean

# Preview lean.
preview-lean:
	$(MAKE) namespace=lean preview-namespace

# Update lean.
update-lean:
	$(MAKE) namespace=lean update-namespace

# Rollout lean.
rollout-lean:
	$(MAKE) namespace=lean rollout-namespace

# Verify all apps.
verify-lean: verify-standort verify-bezeichner verify-web

# Verify standort.
verify-standort:
	curl -svf --header "Content-Type: application/json" --request POST --data {}  https://standort.lean-thoughts.com/v2/location

# Verify bezeichner.
verify-bezeichner:
	curl -svf --header "Content-Type: application/json" --request POST --data '{ "application": "uuid", "count": 10 }'  https://bezeichner.lean-thoughts.com/v1/generate

# Verify web.
verify-web:
	curl -svf https://web.lean-thoughts.com
