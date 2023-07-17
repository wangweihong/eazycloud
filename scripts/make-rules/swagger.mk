# ==============================================================================
# Makefile helper functions for swagger
#

.PHONY: swagger.run
swagger.run: tools.verify.swagger
	@echo "===========> Generating swagger API docs"
	#@swagger generate spec --scan-models -w $(ROOT_DIR)/cmd/genswaggertypedocs -o $(ROOT_DIR)/api/swagger/swagger.yaml
	@swag init --parseDependency --generalInfo ./cmd/gin-swagger-example/example.go --output ./api/swagger/example

.PHONY: swagger.serve
swagger.serve: tools.verify.swagger
	@swagger serve -F=redoc --no-open --port 36666 $(ROOT_DIR)/api/swagger/swagger.yaml


.PHONY: swagger.example
swagger.example: tools.verify.swagger
	@swag init --parseDependency --generalInfo ./cmd/gin-swagger-example/example.go --output ./api/swagger/example

.PHONY: swagger.example.serve
swagger.example.serve: tools.verify.swagger
	@swagger serve -F=redoc --no-open --port 36667 $(ROOT_DIR)/api/swagger/example/swagger.yaml