{ buildGo121Module
, fetchFromGitHub
}:

buildGo121Module {
  pname = "kakoune-pluggo";
  version = "0.1.0";
  src = ./.;
  vendorHash = "sha256:gDIFI6bE2q23yJifzXpeyiN6C4QFTKkFOBjrO9Bohwg=";
}

