# /usr/local/bin/start.sh will start the service

FROM fedora:latest

# Pause indefinitely if asked to do so.
ARG OO_PAUSE_ON_BUILD
RUN test "$OO_PAUSE_ON_BUILD" = "true" && while sleep 10; do true; done || :

ADD scripts/ /usr/local/bin/

RUN dnf install -y golang \
                   gcc \
                   git \
    dnf clean all

ENV GOBIN=/bin \
    GOPATH=/go

# Creating mount points for crio and docker sockets and dependencies.
RUN mkdir -p /host/usr/bin \
             /logs \
             /usr/bin \
             /etc/sysconfig && \
    /usr/bin/go get github.com/rhdedgar/container-info && \
    cd /go/src/github.com/rhdedgar/container-info && \
    /usr/bin/go install && \
    cd && \
    rm -rf /go

USER 0

# Start processes
CMD /usr/local/bin/start.sh
