{ pkgs ? import <nixpkgs> { } }:

pkgs.mkShell {
  buildInputs = [
    pkgs.jq
    pkgs.nixpkgs-fmt

    pkgs.go_1_18
  ];
}
