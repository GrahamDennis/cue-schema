{ pkgs }:
  pkgs.buildGoModule {
    name = "cue-schema";
    src = ./.;
    vendorHash = "sha256-F2RRkU1Nxp8euh1b4iDYbYgQRP7A/4wB0mYyimH4J20=";
}
