
AWS_REGION ?= ap-northeast-1
AWS_PROFILE ?= default
BUCKET_NAME ?= my-tf-s3-signed

_key = index.txt

##@ Run
_build_dir := tmp
_bin = $(_build_dir)/s3sign
.PHONY: run
run: $(_build_dir)/terraform.tfvars $(_bin) ## Run the s3sign command
	$(_bin) -bucket $(BUCKET_NAME) -profile $(AWS_PROFILE) -key $(_key)

##@ Build
$(_bin): $(shell find . -name '*.go')
	go build -o $(_build_dir)/s3sign .

##@ Clean
.PHONY: clean
clean: ## Remove the intermediate files.
	rm -rf $(_build_dir)

.PHONY: remove-aws
remove-aws: $(_build_dir)/terraform.tfvars ## Delete stacks on AWS
	cd deployments/terraform && terraform destroy -var-file=../../$(_build_dir)/terraform.tfvars

##@ Deploy
.PHONY: deploy-aws
deploy-aws: $(_build_dir)/terraform.tfvars ## Deploy stacks to AWS
	cd deployments/terraform && terraform apply -var-file=../../$(_build_dir)/terraform.tfvars

$(_build_dir)/terraform.tfvars:
	mkdir -p $(_build_dir)
	echo 'aws_region = "$(AWS_REGION)"' >> $@
	echo 'aws_profile = "$(AWS_PROFILE)"' >> $@
	echo 'bucket_name = "$(BUCKET_NAME)"' >> $@
	echo 'key = "$(_key)"' >> $@
	# tfvars must ends with a newline
	echo '' >> $@

##@ Help
.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
