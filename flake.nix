{
  description = "A very basic Go flake";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
  };

  outputs = {
    self,
    nixpkgs,
  }: let
    systems = ["x86_64-linux"];
    eachSystem = f: nixpkgs.lib.genAttrs systems (system: f (import nixpkgs {inherit system;}));
  in {
    devShells = eachSystem (pkgs: {
      default = pkgs.mkShellNoCC {
        buildInputs = with pkgs; [
          git
          just

          # Go packages
          go
          golangci-lint
          govulncheck
        ];

        CGO_ENABLE = 0;
      };
    });
  };
}
