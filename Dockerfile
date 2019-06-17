###########################
# Global build args
###########################
# change service and workdir with your 
ARG service=${SERVICE_NAME}
ARG repo=${PROJECT_PATH}
ARG workdir=${repo}/${service}
ARG key_rsa=./id_rsa
############################
# STEP 1 build executable binary
############################

FROM golang:alpine AS build-env

ARG service
ARG workdir
ARG key_rsa

WORKDIR ${workdir}
COPY . ${workdir}


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
RUN cd ${workdir}
RUN dep ensure -v

# Building the app
RUN go build -o app

############################
# STEP 2 build tiny executable container
############################

FROM alpine

ARG service
ARG workdir
ARG key_rsa

WORKDIR ${workdir}
COPY --from=build-env ${workdir}/app ${workdir}/
COPY --from=build-env ${workdir}/migrations/* ${workdir}/migrations/
RUN apk --update add ca-certificates

EXPOSE 8080
ENTRYPOINT ./${workdir}/app