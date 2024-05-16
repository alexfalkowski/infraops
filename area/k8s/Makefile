# Save kubeconfig.
save-config:
	doctl kubernetes cluster kubeconfig save f7a7f2b2-f562-49fe-93e7-2b6ef3038fa5

# Setup all repos.
setup:
	helm repo add nginx-stable https://helm.nginx.com/stable --force-update
	helm repo add jetstack https://charts.jetstack.io --force-update
	helm repo update

# Delete all.
delete:
	kubectl delete namespaces nginx-ingress cert-manager

# Preview deps.
preview:
	helm upgrade --install nginx-ingress-release nginx-stable/nginx-ingress --namespace nginx-ingress --create-namespace --version 1.2.1 --dry-run
	helm upgrade --install cert-manager-release jetstack/cert-manager --namespace cert-manager --create-namespace --version v1.14.5 --set installCRDs=true --dry-run
	kubectl apply -f https://raw.githubusercontent.com/cloudflare/origin-ca-issuer/v0.7.0/deploy/crds/cert-manager.k8s.cloudflare.com_originissuers.yaml -n cert-manager --dry-run=server
	helm upgrade --install origin-ca-issuer-release oci://ghcr.io/cloudflare/origin-ca-issuer-charts/origin-ca-issuer --namespace cert-manager --version 0.5.2 --dry-run

# Update deps.
update:
	helm upgrade --install nginx-ingress-release nginx-stable/nginx-ingress --namespace nginx-ingress --create-namespace --version 1.2.1
	helm upgrade --install cert-manager-release jetstack/cert-manager --namespace cert-manager --create-namespace --version v1.14.5 --set installCRDs=true
	kubectl apply -f https://raw.githubusercontent.com/cloudflare/origin-ca-issuer/v0.7.0/deploy/crds/cert-manager.k8s.cloudflare.com_originissuers.yaml -n cert-manager
	helm upgrade --install origin-ca-issuer-release oci://ghcr.io/cloudflare/origin-ca-issuer-charts/origin-ca-issuer --namespace cert-manager --version 0.5.2
