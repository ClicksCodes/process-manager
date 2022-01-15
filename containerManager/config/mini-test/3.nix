{ pkgs ? import <nixpkgs> { }
, pkgsLinux ? import <nixpkgs> { system = "x86_64-linux"; }
}:

pkgs.dockerTools.buildLayeredImage {
  name = "python-test-image";
  config = {
    Cmd = [ "${pkgsLinux.hello}/bin/hello" ];
  };

  contents = with pkgsLinux; [ hello ];
}
