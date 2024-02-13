{ buildGoModule
, fetchFromGitHub
}:

buildGoModule {
  pname = "kakoune-pluggo";
  version = "0.1.0";
  src = ./.;
  vendorHash = "sha256-MbFvVloyn2BEu5zFc3kW+mUhtGIQTMT2MnDWzCaB5iY=";
}

