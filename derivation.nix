{ buildGo119Module
, fetchFromGitHub
}:

buildGo119Module {
  pname = "nats-editor-kak";
  version = "0.1.0";
  src = ./.;
  vendorSha256 = "gDIFI6bE2q23yJifzXpeyiN6C4QFTKkFOBjrO9Bohwg=";
}

