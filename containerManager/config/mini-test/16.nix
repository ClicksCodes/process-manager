{ pkgs ? import <nixpkgs> { }
, pkgsLinux ? import <nixpkgs> { system = "x86_64-linux"; }
}:

pkgs.ociTools.buildContainer {
    args = [
        (with pkgs;
            writeScript "run.sh" ''
                #!${bash}/bin/bash
                exec ${bash}/bin/bash
            '').outPath
    ];

    mounts = {
        "/data" = {
            type = "none";
            source = "/var/lib/mydata";
            options = [ "bind" ];
        };
    };

    readonly = false;
}