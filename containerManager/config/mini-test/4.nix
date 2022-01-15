{ pkgs ? import <nixpkgs> { }
, pkgsLinux ? import <nixpkgs> { system = "x86_64-linux"; }
}:
let
    repo = pkgsLinux.stdenv.mkDerivation {
        name = "repo";
        src = pkgs.fetchFromGitHub {
            owner = "Minion3665";
            repo = "container";
            rev = "production";
            sha256 = "sha256-3C4BFcn1fOByG1YEzTXWJjo2C4llO3NxN7teLxOEyYA=";
        };

        buildPhase = "echo 'No build phase'";
        installPhase = ''
            mkdir $out/src -p
            cp $src/* $out/src -r
        '';
    };
in pkgs.dockerTools.buildLayeredImage {
  name = "python-test-image";
  config = {
    Cmd = [ "${pkgsLinux.python3}/bin/python3 ${repo}/src/main.py" ];
  };

  contents = with pkgsLinux; [ python3 repo ];
}
