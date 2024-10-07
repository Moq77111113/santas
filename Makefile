
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
	$(MAKE) ui-build
	go run cmd/chmoly/chmoly.go serve

.PHONY: serve-only
serve-only:
	go run cmd/chmoly/chmoly.go serve

.PHONY: ui-install
ui-install:
	cd ui && pnpm install && cd ..

.PHONY: ui-build
ui-build:
	cd ui && pnpm run build && cd ..
