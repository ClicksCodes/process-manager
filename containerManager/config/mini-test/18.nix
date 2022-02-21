{ pkgs ? import <nixpkgs> { }
, pkgsLinux ? import <nixpkgs> { system = "x86_64-linux"; }
}:

pkgs.dockerTools.buildLayeredImage {
  name = "ping-1.1.1.1-and-google";
  config = {
    Cmd = [
    "sh"
    "-c"
    ''
    echo "$(hostname)" > /etc/hostname &&
    echo "127.0.0.1 localhost" >> /etc/hosts &&
    echo "127.0.0.2 $(cat /etc/hostname)" >> /etc/hosts &&
    echo "nameserver 1.1.1.1" >> /etc/resolv.conf &&
    echo "nameserver 1.0.0.1" >> /etc/resolv.conf &&
    ${pkgsLinux.dig}/bin/dig @1.1.1.1 google.com &&
    ${pkgsLinux.iputils}/bin/ping -c 3 1.1.1.1 &&
    ${pkgsLinux.iputils}/bin/ping -c 3 google.com
    ''];
  };

  contents = with pkgsLinux; [ iputils busybox dig ];
}