{ pkgs ? import <nixpkgs> { }
, pkgsLinux ? import <nixpkgs> { system = "x86_64-linux"; }
}:
pkgs.dockerTools.buildImage {
  name = "discord-bot-runner";
  config = {
    Cmd = [ "sh" ];
    WorkingDir = "/root";
  };

  runAsRoot = ''
    #!${pkgsLinux.busybox}/bin/sh
    mkdir /etc
    echo "$(hostname)" > /etc/hostname
    echo "127.0.0.1 localhost" >> /etc/hosts
    echo "127.0.0.2 $(cat /etc/hostname)" >> /etc/hosts
    echo "nameserver 1.1.1.1" >> /etc/resolv.conf
    echo "nameserver 1.0.0.1" >> /etc/resolv.conf
  '';

  contents = [ pkgsLinux.busybox ];
}
