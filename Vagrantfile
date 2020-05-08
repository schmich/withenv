require 'tmpdir'

Vagrant.configure('2') do |config|
  config.vm.box = 'generic/centos7'
  config.vm.hostname = 'withenv'

  config.vm.provision 'shell', inline: <<-SCRIPT
    sudo yum install -y git
    curl -Ls https://dl.google.com/go/go1.14.2.linux-amd64.tar.gz | sudo tar -C /usr/local -xz

    echo >> ~vagrant/.bashrc
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~vagrant/.bashrc
    echo 'cd /withenv' >> ~vagrant/.bashrc
SCRIPT

  config.vm.provider 'virtualbox' do |vb|
    vb.memory = '2048'

    # Move VM host log output to temp dir.
    vb.customize ['modifyvm', :id, '--uartmode1', 'file', File.join(Dir.tmpdir, 'vb')]
  end

  config.vm.synced_folder '.', '/withenv', mount_options: ['dmode=775,fmode=664']

  config.vm.network :public_network
  config.vm.network :private_network, ip: '10.0.0.20'
end
