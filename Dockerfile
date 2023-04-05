FROM ubuntu:latest

ARG GO_VERSION="1.20.2"
ARG KUBECTL_VERSION="1.25.0"
ARG HELM_VERSION="v3.11.2"
ARG AWSCLI_VERSION="2.11.5"

RUN apt-get update && apt-get install -y curl git unzip \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN curl https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz | tar -xz -C /usr/local \
    && export PATH=$PATH:/usr/local/go/bin \
    && echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.bashrc \
    && rm -rf /usr/local/go/doc /usr/local/go/test

RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl \
    && chmod +x ./kubectl \
    && mv ./kubectl /usr/local/bin/kubectl

RUN curl https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash -s -- --version ${HELM_VERSION}

RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64-${AWSCLI_VERSION}.zip" -o "awscliv2.zip" \
    && unzip awscliv2.zip \
    && ./aws/install \
    && rm -rf awscliv2.zip ./aws

# Set the default shell to bash
SHELL ["/bin/bash", "-c"]
