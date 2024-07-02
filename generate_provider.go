package main

import (
	"encoding/json"
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/getkin/kin-openapi/openapi3"
	"log"
	"os"
	"strings"
)

func main() {
	spec, err := loadOpenAPISpec("./internal/spec.json")
	if err != nil {
		fmt.Printf("Error loading spec: %v\n", err)
		os.Exit(1)
	}

	err = generateProvider(spec)
	if err != nil {
		log.Fatalf("Error generating provider: %v", err)
		return
	}

	err = generateResources(spec)
	if err != nil {
		log.Fatalf("Error generating resources: %v", err)
		return
	}

	generateProviderDocumentation(spec)
}

func loadOpenAPISpec(filename string) (*openapi3.T, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var spec openapi3.T
	if err := json.Unmarshal(data, &spec); err != nil {
		return nil, err
	}

	return &spec, nil
}

func generateProvider(spec *openapi3.T) error {
	f := jen.NewFile("qdrant")

	// Add init function
	f.Func().Id("init").Params().Block(
		jen.Qual("github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema", "DescriptionKind").Op("=").Qual("github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema", "StringMarkdown"),
	)

	// Generate Provider function
	f.Func().Id("Provider").Params().Op("*").Qual("github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema", "Provider").Block(
		jen.Return(jen.Op("&").Qual("github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema", "Provider").Values(jen.Dict{
			jen.Id("Schema"): jen.Map(jen.String()).Op("*").Qual("github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema", "Schema").Values(jen.Dict{
				jen.Lit("api_key"): jen.Op("&").Qual("github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema", "Schema").Values(jen.Dict{
					jen.Id("Type"):        jen.Qual("github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema", "TypeString"),
					jen.Id("Required"):    jen.True(),
					jen.Id("DefaultFunc"): jen.Qual("github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema", "EnvDefaultFunc").Call(jen.Lit("QDRANT_CLOUD_API_KEY"), jen.Nil()),
					jen.Id("Description"): jen.Lit("The API Key for Qdrant Cloud API operations."),
				}),
				jen.Lit("api_url"): jen.Op("&").Qual("github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema", "Schema").Values(jen.Dict{
					jen.Id("Type"):        jen.Qual("github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema", "TypeString"),
					jen.Id("Optional"):    jen.True(),
					jen.Id("DefaultFunc"): jen.Qual("github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema", "EnvDefaultFunc").Call(jen.Lit("QDRANT_CLOUD_API_URL"), jen.Lit("https://cloud.qdrant.io/public/v0")),
					jen.Id("Description"): jen.Lit("The URL of the Qdrant Cloud API."),
				}),
				jen.Lit("account_id"): jen.Op("&").Qual("github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema", "Schema").Values(jen.Dict{
					jen.Id("Type"):        jen.Qual("github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema", "TypeString"),
					jen.Id("Optional"):    jen.True(),
					jen.Id("DefaultFunc"): jen.Qual("github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema", "EnvDefaultFunc").Call(jen.Lit("QDRANT_CLOUD_ACCOUNT_ID"), jen.Lit("")),
					jen.Id("Description"): jen.Lit("Default Account Identifier for the Qdrant cloud"),
				}),
			}),
			jen.Id("ResourcesMap"):         jen.Map(jen.String()).Op("*").Qual("github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema", "Resource").Values(generateResourcesMap(spec)),
			jen.Id("DataSourcesMap"):       jen.Map(jen.String()).Op("*").Qual("github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema", "Resource").Values(generateDataSourcesMap(spec)),
			jen.Id("ConfigureContextFunc"): jen.Id("providerConfigure"),
		})),
	)

	// Save generated code to file
	if err := f.Save("qdrant/provider.go"); err != nil {
		return fmt.Errorf("error saving generated provider code: %v", err)
	}

	return nil
}

