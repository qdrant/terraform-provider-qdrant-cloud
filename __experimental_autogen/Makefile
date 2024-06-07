TFPLUGINGEN := ~/go/bin/tfplugingen-framework
TFPLUGINGEN_OPENAPI := ~/go/bin/tfplugingen-openapi
OUTPUT_DIR := internal/provider

cleanup:
	rm -rf $(OUTPUT_DIR)
	rm -rf provider-code-spec.json
	mkdir -p $(OUTPUT_DIR)

scaffold:
	$(TFPLUGINGEN) scaffold resource --name resource_auth_api_keys --output-dir $(OUTPUT_DIR)
	$(TFPLUGINGEN) scaffold resource --name resource_accoun_cluster --output-dir $(OUTPUT_DIR)
	$(TFPLUGINGEN) scaffold data-source --name data_source_accoun_cluster --output-dir $(OUTPUT_DIR)
	$(TFPLUGINGEN) scaffold data-source --name data_source_auth_api_keys --output-dir $(OUTPUT_DIR)
	$(TFPLUGINGEN) scaffold data-source --name data_source_package --output-dir $(OUTPUT_DIR)

gen-spec:
	$(TFPLUGINGEN_OPENAPI) generate  --config ./generator_config.yml --output ./provider-code-spec.json ./openapi_fixed.json
	patch provider-code-spec.json < provider-code-spec.patch

all: cleanup gen-spec scaffold
	$(TFPLUGINGEN) generate resources --input ./provider-code-spec.json --output ./$(OUTPUT_DIR)