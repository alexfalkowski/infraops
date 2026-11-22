override export github_app_private_key_file := $(value github_app_private_key_file)
override export github_webhook_secret_file := $(value github_webhook_secret_file)
override export openai_api_key_file := $(value openai_api_key_file)
override export r2_access_key_id_file := $(value r2_access_key_id_file)
override export r2_secret_access_key_file := $(value r2_secret_access_key_file)

# Run kube-score for lean.
kube-score-lean:
	@./kube-score lean

# Check lean's live configuration without image or host scans.
kubescape-lean:
	@./kubescape --include-namespaces lean

# Delete lean.
delete-lean:
	@kubectl delete namespaces lean

# Create lean.
create-lean:
	@kubectl create namespace lean

# Create Abgleich's GitHub and OpenAI secret in lean (requires github_app_private_key_file, github_webhook_secret_file, and openai_api_key_file).
abgleich-secret:
	@./lean/abgleich

# Create Einblick's R2 secret in lean (requires r2_access_key_id_file and r2_secret_access_key_file).
einblick-secret:
	@./lean/einblick

# Setup lean.
setup-lean: create-lean
