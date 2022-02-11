{ pkgs ? import <nixpkgs> { }
, pkgsLinux ? import <nixpkgs> { system = "x86_64-linux"; }
}:

pkgs.dockerTools.buildLayeredImage {
  name = "ip-a";
  config = {
    Cmd = [ "${pkgsLinux.iproute2}/bin/ip" "a" ];
  };

  contents = with pkgsLinux; [ iproute2 ];
}
