let
  pkgs = import <nixpkgs> { };
  nur = import (builtins.fetchTarball
    "https://github.com/nix-community/NUR/archive/master.tar.gz") {
      inherit pkgs;
    };
in pkgs.mkShell {
  buildInputs = with pkgs; [ go goimports golint nur.repos.xe.gopls ];
  LN_FORMATTER = "text";
}
