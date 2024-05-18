# Delete bezeichner.
delete-bezeichner:
	$(MAKE) namespace=bezeichner delete-namespace

# Setup bezeichner.
setup-bezeichner:
	$(MAKE) namespace=bezeichner setup-namespace

# Preview bezeichner.
preview-bezeichner:
	$(MAKE) namespace=bezeichner preview-namespace

# Update bezeichner.
update-bezeichner:
	$(MAKE) namespace=bezeichner update-namespace
