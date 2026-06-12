# Run kube-score for lean.
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

# Create lean.
create-lean:
	@kubectl create namespace lean

# Setup lean.
setup-lean: create-lean
