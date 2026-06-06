{
  pkgs ? import <nixpkgs> { },
}:

pkgs.mkShell {
  name = "yumeko-dev";

  packages = with pkgs; [
    go
    gopls
    gotools
    gofumpt
    golangci-lint

    air

    sqlite
    git
  ];

  shellHook = ''
    echo "Yumeko dev shell"
    echo "Go: $(go version)"
    echo ""
    echo "Commands:"
    echo "  air              - run with hot reload"
    echo "  go run ./cmd/yumeko"
    echo "  go test ./..."
    echo "  gofmt -w ."
  '';
}
