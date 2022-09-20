{
  description = "TODO: fill me in";
  inputs = {
    flake-utils.url = "github:numtide/flake-utils";
  };
  outputs = { self, nixpkgs, flake-utils }:
    (flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        nats-editor-kak = pkgs.callPackage ./derivation.nix {};
      in {
        packages = {
          default = nats-editor-kak;
          inherit nats-editor-kak;
        };
        checks = {
          test = pkgs.runCommandNoCC "nats-editor-kak-test" {} ''
            mkdir -p $out
            : ${nats-editor-kak}
          '';
        };
    })) // {
      overlays.default = final: prev: {
        nats-editor-kak = prev.callPackage ./derivation.nix {};
      };
    };
}
