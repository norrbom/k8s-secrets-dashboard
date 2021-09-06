FROM golang:1.17-alpine3.14

# Install awscli
RUN apk add --update --no-cache python3 && \
    python3 -m ensurepip && \
    pip3 install --upgrade pip && \
    pip3 install awscli && \
    pip3 cache purge

# https://docs.aws.amazon.com/eks/latest/userguide/install-aws-iam-authenticator.html
# install aws iam authenticator
RUN aws --no-sign-request s3 cp s3://amazon-eks/1.21.2/2021-09-02/bin/linux/amd64/aws-iam-authenticator /usr/bin/aws-iam-authenticator && \
chmod +x /usr/bin/aws-iam-authenticator

WORKDIR /app
COPY go.mod go.sum ./
# download dependencies early to speed up local builds
RUN go mod download
COPY cmd/*.go templates ./
RUN CGO_ENABLED=0 GOOS=linux go test -v ./...
RUN go build -o /app/bin/zcdash

CMD ["/app/bin/zcdash"]
