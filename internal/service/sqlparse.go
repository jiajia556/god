package service

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type fieldInfo struct {
	name     string
	typeName string
	gormTags string
	jsonTag  string
}

func ExtractCreateTables(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var createTables []string
	scanner := bufio.NewScanner(file)
	var currentStmt strings.Builder
	capturing := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "--") || strings.HasPrefix(line, "/*") {
			continue
		}
		if strings.HasPrefix(line, "CREATE TABLE") {
			capturing = true
			currentStmt.WriteString(line + "\n")
			continue
		}

		if strings.HasPrefix(line, ")") && strings.HasSuffix(line, ";") {
			capturing = false
			currentStmt.WriteString(line + "\n")
			createTables = append(createTables, currentStmt.String())
			currentStmt = strings.Builder{}
		}

		if capturing {
			currentStmt.WriteString(line + "\n")
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// 处理最后可能未闭合的语句
	if capturing && currentStmt.Len() > 0 {
		createTables = append(createTables, currentStmt.String())
	}

	return createTables, nil
}

// GenerateStruct generates Go struct definition from SQL create table statement
func GenerateModelStruct(sql string) (string, string, error) {
	tableName, fields, err := parseSQL(sql)
	if err != nil {
		return "", "", err
	}

	return buildStruct(tableName, fields), toCamelCase(tableName), nil
}

func parseSQL(sql string) (string, []fieldInfo, error) {
	tableName, err := extractTableName(sql)
	if err != nil {
		return "", nil, err
	}

	fieldDefinitions, err := extractFieldDefinitions(sql)
	if err != nil {
		return "", nil, err
	}

	var fields []fieldInfo
	for _, def := range fieldDefinitions {
		fi, err := parseField(def)
		if err != nil {
			return "", nil, err
		}
		if fi.name == "" {
			continue
		}
		fields = append(fields, fi)
	}

	return tableName, fields, nil
}

func extractTableName(sql string) (string, error) {
	re := regexp.MustCompile(`(?i)CREATE\s+TABLE\s+[\x60]?(\w+)[\x60]?`)
	matches := re.FindStringSubmatch(sql)
	if len(matches) < 2 {
		return "", fmt.Errorf("table name not found")
	}
	return toCamelCase(matches[1]), nil
}

func extractFieldDefinitions(sql string) ([]string, error) {
	re := regexp.MustCompile(`(?s)\((.*)\)`)
	matches := re.FindStringSubmatch(sql)
	if len(matches) < 2 {
		return nil, fmt.Errorf("field definitions not found")
	}

	// Improved regex to fully capture field definitions
	fieldRe := regexp.MustCompile("[\x60]?(\\w+)[\x60]?\\s+([^,]+)")
	return fieldRe.FindAllString(matches[1], -1), nil
}

func parseField(def string) (fieldInfo, error) {
	// Improved regex to fully capture type description
	re := regexp.MustCompile("[\x60]?(\\w+)[\x60]?\\s+(.+)")
	matches := re.FindStringSubmatch(def)
	if len(matches) < 3 {
		return fieldInfo{}, fmt.Errorf("invalid field definition: %s", def)
	}

	fieldName := matches[1]
	if []byte(fieldName)[0] < 'a' || []byte(fieldName)[0] > 'z' {
		return fieldInfo{}, nil
	}
	typeInfo := strings.ToLower(strings.TrimSpace(matches[2]))

	goType, tags := mapTypeAndTags(typeInfo)
	//isNullable := strings.Contains(typeInfo, "null")

	// Handle pointer types
	//if isNullable && !isSpecialType(goType) {
	//	goType = "*" + goType
	//}

	return fieldInfo{
		name:     toCamelCase(fieldName),
		typeName: goType,
		gormTags: buildGormTags(fieldName, tags),
		jsonTag:  toSnakeCase(fieldName),
	}, nil
}

func mapTypeAndTags(sqlType string) (string, map[string]string) {
	tags := make(map[string]string)
	baseType := regexp.MustCompile(`^(\w+)(?:\(.*?\))?`).FindString(sqlType)
	baseType = strings.ToLower(baseType)
	unsigned := strings.Contains(sqlType, "unsigned")

	var goType string
	switch {
	case strings.HasPrefix(baseType, "tinyint"):
		if unsigned {
			goType = "uint8"
		} else {
			goType = "int8"
		}
	case strings.HasPrefix(baseType, "int") || strings.HasPrefix(baseType, "bigint"):
		if unsigned {
			goType = "uint64"
		} else {
			goType = "int64"
		}
	case strings.HasPrefix(baseType, "decimal"):
		goType = "decimal.Decimal"
	case strings.Contains(baseType, "datetime") || strings.Contains(baseType, "timestamp"):
		goType = "mytime.DateTime"
	case strings.HasPrefix(baseType, "varchar"), strings.HasPrefix(baseType, "text"):
		goType = "string"
	case strings.HasPrefix(baseType, "boolean"):
		goType = "bool"
	default:
		goType = "string"
	}

	// Handle special tags
	//if strings.Contains(sqlType, "auto_increment") {
	//	tags["autoIncrement"] = "true"
	//}
	//if strings.Contains(sqlType, "primary key") {
	//	tags["primaryKey"] = "true"
	//}

	return goType, tags
}

func buildGormTags(fieldName string, tags map[string]string) string {
	parts := []string{"column:" + fieldName}
	for k, v := range tags {
		if v == "true" {
			parts = append(parts, k)
		} else {
			parts = append(parts, fmt.Sprintf("%s:%s", k, v))
		}
	}
	return strings.Join(parts, ";")
}

func isSpecialType(t string) bool {
	return strings.Contains(t, ".") || t == "string"
}

func toCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i := range parts {
		parts[i] = strings.Title(parts[i])
	}
	return strings.Join(parts, "")
}

func toSnakeCase(s string) string {
	return strings.ToLower(s)
}

func buildStruct(tableName string, fields []fieldInfo) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("type %s struct {\n", toCamelCase(tableName)))

	for _, f := range fields {
		sb.WriteString(fmt.Sprintf("    %-8s %-16s `gorm:\"%s\" json:\"%s\"`\n",
			f.name, f.typeName, f.gormTags, f.jsonTag))
	}

	sb.WriteString("}")
	return sb.String()
}
