# Save kubeconfig.
save-config:
	doctl kubernetes cluster kubeconfig save f7a7f2b2-f562-49fe-93e7-2b6ef3038fa5

# Encode a config.
encode-config:
	cat $(kind)/server.yaml | base64 -w 0 | cat ; echo

# Delete all.
delete:
	kubectl delete namespaces standort

# Preview standort.
preview-standort:
	kubectl create namespace standort --dry-run=client
	kubectl apply -f standort/service.yaml -n standort --dry-run=client
	kubectl apply -f standort/ingress.yaml -n standort --dry-run=client

# Update standort.
update-standort:
	kubectl create namespace standort || true
	kubectl apply -f standort/service.yaml -n standort
	kubectl apply -f standort/ingress.yaml -n standort
