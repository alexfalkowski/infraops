override export github_token_file := $(value github_token_file)

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

# Create Abgleich's GitHub token secret in lean (requires github_token_file).
abgleich-secret:
	@test -n "$$github_token_file" || { printf '%s\n' 'github_token_file is required' >&2; exit 1; }
	@kubectl create secret generic abgleich-secret --namespace lean --from-file="github-token=$$github_token_file"

# Setup lean.
setup-lean: create-lean
