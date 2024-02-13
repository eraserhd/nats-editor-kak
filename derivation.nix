{ buildGoModule
, fetchFromGitHub
}:

buildGoModule {
  pname = "kakoune-pluggo";
  version = "0.1.0";
  src = ./.;
  vendorHash = "sha256-JOpp6LvOz4jvMg3V+7u0SpiNgL4s6AgDR9rpinVS73I=";
}

