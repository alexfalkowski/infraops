# Save kubeconfig.
save-config:
	doctl kubernetes cluster kubeconfig save f7a7f2b2-f562-49fe-93e7-2b6ef3038fa5

# Encode a config.
encode-config:
	cat $(kind)/server.yml | base64 -w 0 | cat ; echo

# Preview standort.
preview-standort:
	kubectl apply -f standort.yml --dry-run=client

# Update standort.
update-standort:
	kubectl apply -f standort.yml
