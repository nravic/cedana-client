#!/bin/bash

APT_PACKAGES="wget git make curl libnl-3-dev libnet-dev \
    libbsd-dev python-ipaddress libcap-dev \
    libprotobuf-dev libprotobuf-c-dev protobuf-c-compiler \
    protobuf-compiler python3-protobuf"

install_apt_packages() {
    apt-get update
    apt-get install -y $APT_PACKAGES
}

install_code_server() {
    curl -fsSL https://code-server.dev/install.sh | sh
}

install_bats_core() {
    git clone https://github.com/bats-core/bats-core.git
    cd bats-core
    ./install.sh /usr/local
    cd -
}

install_docker() {
    sudo apt-get update
    sudo apt-get install ca-certificates curl
    sudo install -m 0755 -d /etc/apt/keyrings
    sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
    sudo chmod a+r /etc/apt/keyrings/docker.asc

# Add the repository to Apt sources:
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
        $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
    sudo apt-get update
    sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
}

install_sysbox() {
    wget https://downloads.nestybox.com/sysbox/releases/v0.6.4/sysbox-ce_0.6.4-0.linux_amd64.deb
    apt-get install -y jq
    apt-get install -y ./sysbox-ce_0.6.4-0.linux_amd64.deb
}

print_header() {
    echo "############### $1 ###############"
}

print_env() {
    set +x
    print_header "Environment variables"
    printenv
    print_header "uname -a"
    uname -a || :
    print_header "Mounted file systems"
    cat /proc/self/mountinfo || :
    print_header "Kernel command line"
    cat /proc/cmdline || :
    print_header "Kernel modules"
    lsmod || :
    print_header "Distribution information"
    [ -e /etc/lsb-release ] && cat /etc/lsb-release
    [ -e /etc/redhat-release ] && cat /etc/redhat-release
    [ -e /etc/alpine-release ] && cat /etc/alpine-release
    print_header "ulimit -a"
    ulimit -a
    print_header "Available memory"
    if [ -e /etc/alpine-release ]; then
        # Alpine's busybox based free does not understand -h
        free
    else
        free -h
    fi
    print_header "Available CPUs"
    lscpu || :
    set -x
}

setup_ci() {
    [ -n "$SKIP_CI_SETUP" ] && return
    install_apt_packages
    install_code_server
    install_bats_core

    install_docker
    install_sysbox

    wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz && rm -rf /usr/local/go
    tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz && rm go1.22.0.linux-amd64.tar.gz
    echo '{"client":{"leave_running":false, "task":""}}' >~/.cedana/client_config.json
    # BRANCH_NAME="${CI_BRANCH:-main}"

    export PATH=$PATH:/usr/local/go/bin
    echo "export PATH=$PATH:/usr/local/go/bin" >>/root/.bashrc

    # Install CRIU
    sudo add-apt-repository -y ppa:criu/ppa
    sudo apt-get update && sudo apt-get install -y criu
    # # Install Cedana
    # git clone https://github.com/cedana/cedana && mkdir ~/.cedana
    # git fetch && git checkout ${BRANCH_NAME} && git pull origin ${BRANCH_NAME}

    # Install smoke & bench deps
    cd ../../
    sudo pip3 install -r test/benchmarks/requirements
}

start_cedana() {
    ./build-start-daemon.sh
}
