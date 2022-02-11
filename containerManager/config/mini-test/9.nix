{ pkgs ? import <nixpkgs> { }
, pkgsLinux ? import <nixpkgs> { system = "x86_64-linux"; }
}:

pkgs.dockerTools.buildLayeredImage {
  name = "ping-cloudflare-dns";
  config = {
    Cmd = [ "${pkgsLinux.iputils}/bin/ping" "1.1.1.1" ];
  };

  contents = with pkgsLinux; [ iputils ];
}
