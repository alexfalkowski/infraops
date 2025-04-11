# Run kubescore for lean.
kube-score-lean:
	@kubectl api-resources --verbs=list --namespaced -o name \
		| xargs -I{} bash -c "kubectl get {} --namespace lean -oyaml && echo ---" \
		| kube-score score --ignore-test deployment-has-host-podantiaffinity  -

# Run kubescape for lean.
kubescape-lean:
	@kubescape scan --include-namespaces lean

# Delete lean.
delete-lean:
	@kubectl delete namespaces lean

# Create lean
create-lean:
	@kubectl create namespace lean

# Setup otlp.
setup-otlp:
	@kubectl create secret generic otlp-secret --from-file ~/keys/otlp --namespace lean

# Setup konfig.
setup-konfig:
	@kubectl create secret generic konfig-secret --from-file ~/keys/konfig --namespace lean

# Setup Github.
setup-gh:
	@kubectl create secret generic gh-secret --from-literal=token=$(GITHUB_TOKEN) --namespace lean

# Setup lean.
setup-lean: create-lean setup-otlp setup-konfig setup-gh

# Rollout lean.
rollout-lean: rollout-konfig rollout-standort rollout-bezeichner rollout-web

# Rollout konfig.
rollout-konfig:
	@kubectl rollout restart deployment/konfig -n lean

# Rollout standort.
rollout-standort:
	@kubectl rollout restart deployment/standort -n lean

# Rollout bezeichner.
rollout-bezeichner:
	@kubectl rollout restart deployment/bezeichner -n lean

# Rollout web.
rollout-web:
	@kubectl rollout restart deployment/web -n lean

# Verify all apps.
verify-lean: verify-standort verify-bezeichner verify-web

# Verify standort.
verify-standort:
	@curl -svf --header "Content-Type: application/json" --request POST --data {}  https://standort.lean-thoughts.com/standort.v2.Service/GetLocation

# Verify bezeichner.
verify-bezeichner:
	@curl -svf --header "Content-Type: application/json" --request POST --data '{ "application": "ulid", "count": 10 }'  https://bezeichner.lean-thoughts.com/bezeichner.v1.Service/GenerateIdentifiers

# Verify web.
verify-web:
	@curl -svf https://web.lean-thoughts.com
