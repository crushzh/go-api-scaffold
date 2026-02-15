// Module code generator
//
// Usage:
//   go run cmd/gen/main.go -name order -cn Order
//   make gen name=order cn=Order
//
// Output:
//   internal/handler/order_handler.go    — HTTP CRUD endpoints
//   internal/service/order_service.go    — Business logic layer
//   internal/model/order.go             — Data model
//   internal/store/order_repo.go        — Data repository
//   Also auto-registers routes in router.go
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
)

type ModuleData struct {
	Name        string // order
	PascalName  string // Order
	CamelName   string // order
	SnakeName   string // order
	KebabName   string // order
	PluralName  string // orders
	ChineseName string // Order
	ModulePath  string // go-api-scaffold (go module name)
}

func main() {
	name := flag.String("name", "", "module name (lowercase, e.g. order)")
	cn := flag.String("cn", "", "display name (e.g. Order)")
	modulePath := flag.String("module", "", "Go module path (auto-detected)")
	flag.Parse()

	if *name == "" {
		fmt.Fprintln(os.Stderr, "error: -name is required")
		fmt.Fprintln(os.Stderr, "Usage: go run cmd/gen/main.go -name order -cn Order")
		os.Exit(1)
	}

	if *cn == "" {
		*cn = *name
	}

	// Auto-detect module path
	if *modulePath == "" {
		*modulePath = detectModulePath()
	}

	data := ModuleData{
		Name:        strings.ToLower(*name),
		PascalName:  toPascalCase(*name),
		CamelName:   toCamelCase(*name),
		SnakeName:   toSnakeCase(*name),
		KebabName:   toKebabCase(*name),
		PluralName:  toPlural(strings.ToLower(*name)),
		ChineseName: *cn,
		ModulePath:  *modulePath,
	}

	fmt.Printf("generating module: %s (%s)\n", data.PascalName, data.ChineseName)

	// Generate files
	files := []struct {
		tmpl string
		out  string
	}{
		{"templates/handler.go.tmpl", fmt.Sprintf("internal/handler/%s_handler.go", data.SnakeName)},
		{"templates/service.go.tmpl", fmt.Sprintf("internal/service/%s_service.go", data.SnakeName)},
		{"templates/model.go.tmpl", fmt.Sprintf("internal/model/%s.go", data.SnakeName)},
		{"templates/store.go.tmpl", fmt.Sprintf("internal/store/%s_repo.go", data.SnakeName)},
	}

	for _, f := range files {
		if err := generateFile(f.tmpl, f.out, data); err != nil {
			fmt.Fprintf(os.Stderr, "failed to generate %s: %v\n", f.out, err)
			os.Exit(1)
		}
		fmt.Printf("  + %s\n", f.out)
	}

	// Auto-register route
	if err := appendRoute(data); err != nil {
		fmt.Fprintf(os.Stderr, "  ! auto-register route failed: %v (add manually)\n", err)
	} else {
		fmt.Println("  + route registered in router.go")
	}

	// Auto-register model migration
	if err := appendMigration(data); err != nil {
		fmt.Fprintf(os.Stderr, "  ! auto-register migration failed: %v (add manually)\n", err)
	} else {
		fmt.Println("  + migration registered in store.go")
	}

	fmt.Printf("\nmodule %s generated successfully!\n", data.PascalName)
	fmt.Println("\nNext steps:")
	fmt.Printf("  1. edit internal/model/%s.go — add model fields\n", data.SnakeName)
	fmt.Printf("  2. edit internal/service/%s_service.go — implement business logic\n", data.SnakeName)
	fmt.Printf("  3. run make docs — update Swagger docs\n")
}

func generateFile(tmplPath, outPath string, data ModuleData) error {
	// Check if target file already exists
	if _, err := os.Stat(outPath); err == nil {
		return fmt.Errorf("file already exists: %s", outPath)
	}

	// Read template
	tmplContent, err := os.ReadFile(tmplPath)
	if err != nil {
		return fmt.Errorf("read template: %w", err)
	}

	// Parse and execute template
	t, err := template.New(filepath.Base(tmplPath)).Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}

	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return t.Execute(f, data)
}