func generateResourcesMap(spec *openapi3.T) jen.Dict {
	resources := jen.Dict{}
	for path, pathItem := range spec.Paths {
		resourceName := extractResourceName(path)
		sanitizedResourceName := strings.ReplaceAll(resourceName, "-", "_")
		if pathItem.Post != nil || pathItem.Put != nil {
			resources[jen.Lit(fmt.Sprintf("qdrant-cloud_%s", strings.ToLower(sanitizedResourceName)))] = jen.Id(fmt.Sprintf("resource%s", sanitizedResourceName)).Call()
		}
	}
	return resources
}

func generateDataSourcesMap(spec *openapi3.T) jen.Dict {
	dataSources := jen.Dict{}
	for path, pathItem := range spec.Paths {
		resourceName := extractResourceName(path)
		sanitizedResourceName := strings.ReplaceAll(resourceName, "-", "_")
		if pathItem.Get != nil {
			dataSources[jen.Lit(fmt.Sprintf("qdrant-cloud_%s", strings.ToLower(sanitizedResourceName)))] = jen.Id(fmt.Sprintf("dataSource%s", sanitizedResourceName)).Call()
		}
	}
	return dataSources
}

func extractResourceName(path string) string {
	parts := strings.Split(path, "/")
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] != "" && parts[i] != "{" && !strings.Contains(parts[i], "}") {
			if parts[i] == "api-keys" {
				return "ApiKey"
			}
			name := strings.ReplaceAll(parts[i], "-", "_")
			name = strings.TrimSuffix(name, "s")
			return strings.Title(name)
		}
	}
	return "Resource"
}

