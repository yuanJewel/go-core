FROM alpine:3.20.3
ARG VERSION
ARG BUILDUSER
LABEL PROJECT="SmartLyu-go-core"    \
      VERSION="$VERSION"            \
      AUTHOR="$BUILDUSER"
MAINTAINER SmartLyu "luyu151111@gamil.com"
ENV LC_ALL en_US.UTF-8
ADD build/go-core /bin/go-core
RUN chmod +x /bin/go-core
ENTRYPOINT ["go-core"]
