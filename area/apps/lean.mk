# Run kube-score for lean.
kube-score-lean:
	@./kube-score-live lean

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
