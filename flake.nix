{
  description = "GamesDB mirror in go";
  nixConfig.bash-prompt = "flake:gamesdb-mirror $ ";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let pkgs = nixpkgs.legacyPackages.${system}; in
      rec {
        devShell = import ./shell.nix { inherit pkgs; };
      }
    );
}
