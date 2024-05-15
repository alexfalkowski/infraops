# Save kubeconfig.
save-config:
	doctl kubernetes cluster kubeconfig save f7a7f2b2-f562-49fe-93e7-2b6ef3038fa5

# Setup all repos.
setup:
	helm repo add nginx-stable https://helm.nginx.com/stable
	helm repo update

# Preview deps.
preview:
	helm install nginx-ingress-release nginx-stable/nginx-ingress --namespace nginx-ingress --create-namespace --dry-run


# Update deps.
update:
	helm install nginx-ingress-release nginx-stable/nginx-ingress --namespace nginx-ingress --create-namespace
<<<<<<< Updated upstream
=======
	helm install cert-manager-release jetstack/cert-manager --namespace cert-manager --create-namespace --version v1.14.5 --set installCRDs=true
>>>>>>> Stashed changes

