{ pkgs ? import <nixpkgs> { }
, pkgsLinux ? import <nixpkgs> { system = "x86_64-linux"; }
}:

pkgs.dockerTools.buildLayeredImage {
  name = "ip-route";
  config = {
    Cmd = [ "${pkgsLinux.iproute2}/bin/ip" "route" ];
    services.resolved.enable = true;
  };

  contents = with pkgsLinux; [ iproute2 ];
}
