FROM ubuntu:latest

ARG KUBECTL_VERSION="v1.25.0"
ARG HELM_VERSION="v3.11.2"
ARG AWSCLI_VERSION="2.11.5"
ARG GO_VERSION="1.20.2"

RUN apt-get update && apt-get install -y curl git unzip \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN curl -LO https://dl.k8s.io/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl \
    && chmod +x ./kubectl \
    && mv ./kubectl /usr/local/bin/kubectl

RUN curl https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash -s -- --version ${HELM_VERSION}

RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64-${AWSCLI_VERSION}.zip" -o "awscliv2.zip" \
    && unzip awscliv2.zip \
    && ./aws/install \
    && rm -rf awscliv2.zip ./aws

RUN curl https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz | tar -xz -C /usr/local \
    && rm -rf /usr/local/go/doc /usr/local/go/test

ENV PATH="$PATH:/usr/local/go/bin"

# Set the default shell to bash
SHELL ["/bin/bash", "-c"]
