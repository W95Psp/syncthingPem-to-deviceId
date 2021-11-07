{
  description = "Extract Syncthing device ID from PEM files";

  inputs.flake-utils.url = "github:numtide/flake-utils";

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let pkgs = nixpkgs.legacyPackages.${system}; in
      rec {
        packages = flake-utils.lib.flattenTree {
          syncthing-pem2id = pkgs.stdenv.mkDerivation {
            name = "syncthing-pem2id";
            src = ./.;
            buildInputs = [ pkgs.go ];
            buildPhase = ''
              export HOME=$(realpath .)
              export GOPATH=$(realpath .)
              go build main.go
            '';
            installPhase = "mv main $out";
          };
        };
        defaultPackage = packages.syncthing-pem2id;
      }
    );
}
