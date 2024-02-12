{
  description = "TODO: fill me in";
  inputs = {
    flake-utils.url = "github:numtide/flake-utils";
  };
  outputs = { self, nixpkgs, flake-utils }:
    (flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        kakoune-pluggo = pkgs.callPackage ./derivation.nix {};
      in {
        packages = {
          default = kakoune-pluggo;
          inherit kakoune-pluggo;
        };
        checks = {
          test = pkgs.runCommand "kakoune-pluggo-test" {} ''
            mkdir -p $out
            : ${kakoune-pluggo}
          '';
        };
        devShells.default = pkgs.mkShell {
          nativeBuildInputs = with pkgs; [
            go_1_21
          ];
        };
    })) // {
      overlays.default = final: prev: {
        kakoune-pluggo = prev.callPackage ./derivation.nix {};
      };
    };
}
