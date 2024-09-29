FROM alpine:3.20.3
ARG VERSION
ARG BUILDUSER
LABEL PROJECT="yuan-jewel"    \
      VERSION="$VERSION"            \
      AUTHOR="$BUILDUSER"
MAINTAINER yuanJewel "luyu151111@gamil.com"
ENV LC_ALL en_US.UTF-8
ADD build/yuan /bin/yuan
RUN chmod +x /bin/yuan
WORKDIR /opt
ENTRYPOINT ["yuan"]
