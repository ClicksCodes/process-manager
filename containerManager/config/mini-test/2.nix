{ pkgs ? import <nixpkgs> { }
, pkgsLinux ? import <nixpkgs> { system = "x86_64-linux"; }
}:

pkgs.dockerTools.buildImage {
  name = "python-test-image";
  config = {
    Cmd = [ "${pkgsLinux.python3}/bin/python3" ];
  };
  contents = with pkgs; [
    python3
  ];
}
