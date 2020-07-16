# reference file - https://github.com/containers/buildah/blob/v1.14.10/contrib/buildahimage/stable/Dockerfile
FROM registry.access.redhat.com/ubi8

RUN useradd build && yum -y update --security --sec-severity=Important --sec-severity=Critical && \
        yum -y install fuse-overlayfs --setopt=tsflags=nodocs && \
        rm -rf /var/cache /var/log/yum.*

# Set up environment variables to note that this is
# not starting with usernamespace and default to
# isolate the filesystem with chroot.
ENV _BUILDAH_STARTED_IN_USERNS="" BUILDAH_ISOLATION=chroot

# install operator binary
COPY image-labels-test /usr/local/bin/image-labels-test
