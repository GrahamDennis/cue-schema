{ pkgs, lib }:
  pkgs.buildGoModule {
    name = "cue-schema";
    src = lib.fileset.toSource {
      root = ./.;
      fileset = lib.fileset.unions [
        ./cmd
        ./go.mod
        ./go.sum
        ./main.go
      ];
    };
    vendorHash = "sha256-F2RRkU1Nxp8euh1b4iDYbYgQRP7A/4wB0mYyimH4J20=";
}
