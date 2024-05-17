# Save kubeconfig.
save-config:
	doctl kubernetes cluster kubeconfig save f7a7f2b2-f562-49fe-93e7-2b6ef3038fa5

# Delete all.
delete:
	kubectl delete namespaces nginx-ingress cert-manager grafana

# Setup everything.
setup: setup-helm setup-grafana

# Setup helm repos.
setup-helm:
	helm repo add nginx-stable https://helm.nginx.com/stable
	helm repo add jetstack https://charts.jetstack.io
	helm repo add grafana https://grafana.github.io/helm-charts
	helm repo update

# Setup grafana.
setup-grafana:
	kubectl create secret generic metrics-secret --from-literal=username=1289421 --from-literal=password=$(GRAFANA_METRICS_TOKEN) --namespace grafana
	kubectl create secret generic logs-secret --from-literal=username=743756 --from-literal=password=$(GRAFANA_LOGS_TOKEN) --namespace grafana

# Preview grafana.
preview-grafana:
	kubectl apply -f grafana/kube.yaml --namespace grafana --dry-run=server
	kubectl apply -f grafana/cluster.yaml --namespace grafana --dry-run=server
	kubectl apply -f grafana/metrics.yaml --namespace grafana --dry-run=server
	kubectl apply -f grafana/logs.yaml --namespace grafana --dry-run=server

# Update grafana.
update-grafana:
	kubectl apply -f grafana/kube.yaml --namespace grafana
	kubectl apply -f grafana/cluster.yaml --namespace grafana
	kubectl apply -f grafana/metrics.yaml --namespace grafana
	kubectl apply -f grafana/logs.yaml --namespace grafana

# Preview helm.
preview-helm:
	helm upgrade --install nginx-ingress-release nginx-stable/nginx-ingress --namespace nginx-ingress --create-namespace --version 1.2.1 --dry-run
	helm upgrade --install cert-manager-release jetstack/cert-manager --namespace cert-manager --create-namespace --version v1.14.5 --set installCRDs=true --dry-run
	helm upgrade --install grafana-agent-release --create-namespace grafana/grafana-agent-operator -n grafana --version v0.3.22 --dry-run

# Update helm.
update-helm:
	helm upgrade --install nginx-ingress-release nginx-stable/nginx-ingress --namespace nginx-ingress --create-namespace --version 1.2.1
	helm upgrade --install cert-manager-release jetstack/cert-manager --namespace cert-manager --create-namespace --version v1.14.5 --set installCRDs=true
	helm upgrade --install grafana-agent-release --create-namespace grafana/grafana-agent-operator -n grafana --version v0.3.22

# Preview deps.
preview: preview-helm preview-grafana

# Update deps.
update: update-helm update-grafana
