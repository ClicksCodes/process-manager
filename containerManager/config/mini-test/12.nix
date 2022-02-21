{ pkgs ? import <nixpkgs> { }
, pkgsLinux ? import <nixpkgs> { system = "x86_64-linux"; }
}:

pkgs.dockerTools.buildLayeredImage {
  name = "ping-google";
  config = {
    Cmd = [ "${pkgsLinux.iputils}/bin/ping" "google.com" ];
  };

  contents = with pkgsLinux; [ iputils ];
}
