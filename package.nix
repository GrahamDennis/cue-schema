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
    vendorHash = "sha256-RbaaxFtUPp8oaK0XS3g07XWz2MQkx4uAi8tGFND2Lhk=";
}
