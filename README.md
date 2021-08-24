## k8s Secrets Dashboard
A dashboard showing Sealed Secrets vs Vault secrets, using naming conventions.

### Prerequisites

* terraform 1
* go 1.17 or greater

### Run
```bash
export KUBECONFIG_SI1=<PATH_TO_KUEBCONFIG_FILE>
go run main.go report.go utils.go
```
### Docker
```bash
docker build -t zcreport . && \
docker run -it --rm -e KUBECONFIG_SI1=/kubernetes-configuration/<KUEBCONFIG_FILE> \
-v $(pwd)/data:/app/data -v $(pwd)/templates:/app/templates \
-v <SOMEPATH>/kubernetes-configuration:/kubernetes-configuration zcreport
```
### Publish to S3 DEV
```bash
docker run \
-w /apps/terraform/dev/static-web-site \
--rm -it \
-e AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID \
-e AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY \
-v ${PWD}:/apps alpine/terragrunt:1.0.5 \
/bin/sh -c "cp /apps/data/* ./content && terragrunt apply --terragrunt-non-interactive -auto-approve"
```
