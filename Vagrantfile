# -*- mode: ruby -*-
# vi: set ft=ruby :

# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  # All Vagrant configuration is done here. The most common configuration
  # options are documented and commented below. For a complete reference,
  # please see the online documentation at vagrantup.com.

  # Every Vagrant virtual environment requires a box to build off of.
  config.vm.box = "ubuntu/trusty64"

  config.vm.hostname = "boxr-master"

  # Create a forwarded port mapping which allows access to a specific port
  # within the machine from a port on the host machine. In the example below,
  # accessing "localhost:8080" will access port 80 on the guest machine.
  # config.vm.network "forwarded_port", guest: 3000, host: 8080

  # Create a private network, which allows host-only access to the machine
  # using a specific IP.
  config.vm.network "private_network", ip: "10.10.10.10"

  # Create a public network, which generally matched to bridged network.
  # Bridged networks make the machine appear as another physical device on
  # your network.
  # config.vm.network "public_network"

  # If true, then any SSH connections made will enable agent forwarding.
  # Default value: false
  # config.ssh.forward_agent = true

  # Share an additional folder to the guest VM. The first argument is
  # the path on the host to the actual folder. The second argument is
  # the path on the guest to mount the folder. And the optional third
  # argument is a set of non-required options.
  config.vm.synced_folder ".", "/opt/go/src/github.com/Secret-Ironman/boxr"

  config.vm.provision "shell", inline: <<-EOF
    set -e

    # System packages
    echo "Installing Base Packages"
    export DEBIAN_FRONTEND=noninteractive
    sudo apt-get update -qq
    sudo apt-get install -qqy --force-yes build-essential bzr git mercurial vim

    # Install Go
    GOVERSION="1.3.3"
    GOTARBALL="go${GOVERSION}.linux-amd64.tar.gz"
    export GOROOT=/usr/local/go
    export GOPATH=/opt/go
    export PATH=$PATH:$GOROOT/bin:$GOPATH/bin

    echo "Installing Go $GOVERSION"
    if [ ! $(which go) ]; then
        echo "    Downloading $GOTARBALL"
        wget --quiet --directory-prefix=/tmp https://storage.googleapis.com/golang/$GOTARBALL

        echo "    Extracting $GOTARBALL to $GOROOT"
        sudo tar -C /usr/local -xzf /tmp/$GOTARBALL

        echo "    Configuring GOPATH"
        sudo mkdir -p $GOPATH/src $GOPATH/bin $GOPATH/pkg
        sudo chown -R vagrant $GOPATH

        echo "    Configuring env vars"
        echo "export PATH=\$PATH:$GOROOT/bin:$GOPATH/bin" | sudo tee /etc/profile.d/golang.sh > /dev/null
        echo "export GOROOT=$GOROOT" | sudo tee --append /etc/profile.d/golang.sh > /dev/null
        echo "export GOPATH=$GOPATH" | sudo tee --append /etc/profile.d/golang.sh > /dev/null
    fi

    # Auto cd to boxr install dir
    echo "cd /opt/go/src/github.com/Secret-Ironman/boxr\nexport BOXR_PORT=3000" >> /home/vagrant/.bashrc

    # add user to docker
    groupadd docker
    gpasswd -a vagrant docker

    # Cleanup
    sudo apt-get autoremove -y --force-yes

    echo <<DONE
    PROVISIONING COMPLETE:
        vagrant ssh
    DONE
  EOF

  config.vm.provision "shell", inline: <<-EOF
    sudo apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 36A1D7869245C8950F966E92D8576A8BA88D21E9
    sudo sh -c 'echo deb http://get.docker.io/ubuntu docker main' > /etc/apt/sources.list.d/docker.list
    sudo apt-get update
    sudo apt-get install -y linux-image-generic-lts-raring linux-headers-generic-lts-raring
    sudo apt-get install -y lxc-docker-1.0.1
  EOF
end
