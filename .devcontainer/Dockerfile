# Development / test container for rclgo.
#
# Mirrors the GitHub Actions CI environment (see .github/workflows/main.yaml):
# ROS 2 Lyrical as the base, plus Go, make and the ROS interface packages the
# test suite depends on. Build args let you bump versions without editing the
# file.
ARG ROS_DISTRO=lyrical
FROM ros:${ROS_DISTRO}-ros-base

ARG ROS_DISTRO=lyrical
ENV ROS_DISTRO=${ROS_DISTRO}

# Go toolchain. go.mod declares `go 1.26`; with GOTOOLCHAIN=auto (the default)
# the bootstrap toolchain installed here will fetch the exact version go.mod
# asks for on first build. Override GO_VERSION at build time to pin a newer
# bootstrap.
ARG GO_VERSION=1.27.3
ARG GOLANGCI_LINT_VERSION=v2.15.0

ENV DEBIAN_FRONTEND=noninteractive

# Base build tooling and the ROS interface packages exercised by the tests.
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        build-essential \
        ca-certificates \
        curl \
        git \
        make \
        pkg-config \
        sudo \
        ros-${ROS_DISTRO}-action-msgs \
        ros-${ROS_DISTRO}-example-interfaces \
        ros-${ROS_DISTRO}-test-msgs \
        ros-${ROS_DISTRO}-rmw-zenoh-cpp \
    && rm -rf /var/lib/apt/lists/*

# Install Go from the official tarball (the apt golang is usually too old).
RUN ARCH="$(dpkg --print-architecture)" \
    && case "${ARCH}" in \
        amd64) GOARCH=amd64 ;; \
        arm64) GOARCH=arm64 ;; \
        *) echo "unsupported arch: ${ARCH}" >&2; exit 1 ;; \
    esac \
    && curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-${GOARCH}.tar.gz" -o /tmp/go.tgz \
    && tar -C /usr/local -xzf /tmp/go.tgz \
    && rm /tmp/go.tgz

# Non-root development user. The devcontainer runs as this user, and the
# devcontainer tooling (updateRemoteUserUID) remaps its UID/GID to match the
# host user at container start, so files created in the bind-mounted workspace
# are owned by you rather than root. Ubuntu noble based images already ship an
# `ubuntu` user at UID 1000; create the user if the base image lacks it.
ARG USERNAME=ubuntu
ARG USER_UID=1000
ARG USER_GID=${USER_UID}
RUN if ! id -u "${USERNAME}" >/dev/null 2>&1; then \
        groupadd --gid "${USER_GID}" "${USERNAME}" \
        && useradd --uid "${USER_UID}" --gid "${USER_GID}" -m -s /bin/bash "${USERNAME}"; \
    fi \
    && echo "${USERNAME} ALL=(root) NOPASSWD:ALL" > "/etc/sudoers.d/${USERNAME}" \
    && chmod 0440 "/etc/sudoers.d/${USERNAME}" \
    && mkdir -p "/home/${USERNAME}/go" \
    && chown -R "${USERNAME}:" "/home/${USERNAME}/go"

ENV GOPATH=/home/${USERNAME}/go
ENV PATH=/usr/local/go/bin:${GOPATH}/bin:${PATH}

# golangci-lint for `make lint`.
RUN curl -fsSL https://raw.githubusercontent.com/golangci/golangci-lint/main/install.sh \
        | sh -s -- -b /usr/local/bin "${GOLANGCI_LINT_VERSION}"

# cgocheck=0 is required by rclgo's cgo bindings (see ros2.env).
ENV GODEBUG=cgocheck=0
ENV GOFLAGS=-buildvcs=false

# Current Kilted rosidl headers transitively include "rcutils/allocator.h", but
# the per-package ROS include layout puts it under .../include/rcutils/rcutils/.
# CGO_CFLAGS is appended to each package's own #cgo flags, so adding the rcutils
# roots here lets test fixtures with hard-coded include lists still compile.
ENV CGO_CFLAGS="-I/opt/ros/${ROS_DISTRO}/include/rcutils -I/usr/include/rcutils"

# Use the Zenoh middleware as the default RMW. Note: rmw_zenoh_cpp needs a Zenoh
# router for discovery — start one in the container with
#   ros2 run rmw_zenoh_cpp rmw_zenohd
# (or unset RMW_IMPLEMENTATION to fall back to the default Fast-DDS).
ENV RMW_IMPLEMENTATION=rmw_zenoh_cpp

# Source the ROS 2 environment for every interactive bash shell so AMENT_PREFIX_PATH,
# LD_LIBRARY_PATH, etc. are set when running tests by hand.
RUN echo 'source /opt/ros/${ROS_DISTRO}/setup.bash' >> /etc/bash.bashrc

# Entrypoint that sources ROS 2 before exec'ing the given command, so non-login
# invocations (e.g. `docker run rclgo make test`) get the environment too.
COPY ros_entrypoint.sh /usr/local/bin/ros_entrypoint.sh
RUN chmod +x /usr/local/bin/ros_entrypoint.sh
ENTRYPOINT ["/usr/local/bin/ros_entrypoint.sh"]
CMD ["bash"]