// appendRoute inserts route registration code at the marker comment in router.go
func appendRoute(data ModuleData) error {
	routerFile := "internal/handler/router.go"
	content, err := os.ReadFile(routerFile)
	if err != nil {
		return err
	}

	marker := "// GEN:ROUTE_REGISTER - Auto-appended by code generator, do not remove"
	routeCode := fmt.Sprintf(`			// %s module
			%sHandler := New%sHandler(%sSvc)
			%s := authorized.Group("/%s")
			{
				%s.GET("", %sHandler.List)
				%s.POST("", %sHandler.Create)
				%s.GET("/:id", %sHandler.Get)
				%s.PUT("/:id", %sHandler.Update)
				%s.DELETE("/:id", %sHandler.Delete)
			}

			`,
		data.PascalName,
		data.CamelName, data.PascalName, data.CamelName,
		data.PluralName, data.PluralName,
		data.PluralName, data.CamelName,
		data.PluralName, data.CamelName,
		data.PluralName, data.CamelName,
		data.PluralName, data.CamelName,
		data.PluralName, data.CamelName,
	)

	newContent := strings.Replace(string(content), marker, routeCode+marker, 1)
	if newContent == string(content) {
		return fmt.Errorf("route marker comment not found")
	}

	return os.WriteFile(routerFile, []byte(newContent), 0o644)
}

// appendMigration inserts model migration at the marker comment in store.go
func appendMigration(data ModuleData) error {
	storeFile := "internal/store/store.go"
	content, err := os.ReadFile(storeFile)
	if err != nil {
		return err
	}

	marker := "// GEN:MODEL_MIGRATE - Auto-appended by code generator, do not remove"
	migrationCode := fmt.Sprintf("\t\t&model.%s{},\n\t\t", data.PascalName)

	newContent := strings.Replace(string(content), marker, migrationCode+marker, 1)
	if newContent == string(content) {
		return fmt.Errorf("migration marker comment not found")
	}

	// Ensure model package is imported
	if !strings.Contains(newContent, `"go-api-scaffold/internal/model"`) && !strings.Contains(newContent, fmt.Sprintf(`"%s/internal/model"`, data.ModulePath)) {
		newContent = strings.Replace(newContent, `"go-api-scaffold/pkg/logger"`,
			fmt.Sprintf(`"go-api-scaffold/internal/model"\n\t"go-api-scaffold/pkg/logger"`), 1)
	}

	return os.WriteFile(storeFile, []byte(newContent), 0o644)
}

// appendServiceInit prints a reminder to register the service in main.go
// Note: main.go is not auto-modified because service init may require different params

func detectModulePath() string {
	content, err := os.ReadFile("go.mod")
	if err != nil {
		return "go-api-scaffold"
	}
	for _, line := range strings.Split(string(content), "\n") {
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module "))
		}
	}
	return "go-api-scaffold"
}

// ========================
// Naming conversion utilities
// ========================

func toPascalCase(s string) string {
	parts := splitWords(s)
	for i, p := range parts {
		parts[i] = strings.Title(strings.ToLower(p))
	}
	return strings.Join(parts, "")
}

func toCamelCase(s string) string {
	pascal := toPascalCase(s)
	if pascal == "" {
		return ""
	}
	runes := []rune(pascal)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

func toSnakeCase(s string) string {
	parts := splitWords(s)
	for i, p := range parts {
		parts[i] = strings.ToLower(p)
	}
	return strings.Join(parts, "_")
}

func toKebabCase(s string) string {
	parts := splitWords(s)
	for i, p := range parts {
		parts[i] = strings.ToLower(p)
	}
	return strings.Join(parts, "-")
}

func toPlural(s string) string {
	if strings.HasSuffix(s, "s") || strings.HasSuffix(s, "x") || strings.HasSuffix(s, "ch") || strings.HasSuffix(s, "sh") {
		return s + "es"
	}
	if strings.HasSuffix(s, "y") && len(s) > 1 {
		prev := s[len(s)-2]
		if prev != 'a' && prev != 'e' && prev != 'i' && prev != 'o' && prev != 'u' {
			return s[:len(s)-1] + "ies"
		}
	}
	return s + "s"
}

func splitWords(s string) []string {
	// Supports snake_case, kebab-case, camelCase, PascalCase
	s = strings.ReplaceAll(s, "-", "_")
	parts := strings.Split(s, "_")
	var result []string
	for _, p := range parts {
		// Split camelCase
		var current []rune
		for i, r := range p {
			if unicode.IsUpper(r) && i > 0 {
				result = append(result, string(current))
				current = nil
			}
			current = append(current, r)
		}
		if len(current) > 0 {
			result = append(result, string(current))
		}
	}
	return result
}
