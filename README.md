# PRISM


.PHONY: docs run dev

docs:
	swag init --dir cmd/server,handlers,dto --generalInfo main.go

run:
	go run cmd/server/main.go

dev: docs run