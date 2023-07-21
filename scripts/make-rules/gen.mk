# ==============================================================================
# Makefile helper functions for generate necessary files
#

.PHONY: gen.run
gen.run: gen.clean gen.errcode
# gen.run: gen.clean gen.errcode gen.docgo.doc

.PHONY: gen.errcode
gen.errcode: gen.errcode.code gen.errcode.doc

.PHONY: gen.errcode.code
gen.errcode.code: tools.verify.codegen
	@echo "===========> Generating error code go source files to path:${ROOT_DIR}/internal/pkg/code"
	@codegen -type=int ${ROOT_DIR}/internal/pkg/code

.PHONY: gen.errcode.doc
gen.errcode.doc: tools.verify.codegen
	@echo "===========> Generating error code markdown documentation:${ROOT_DIR}/docs/guide/zh-CN/api/error_code_generated.md"
	@codegen -type=int -doc \
		-output ${ROOT_DIR}/docs/guide/zh-CN/api/error_code_generated.md ${ROOT_DIR}/internal/pkg/code

.PHONY: gen.docgo.doc
gen.docgo.doc:
	@echo "===========> Generating missing doc.go for go packages"
	@${ROOT_DIR}/scripts/gendoc.sh

.PHONY: gen.docgo.check
gen.docgo.check: gen.docgo.doc
	@n="$$(git ls-files --others '*/doc.go' | wc -l)"; \
	if test "$$n" -gt 0; then \
		git ls-files --others '*/doc.go' | sed -e 's/^/  /'; \
		echo "$@: untracked doc.go file(s) exist in working directory" >&2 ; \
		false ; \
	fi

# 根据configs/*.yaml中的模板,
.PHONY: gen.defaultconfigs
gen.defaultconfigs:
	@${ROOT_DIR}/scripts/gen_default_config.sh ${COMPONENTS}

.PHONY: gen.ca.%
gen.ca.%:
	$(eval Certifcate := $(word 1,$(subst ., ,$*)))
	@echo "===========> Generating Certifcate files for $(Certifcate),Subjects:$(CERTIFICATES_SUBJECT)"
	@echo "===========> OUTPUT_DIR:$(OUTPUT_DIR)/cert"
	@${ROOT_DIR}/scripts/gencerts.sh generate_certificate $(OUTPUT_DIR)/cert $(Certifcate) $(CERTIFICATES_SUBJECT)

.PHONY: gen.ca
gen.ca: $(addprefix gen.ca., $(CERTIFICATES))

.PHONY: gen.clean
gen.clean:
	@$(FIND) -path ${ROOT_DIR}/internal/pkg/code -type f -name '*_generated.go' -delete