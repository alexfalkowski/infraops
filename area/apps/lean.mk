otlp_secret := $(shell echo -n "Basic 790760:$(GRAFANA_OTLP_TOKEN)" | base64 -w 0)

# Run kubescore for lean.
kube-score-lean:
	kubectl api-resources --verbs=list --namespaced -o name \
		| xargs -I{} bash -c "kubectl get {} --namespace lean -oyaml && echo ---" \
		| kube-score score --ignore-test deployment-has-host-podantiaffinity  -

# Delete lean.
delete-lean:
	kubectl delete namespaces lean

# Setup lean.
setup-lean:
	kubectl create namespace lean
	kubectl create secret generic otlp-secret --from-literal=token=$(otlp_secret) --namespace lean
	kubectl create secret generic konfig-secret --from-literal=token=$(KONFIG_TOKEN) --namespace lean
	kubectl create secret generic gh-secret --from-literal=token=$(GITHUB_TOKEN) --namespace lean

# Rollout lean.
rollout-lean: rollout-konfig rollout-standort rollout-bezeichner rollout-web

# Rollout konfig.
rollout-konfig:
	kubectl rollout restart deployment/konfig -n lean

# Rollout standort.
rollout-standort:
	kubectl rollout restart deployment/standort -n lean

# Rollout bezeichner.
rollout-bezeichner:
	kubectl rollout restart deployment/bezeichner -n lean

# Rollout web.
rollout-web:
	kubectl rollout restart deployment/web -n lean

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
