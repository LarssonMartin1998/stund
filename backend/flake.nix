{
  description = "Go/Chi backend";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs { inherit system; };
        CGO = "1";
        go = pkgs.go;
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "backend";
          version = "0.0.1";
          src = ./.;

          vendorHash = "sha256-1OoBfO2KqtyTVFNLr+ftE5KIS50MyR71rO8zB/FFjVQ=";

          env = {
            CGO_ENABLED = CGO;
          };

          doCheck = true;
        };

        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go
            delve
            gotools
            air
          ];

          env = {
            CGO_ENABLED = CGO;
          };

          shellHook = ''
            echo "🦫 $(${go}/bin/go version) ready!"
            echo "CGO_ENABLED: $CGO_ENABLED"
          '';
        };
      }
    );
}
