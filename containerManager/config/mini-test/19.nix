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
            sha256 = "sha256-wGvftnTv+79lfnPpKeOSIr44pCqHEW02XVOxGpnXqaM=";
        };

        buildPhase = "echo 'No build phase'";
        installPhase = ''
            mkdir $out/src -p
            cp $src/* $out/src -r
        '';
    };
    DISCORD_TOKEN = (import /home/minion/Private/create-machine-programmers-discord-token.nix {}).token;
in pkgs.dockerTools.buildImage {
  name = "discord-bot-runner";
  config = {
    Env = [
      "DISCORD_TOKEN=${DISCORD_TOKEN}"
      "PATH=${pkgsLinux.busybox}/bin:${pkgsLinux.nodejs-17_x}/bin"
    ];
    Entrypoint = [ "${pkgsLinux.nodejs-17_x}/bin/npm" ];
    Cmd = [ "run" "container" ];
    WorkingDir = "${repo}/src";

  };

  runAsRoot = ''
    #!${pkgsLinux.busybox}/bin/sh
    mkdir /etc
    echo "$(hostname)" > /etc/hostname
    echo "127.0.0.1 localhost" >> /etc/hosts
    echo "127.0.0.2 $(cat /etc/hostname)" >> /etc/hosts
    echo "nameserver 1.1.1.1" >> /etc/resolv.conf
    echo "nameserver 1.0.0.1" >> /etc/resolv.conf
  '';

  contents = [ pkgsLinux.busybox repo ];
}
