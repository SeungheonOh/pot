{ nixpkgs ? import <nixpkgs> { }}:
let 
  pkgs = with nixpkgs; [
    pkgconfig
    go
    xorg.libX11
  ];

in 
  nixpkgs.stdenv.mkDerivation{
    name = "PixelOnTerminal";
    buildInputs = pkgs;
    meta = {
      description = "Pixels but in terminal";
      homepage = https://github.com/SeungheonOh/PixelOnTerminal/;
      license = nixpkgs.stdenv.lib.licenses.mit;
      maintainers = [ "Seungheon Oh - dan.oh0721@gmail.com" ];
    };
  }
