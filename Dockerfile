###########################
# Global build args
###########################
# change service and workdir with your 
ARG service=GoBackbone 
ARG workdir=/go/src/github.com/orov.io
ARG key_rsa=./id_rsa
############################
# STEP 1 build executable binary
############################

FROM golang:alpine 

ARG service
ARG workdir
ARG key_rsa

WORKDIR ${workdir}/${service}
COPY . ${workdir}/${service}

# Install git.
## Git is required for fetching dependencies.
RUN apk update && apk add --no-cache git openssh-client

## Update ssh hosts
#Uncomment this if you have private bitbucket dependencies
#RUN mkdir -p /root/.ssh && \
#    chmod 0700 /root/.ssh && \
#    ssh-keyscan bitbucket.org > /root/.ssh/known_hosts
## Add the keys
#COPY ${key_rsa} /root/.ssh/id_rsa

# Update project dependencies
## Installig dep
RUN go get -u github.com/golang/dep/...
## Updating dependencies
RUN cd ${workdir}/${service}
RUN dep ensure -v

# Building the app
RUN go build -o main

# Configuring the server
EXPOSE 8080
ENTRYPOINT "./main"