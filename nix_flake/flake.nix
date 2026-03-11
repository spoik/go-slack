{
  description = "Flake for remote server configuration to serve the chat app.";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-25.11";
    colmena.url = "github:zhaofengli/colmena";
  };

  outputs = { self, nixpkgs, colmena, ... }: {
    colmenaHive = colmena.lib.makeHive {
      meta = {
        nixpkgs = import nixpkgs {
          system = "x86_64-linux";
        };
      };

      "web-server" = { name, nodes, ... }: {
        deployment = {
          targetHost = "main";
          targetUser = "spike";
          buildOnTarget = true;
        };

        imports = [ ./hosts/web_server/configuration.nix ];
        nixpkgs.hostPlatform = "aarch64-linux";
      };
    };
  };
}
