
.PHONY: install-orm
orm-install:
	go get entgo.io/ent/cmd/ent

# Generate orm code
.PHONY: orm-gen
orm-gen:
	go generate ./ent

# Run app
.PHONY: serve
serve:
	go run cmd/chmoly/chmoly.go serve