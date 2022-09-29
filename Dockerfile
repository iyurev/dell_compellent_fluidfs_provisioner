FROM golang:1.11.4
RUN apt update && apt -y install nfs-common
ADD ./build/fluidfs_provisioner /bin/
USER 46
