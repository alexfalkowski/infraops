include ../../bin/build/make/help.mak

# Save kubeconfig.
save-config:
	@doctl kubernetes cluster kubeconfig save e18082d3-d874-4e84-8703-c15e6d6fea3a

# Delete add-on namespaces.
delete:
	@kubectl delete namespaces contour circleci metrics-server

# Set up Helm repos and install add-ons (requires CIRCLECI_K8S_TOKEN).
setup: add-helm-repos update-helm-repos install-helm

# Add Helm repos.
add-helm-repos:
	@helm repo add contour https://projectcontour.github.io/helm-charts/
	@helm repo add circleci https://circleci-public.github.io/cci-k8s-release-agent
	@helm repo add metrics-server https://kubernetes-sigs.github.io/metrics-server/

# Update Helm repos.
update-helm-repos:
	@helm repo update

# Install CircleCI release agent (requires CIRCLECI_K8S_TOKEN).
install-circleci:
	@helm upgrade --install circleci-release circleci/circleci-release-agent --set tokenSecret.token=$(CIRCLECI_K8S_TOKEN) --create-namespace --namespace circleci --set managedNamespaces="{lean}" --set rbac.clusterRoles=false --version 1.8.2 -f circleci/values.yaml --wait --timeout 5m

# Install metrics server.
install-metrics-server:
	@helm upgrade --install metrics-server-release metrics-server/metrics-server --create-namespace --namespace metrics-server --version 3.14.0 --wait --timeout 5m

# Create the shared self-signed ingress certificate in the existing lean namespace.
create-ingress-cert:
	@./contour/certificate

# Install Contour and Envoy.
install-contour:
	@helm upgrade --install contour-release contour/contour --namespace contour --create-namespace --version 0.9.0 -f contour/values.yaml --wait --timeout 5m

# Render the pinned Contour chart without changing the cluster.
render-contour:
	@helm template contour-release contour/contour --namespace contour --version 0.9.0 -f contour/values.yaml

# Show Envoy's public load-balancer address for Cloudflare DNS.
contour-address:
	@kubectl get service contour-release-envoy --namespace contour -o wide

# Install Helm add-ons.
install-helm: install-circleci install-metrics-server install-contour

# List all pods.
pods:
	@kubectl get pods --all-namespaces
