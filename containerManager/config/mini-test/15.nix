{ pkgs ? import <nixpkgs> { }
, pkgsLinux ? import <nixpkgs> { system = "x86_64-linux"; }
}:

pkgs.dockerTools.buildLayeredImage {
  name = "dig-google";
  config = {
    Cmd = [ "${pkgsLinux.coreutils}/bin/cat" "/etc/resolv.conf" ];

  };

  contents = with pkgsLinux; [ coreutils ];
}
