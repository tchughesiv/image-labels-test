# reference file - https://github.com/containers/buildah/blob/v1.14.10/contrib/buildahimage/stable/Dockerfile
FROM registry.access.redhat.com/ubi8

RUN useradd build && yum -y update --security --sec-severity=Important --sec-severity=Critical && \
        yum -y install containers-common fuse-overlayfs --setopt=tsflags=nodocs && \
        rm -rf /var/cache /var/log/yum.*

# Adjust storage.conf to enable Fuse storage.
RUN sed -i -e 's|^#mount_program|mount_program|g' -e '/additionalimage.*/a "/var/lib/shared",' -e 's|^# mountopt[[:space:]]*=.*$|mountopt = "nodev,fsync=0"|g' /etc/containers/storage.conf

# Set up environment variables to note that this is
# not starting with usernamespace and default to
# isolate the filesystem with chroot.
ENV _BUILDAH_STARTED_IN_USERNS="" BUILDAH_ISOLATION=chroot

# install operator binary
COPY image-labels-test /usr/local/bin/image-labels-test