func sanitizeResourceName(name string) string {
	return strings.ReplaceAll(name, "-", "_")
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func extractResources(spec *openapi3.T) map[string]*openapi3.PathItem {
	resources := make(map[string]*openapi3.PathItem)
	for path, pathItem := range spec.Paths {
		resourceName := extractResourceName(path)
		resources[resourceName] = pathItem
	}
	return resources
}

func generateProviderDocumentation(spec *openapi3.T) {
	var doc strings.Builder

	doc.WriteString("# Qdrant Cloud Terraform Provider\n\n")
	doc.WriteString("This provider allows you to interact with Qdrant Cloud resources.\n\n")

	doc.WriteString("## Provider Configuration\n\n")
	doc.WriteString("```hcl\n")
	doc.WriteString("provider \"qdrant-cloud\" {\n")
	doc.WriteString("  api_key = var.qdrant_cloud_api_key\n")
	doc.WriteString("}\n")
	doc.WriteString("```\n\n")

	resources := []string{"AccountsAuthKey"} // Add more resources here as needed
	for _, resourceName := range resources {
		generateResourceDocumentation(&doc, spec, resourceName)
	}

	err := os.WriteFile("docs/provider.md", []byte(doc.String()), 0644)
	if err != nil {
		fmt.Printf("Error saving documentation: %v\n", err)
		os.Exit(1)
	}
}

func generateResourceDocumentation(doc *strings.Builder, spec *openapi3.T, resourceName string) {
	doc.WriteString(fmt.Sprintf("## Resource: qdrant-cloud_%s\n\n", strings.ToLower(resourceName)))
	doc.WriteString(fmt.Sprintf("Manages a %s for a Qdrant Cloud account.\n\n", resourceName))
	doc.WriteString("### Example Usage\n\n")
	doc.WriteString("```hcl\n")
	doc.WriteString(fmt.Sprintf("resource \"qdrant-cloud_%s\" \"example\" {\n", strings.ToLower(resourceName)))

	schemaRef := spec.Components.Schemas[resourceName]
	if schemaRef != nil && schemaRef.Value != nil {
		for propName := range schemaRef.Value.Properties {
			doc.WriteString(fmt.Sprintf("  %s = \"example_%s\"\n", propName, propName))
		}
	} else {
		doc.WriteString("  # Schema information not available\n")
	}

	doc.WriteString("}\n")
	doc.WriteString("```\n\n")

	doc.WriteString("### Argument Reference\n\n")
	if schemaRef != nil && schemaRef.Value != nil {
		for propName, propSchema := range schemaRef.Value.Properties {
			doc.WriteString(fmt.Sprintf("* `%s` - ", propName))
			if propSchema.Value != nil {
				doc.WriteString(propSchema.Value.Description + " ")
				doc.WriteString(fmt.Sprintf("(Type: %s) ", propSchema.Value.Type))
				if contains(schemaRef.Value.Required, propName) {
					doc.WriteString("(Required) ")
				}
				if propSchema.Value.Default != nil {
					doc.WriteString(fmt.Sprintf("(Default: %v) ", propSchema.Value.Default))
				}
			}
			doc.WriteString("\n")
		}
	} else {
		doc.WriteString("* Schema information not available\n")
	}
	doc.WriteString("\n")
}

func generateResources(spec *openapi3.T) error {
	for path, pathItem := range spec.Paths {
		resourceName := extractResourceName(path)
		var methods []string
		if pathItem.Get != nil {
			methods = append(methods, "GET")
		}
		if pathItem.Post != nil {
			methods = append(methods, "POST")
		}
		if pathItem.Put != nil {
			methods = append(methods, "PUT")
		}
		if pathItem.Delete != nil {
			methods = append(methods, "DELETE")
		}
		if err := generateCRUDFunction(spec, resourceName, methods, path); err != nil {
			return fmt.Errorf("error generating CRUD function for resource %s: %v", resourceName, err)
		}
	}
	return nil
}

func generateCRUDFunction(spec *openapi3.T, resourceName string, methods []string, path string) error {
	for _, method := range methods {
		op := findOperation(spec, path, method)
		if op != nil {
			switch method {
			case "GET":
				if err := generateReadFunction(resourceName, op); err != nil {
					return fmt.Errorf("error generating read function for resource %s: %v", resourceName, err)
				}
				//case "POST":
				//	if err := generateCreateFunction(resourceName, op); err != nil {
				//		return fmt.Errorf("error generating create function for resource %s: %v", resourceName, err)
				//	}
				//case "PUT":
				//	if err := generateUpdateFunction(resourceName, op); err != nil {
				//		return fmt.Errorf("error generating update function for resource %s: %v", resourceName, err)
				//	}
				//case "DELETE":
				//	if err := generateDeleteFunction(resourceName, op); err != nil {
				//		return fmt.Errorf("error generating delete function for resource %s: %v", resourceName, err)
				//	}
			}
		}
	}
	return nil
}

func generateDeleteFunction(resourceName string, op *openapi3.Operation) error {
	return nil
}

func generateUpdateFunction(resourceName string, op *openapi3.Operation) error {
	return nil
}

func generateCreateFunction(resourceName string, op *openapi3.Operation) error {
	return nil
}

func findOperation(spec *openapi3.T, path, method string) *openapi3.Operation {
	pathItem := spec.Paths[path]
	if pathItem == nil {
		return nil
	}
	switch method {
	case "GET":
		return pathItem.Get
	case "POST":
		return pathItem.Post
	case "PUT":
		return pathItem.Put
	case "DELETE":
		return pathItem.Delete
	}
	return nil
}

func generateReadFunction(resourceName string, op *openapi3.Operation) error {
	sanitizedName := sanitizeResourceName(resourceName)
	f := jen.NewFile("qdrant")

	f.Func().Id(fmt.Sprintf("resource%sRead", sanitizedName)).Params(
		jen.Id("ctx").Qual("context", "Context"),
		jen.Id("d").Op("*").Qual("github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema", "ResourceData"),
		jen.Id("m").Interface(),
	).Qual("github.com/hashicorp/terraform-plugin-sdk/v2/diag", "Diagnostics").Block(
		jen.Const().Id("errorPrefix").Op("=").Lit(fmt.Sprintf("error reading %s", sanitizedName)),
		jen.List(jen.Id("apiClient"), jen.Id("diagnostics")).Op(":=").Id("getClient").Call(jen.Id("m")),
		jen.If(jen.Id("diagnostics").Dot("HasError").Call()).Block(
			jen.Return(jen.Id("diagnostics")),
		),

		// Making the API call using the client
		jen.List(jen.Id("resp"), jen.Err()).Op(":=").Id("apiClient").Dot(fmt.Sprintf("Get%sWithResponse", sanitizedName)).Call(jen.Id("ctx")),
		jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Return(jen.Qual("github.com/hashicorp/terraform-plugin-sdk/v2/diag", "FromErr").Call(
				jen.Qual("fmt", "Errorf").Call(jen.Lit("%s: %v"), jen.Id("errorPrefix"), jen.Err()),
			)),
		),

		jen.If(jen.Id("resp").Dot("JSON422").Op("!=").Nil()).Block(
			jen.Return(jen.Qual("github.com/hashicorp/terraform-plugin-sdk/v2/diag", "FromErr").Call(
				jen.Qual("fmt", "Errorf").Call(jen.Lit("%s: %v"), jen.Id("errorPrefix"), jen.Id("getError").Call(jen.Id("resp").Dot("JSON422"))),
			)),
		),
		jen.If(jen.Id("resp").Dot("StatusCode").Call().Op("!=").Lit(200)).Block(
			jen.Return(jen.Qual("github.com/hashicorp/terraform-plugin-sdk/v2/diag", "FromErr").Call(
				jen.Qual("fmt", "Errorf").Call(jen.Id("getErrorMessage").Call(jen.Id("errorPrefix"), jen.Id("resp").Dot("HTTPResponse"))),
			)),
		),

		// Flatten and set data
		jen.Id("flattened").Op(":=").Id("flatten").Call(jen.Id("resp").Dot("JSON200")),
		jen.For(jen.List(jen.Id("k"), jen.Id("v")).Op(":=").Range().Id("flattened")).Block(
			jen.If(jen.Err().Op(":=").Id("d").Dot("Set").Call(jen.Id("k"), jen.Id("v")), jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Qual("github.com/hashicorp/terraform-plugin-sdk/v2/diag", "FromErr").Call(
					jen.Qual("fmt", "Errorf").Call(jen.Lit("%s: %v"), jen.Id("errorPrefix"), jen.Err()),
				)),
			),
		),

		jen.Id("d").Dot("SetId").Call(jen.Qual("time", "Now").Call().Dot("Format").Call(jen.Qual("time", "RFC3339"))),
		jen.Return(jen.Nil()),
	)

	// Save the generated code to a file
	fileName := fmt.Sprintf("qdrant/generated_resource_%s_read.go", strings.ToLower(sanitizedName))
	if err := f.Save(fileName); err != nil {
		return fmt.Errorf("error saving generated read function for resource %s: %v", sanitizedName, err)
	}

	return nil
}

func generateFlattenFunction(resourceName string) error {
	sanitizedName := sanitizeResourceName(resourceName)
	f := jen.NewFile("qdrant")

	f.Func().Id(fmt.Sprintf("flatten%s", sanitizedName)).Params(
		jen.Id("response").Op("*").Qual("terraform-provider-qdrant-cloud/v1/internal/client", fmt.Sprintf("%sOut", sanitizedName)),
	).Map(jen.String()).Interface().Block(
		jen.Return(jen.Id("flatten").Call(jen.Id("response"))),
	)

	// Save the generated code to a file
	fileName := fmt.Sprintf("qdrant/generated_flatten_%s.go", strings.ToLower(sanitizedName))
	if err := f.Save(fileName); err != nil {
		return fmt.Errorf("error saving generated flatten function for resource %s: %v", sanitizedName, err)
	}

	return nil
}
