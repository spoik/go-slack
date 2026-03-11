{
  description = "Flake for remote server configuration to serve the chat app.";

  inputs.nixpkgs.url = "github:nixos/nixpkgs/nixos-25.11";

  outputs = { self, nixpkgs }: {
    # colmena = import ./colmena.nix { inherit nixpkgs self; };
    colmena = {
      meta = {
        nixpkgs = import nixpkgs { system = "x86_64-linux"; };
        # Allows sharing arguments across all nodes
        specialArgs = { inherit self; };
      };

      "web-server" = { name, nodes, ... }: {
        deployment = {
          targetHost = "main";
          targetUser = "spike";
          buildOnTarget = true;
        };

        imports = [ ./hosts/web_sever/configuration.nix ];
        nixpkgs.hostPlatform = "aarch64-linux";
      };
    };
    
    # Optional: Keep standard nixosConfigurations so nixos-rebuild still works
    # nixosConfigurations = self.colmena;

    # nixosConfigurations.web-server = nixpkgs.lib.nixosSystem {
    #   system = "aarch64-linux";
    #   modules = [ ./hosts/web_server/configuration.nix ];
    # };
  };
}
