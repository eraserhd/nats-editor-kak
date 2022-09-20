{ buildGo119Module
, fetchFromGitHub
}:

buildGo119Module {
  pname = "nats-editor-kak";
  version = "0.1.0";
  src = ./.;
  vendorSha256 = "BzVGuisFW2iJWttsx9OUcHvLOmazEYCw+SDORYZmS6o=";
}

