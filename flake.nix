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

        CGO_ENABLED = 0;
      };
    });

    packages = eachSystem (pkgs: rec {
      default = metastructs;

      metastructs = pkgs.buildGoModule (finalAttrs: {
        pname = "metastructs";
        version = "0.0.2";

        src = ./.;

        vendorHash = "sha256-JeMJffFR9ZcVHd0mJu0Xj3Ja45uDwgGoP18tWLiEfrg=";

        meta = {
          description = "Code generator for implementing boilerplate struct methods";
          homepage = "https:///github.com/kytnacode/metastructs";
          license = pkgs.lib.licenses.mit;
        };

        env.CGO_ENABLED = 0;
      });
    });
  };
}
