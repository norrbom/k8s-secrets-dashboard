## k8s Secrets Dashboard
Page showing progress of services moving over from Sealed Secrets to Vault and Zero Config.

### Prerequisites

* terraform 1
* go 1.17 or greater

### Test
```bash
go test ./...
```
### Manage Dependencies
```bash
go get -v ./cmd
# commit go.mod and go.sum
```
### Run
```bash
export KUBECONFIG_TEST=<PATH_TO_KUEBCONFIG_FILE>
export KUBECONFIG_PROD=<PATH_TO_KUEBCONFIG_FILE>
export GIT_USER=<GIT_USER>
export GIT_PASS=<GIT_PASS>
go run main.go report.go gitops.go utils.go config.go
```
### Docker
```bash
docker build -t zcreport . && \
docker run -it --rm -e GIT_USER=<GIT_USER> -e GIT_PASS=<GIT_PASS> \
-e KUBECONFIG_TEST=/kubernetes-configuration/<KUEBCONFIG_FILE> \
-e KUBECONFIG_PROD=/kubernetes-configuration/<KUEBCONFIG_FILE> \
-v $(pwd)/data:/app/data -v $(pwd)/templates:/app/templates -v <SOMEPATH>:/kubernetes-configuration zcreport
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
