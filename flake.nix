{
  description =
    "A CLI for calculating the commute time from one address to one address or several other addresses. ";

  inputs = {
    flake-utils.url = "github:numtide/flake-utils";
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  };

  outputs = { nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let pkgs = nixpkgs.legacyPackages.${system};
      in {
        packages.default = pkgs.buildGoModule {
          pname = "traveltime";
          version = "0.0.1";
          src = pkgs.fetchFromGitHub {
            owner = "toffernator";
            repo = "traveltime";
            rev = "main";
            hash = "sha256-oywarRD+kYXWnqjtEC/IA6pCB8juUHpxVIj0NnjGYOY=";
          };

          # https://nixos.org/manual/nixpkgs/stable/#ex-buildGoModule
          vendorHash = "sha256-sAn0KozBfeYjjIvIdgSWJyfpN6x8uTLmdKrMDsn/6jA=";
        };

        devShells.default = pkgs.mkShell {
          packages = with pkgs; [ go google-cloud-sdk cobra-cli ];
        };
      });
}

