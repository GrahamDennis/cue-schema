{
  description = "cpp-package";

  inputs = {
    flake-utils.url = "github:numtide/flake-utils";
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-24.11";
  };

  outputs =
    inputs@{
      flake-utils,
      nixpkgs,
      ...
    }:
    let
      perSystemOutputs = flake-utils.lib.eachDefaultSystem (
        system:
        let
          pkgs = import nixpkgs {
            inherit system;
          };
          cue-schema = pkgs.callPackage ./package.nix {};
        in
        {
          packages = {
            inherit cue-schema;
            default = cue-schema;
          };
        }
      );
    in
    perSystemOutputs
    // {
      inherit inputs;
    };
}
