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
    ${pkgsLinux.iputils}/bin/ping -c 3 1.1.1.1 &&
    ${pkgsLinux.iputils}/bin/ping -c 3 google.com
    ''];
  };

  contents = with pkgsLinux; [ iputils busybox ];
}