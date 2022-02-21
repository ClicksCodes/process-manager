{ pkgs ? import <nixpkgs> { }
, pkgsLinux ? import <nixpkgs> { system = "x86_64-linux"; }
}:

pkgs.dockerTools.buildLayeredImage {
  name = "dig-google";
  config = {
    Cmd = [ "${pkgsLinux.dig}/bin/dig" "@1.1.1.1" "google.com" ];
  };

  contents = with pkgsLinux; [ dig ];
}
