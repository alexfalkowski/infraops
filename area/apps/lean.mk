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

# Setup Github.
setup-gh:
	@kubectl create secret generic gh-secret --from-literal=token=$(GITHUB_TOKEN) --namespace lean

# Setup lean.
setup-lean: create-lean setup-gh

# Rollout lean.
rollout-lean: rollout-standort rollout-bezeichner rollout-web rollout-monitoror

# Rollout standort.
rollout-standort:
	@kubectl rollout restart deployment/standort -n lean

# Rollout bezeichner.
rollout-bezeichner:
	@kubectl rollout restart deployment/bezeichner -n lean

# Rollout web.
rollout-web:
	@kubectl rollout restart deployment/web -n lean

# Rollout monitoror.
rollout-monitoror:
	@kubectl rollout restart deployment/monitoror -n lean

# Verify all apps.
verify-lean: verify-standort verify-bezeichner verify-web verify-monitoror

# Verify standort.
verify-standort:
	@curl -svf --header "Content-Type: application/json" --request POST --data {}  https://standort.lean-thoughts.com/standort.v2.Service/GetLocation

# Verify bezeichner.
verify-bezeichner:
	@curl -svf --header "Content-Type: application/json" --request POST --data '{ "application": "ulid", "count": 10 }'  https://bezeichner.lean-thoughts.com/bezeichner.v1.Service/GenerateIdentifiers

# Verify web.
verify-web:
	@curl -svf https://web.lean-thoughts.com

# Verify monitoror.
verify-monitoror:
	@curl -svf https://monitoror.lean-thoughts.com

# Load test lean.
load-lean: load-standort load-bezeichner load-web load-monitoror

# Load test standort.
load-standort:
	@echo "POST https://standort.lean-thoughts.com/standort.v2.Service/GetLocation" | vegeta attack -duration=30s -body "lean/standort.json" -header "Content-Type: application/json" | tee "lean/standort.bin" | vegeta report

# Load test bezeichner.
load-bezeichner:
	@echo "POST https://bezeichner.lean-thoughts.com/bezeichner.v1.Service/GenerateIdentifiers" | vegeta attack -duration=30s -body "lean/bezeichner.json" -header "Content-Type: application/json" | tee "lean/bezeichner.bin" | vegeta report

# Load test web.
load-web:
	@echo "GET https://web.lean-thoughts.com" | vegeta attack -duration=30s | tee "lean/web.bin" | vegeta report

# Load test monitoror.
load-monitoror:
	@echo "GET https://monitoror.lean-thoughts.com" | vegeta attack -duration=30s | tee "lean/monitoror.bin" | vegeta report
