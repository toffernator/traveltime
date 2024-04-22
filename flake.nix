{
  description =
    "A CLI for calculating the commute time from one address to one address or several other addresses. ";

  inputs = {
    flake-utils.url = "github:numtide/flake-utils";
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  };

  outputs = { nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        version = "0.0.1";
        pname = "traveltime";
      in {
        packages = {
          default = pkgs.buildGoModule {
            inherit pname version;

            src = pkgs.fetchFromGitHub {
              owner = "toffernator";
              repo = pname;
              rev = "main"; # TODO "v${pname}"
              hash = "sha256-ZRz0C4x2LFpIk1kQgu86p+cIeqXKnfozlBoaY/NRw6E=";
            };

            # https://nixos.org/manual/nixpkgs/stable/#ex-buildGoModule
            vendorHash = "sha256-sAn0KozBfeYjjIvIdgSWJyfpN6x8uTLmdKrMDsn/6jA=";
          };

          local = pkgs.buildGoModule {
            inherit pname version;
            src = ./.;
            vendorHash = "sha256-sAn0KozBfeYjjIvIdgSWJyfpN6x8uTLmdKrMDsn/6jA=";
          };
        };

        devShells.default = pkgs.mkShell {
          packages = with pkgs; [ go google-cloud-sdk cobra-cli ];
        };
      });
}

