override export github_app_private_key_file := $(value github_app_private_key_file)
override export github_webhook_secret_file := $(value github_webhook_secret_file)

# Run kube-score for lean.
kube-score-lean:
	@./kube-score lean

# Run kubescape for lean.
kubescape-lean:
	@kubescape scan --include-namespaces lean

# Delete lean.
delete-lean:
	@kubectl delete namespaces lean

# Create lean.
create-lean:
	@kubectl create namespace lean

# Create Abgleich's GitHub secret in lean (requires github_app_private_key_file and github_webhook_secret_file).
abgleich-secret:
	@./lean/abgleich

# Setup lean.
setup-lean: create-lean
