((nil
  (eval progn
        (setq-local
         compilation-read-command nil
         projectile-project-compilation-command "make build"
         projectile-project-test-cmd "go test -race ./...")))
 (go-mode (eval progn (setq-local lsp-go-use-gofumpt t))))
