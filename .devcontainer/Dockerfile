# Humble base image
FROM ghcr.io/tiiuae/fog-ros-baseimage:sha-7fca749

RUN FOG_DEB_REPO="https://ssrc.jfrog.io/artifactory/ssrc-deb-public-local" \
  && echo "deb [trusted=yes] ${FOG_DEB_REPO} $(lsb_release -cs) fog-sw" > /etc/apt/sources.list.d/fogsw.list \
  && echo "deb [trusted=yes] ${FOG_DEB_REPO} $(lsb_release -cs) fog-sw-sros" > /etc/apt/sources.list.d/fogsw-sros.list \
  && apt-get update \
  && apt-get dist-upgrade -y

RUN apt-get install -y --no-install-recommends --allow-downgrades \
  openssh-client \
  libclang-dev \
  ros-${ROS_DISTRO}-fognav-msgs \
  ros-${ROS_DISTRO}-geographic-msgs \
  ros-${ROS_DISTRO}-geometry-msgs \
  ros-${ROS_DISTRO}-diagnostic-msgs \
  ros-${ROS_DISTRO}-map-msgs \
  ros-${ROS_DISTRO}-fastrtps=2.6.6-1jammy.20230919.195605 \
  ros-${ROS_DISTRO}-rosbag2 \
  ros-${ROS_DISTRO}-diagnostic-msgs \
  && rm -rf /var/lib/apt/lists/*

# FIXME Use your username instead of "developer".
ARG USERNAME=username
ARG USER_UID=1002
ARG USER_GID=$USER_UID

RUN groupadd --gid $USER_GID $USERNAME \
  && useradd --uid $USER_UID --gid $USER_GID -m $USERNAME \
  && apt-get update \
  && apt-get install -y sudo \
  && echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME \
  && chmod 0440 /etc/sudoers.d/$USERNAME

USER $USERNAME
