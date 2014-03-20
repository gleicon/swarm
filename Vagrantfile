# -*- mode: ruby -*-
# vi: set ft=ruby :

# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
#
# Place this Vagrantfile in your src folder and run:
#
#     vagrant up
#
# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

# Vagrantfile API/syntax version.
VAGRANTFILE_API_VERSION = "2"

# See https://code.google.com/p/go/downloads/list
GO_ARCHIVES = {
  "linux" => "go1.2.linux-amd64.tar.gz",
}

INSTALL = {
  "linux" => "apt-get update -qq; apt-get install -qq -y git mercurial bzr curl",
}

# location of the Vagrantfile
def src_path
  File.dirname(__FILE__)
end

# shell script to bootstrap Go
def bootstrap(box)
  install = INSTALL[box]
  archive = GO_ARCHIVES[box]

  profile = <<-PROFILE
    export GOPATH=$HOME
    export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
    export CDPATH=.:$GOPATH/src/github.com:$GOPATH/src/code.google.com/p:$GOPATH/src/bitbucket.org:$GOPATH/src/launchpad.net
  PROFILE

  <<-SCRIPT
  #{install}

  if ! [ -f /home/vagrant/#{archive} ]; then
    response=$(curl -O# https://go.googlecode.com/files/#{archive})
  fi
  tar -C /usr/local -xzf #{archive}

  echo '#{profile}' >> /home/vagrant/.profile

  echo "\nRun: vagrant ssh #{box} -c 'cd project/path; go test ./...'"
  SCRIPT
end

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|

  config.vm.define "linux" do |linux|
    linux.vm.box = "precise64"
    linux.vm.box_url = "http://files.vagrantup.com/precise64.box"
    linux.vm.provision :shell, :inline => bootstrap("linux")
  end

end
