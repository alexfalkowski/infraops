# Save kubeconfig.
save-config:
	doctl kubernetes cluster kubeconfig save f7a7f2b2-f562-49fe-93e7-2b6ef3038fa5

# Delete all.
delete: delete-standort

# Setup all apps.
setup: setup-standort

# Delete standort.
delete-standort:
	kubectl delete namespaces standort

# Setup standort.
setup-standort:
	kubectl create namespace standort

# Preview standort.
preview-standort:
	kubectl apply -f standort -n standort --dry-run=server

# Update standort.
update-standort:
	kubectl apply -f standort -n standort
