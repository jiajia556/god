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

// 改进：基于括号层级与引号状态提取字段定义，避免被类型/注释/索引内的逗号误分割
func extractFieldDefinitions(sql string) ([]string, error) {
	// 找到第一个左括号（字段定义起点）
	start := strings.Index(sql, "(")
	if start < 0 {
		return nil, fmt.Errorf("field definitions not found")
	}

	// 从 start 向后找到匹配的右括号（考虑嵌套）
	level := 0
	end := -1
	for i := start; i < len(sql); i++ {
		ch := sql[i]
		if ch == '(' {
			level++
		} else if ch == ')' {
			level--
			if level == 0 {
				end = i
				break
			}
		}
	}
	if end == -1 || end <= start {
		return nil, fmt.Errorf("field definitions not found (unmatched parentheses)")
	}

	inner := sql[start+1 : end]

	// 按最外层逗号分割（忽略括号内和引号内的逗号）
	defs := splitFieldDefinitions(inner)

	// 清理并返回
	out := make([]string, 0, len(defs))
	for _, d := range defs {
		s := strings.TrimSpace(d)
		if s != "" {
			out = append(out, s)
		}
	}
	return out, nil
}

// splitFieldDefinitions: 在没有正则限制下按最外层逗号切分定义
// 新增：考虑单/双引号以及反斜杠转义，避免在引号内分割
func splitFieldDefinitions(body string) []string {
	var defs []string
	var cur strings.Builder
	level := 0
	inSingle := false
	inDouble := false
	escaped := false

	for i := 0; i < len(body); i++ {
		ch := body[i]

		// 处理转义：如果前一个字符是反斜杠，则当前字符是被转义的
		if ch == '\\' && !escaped {
			escaped = true
			cur.WriteByte(ch)
			continue
		}

		// 如果在单引号或双引号内，被转义字符不作为引号结束
		if !escaped {
			if ch == '\'' && !inDouble {
				inSingle = !inSingle
				cur.WriteByte(ch)
				continue
			}
			if ch == '"' && !inSingle {
				inDouble = !inDouble
				cur.WriteByte(ch)
				continue
			}
		}

		// 非引号中且非转义，处理括号层级
		if !inSingle && !inDouble {
			if ch == '(' {
				level++
			} else if ch == ')' {
				if level > 0 {
					level--
				}
			}
		}

		// 分割条件：最外层逗号（level==0）且不在任意引号内
		if ch == ',' && level == 0 && !inSingle && !inDouble {
			part := strings.TrimSpace(cur.String())
			if part != "" {
				defs = append(defs, part)
			}
			cur.Reset()
			escaped = false
			continue
		}

		// 正常写入当前字符
		cur.WriteByte(ch)
		// 重置 escaped 标志（只有紧接着的字符被转义）
		escaped = false
	}

	// 追加剩余
	if s := strings.TrimSpace(cur.String()); s != "" {
		defs = append(defs, s)
	}
	return defs
}

func parseField(def string) (fieldInfo, error) {
	// Improved regex to fully capture type description
	re := regexp.MustCompile("[\x60]?(\\w+)[\x60]?\\s+(.+)")
	matches := re.FindStringSubmatch(def)
	if len(matches) < 3 {
		return fieldInfo{}, fmt.Errorf("invalid field definition: %s", def)
	}

	fieldName := matches[1]
	// 仅对以小写字母开头的 Column 名进行生成��保持旧逻辑）
	if len(fieldName) == 0 || []byte(fieldName)[0] < 'a' || []byte(fieldName)[0] > 'z' {
		return fieldInfo{}, nil
	}
	typeInfo := strings.ToLower(strings.TrimSpace(matches[2]))

	// 保留原有类型映射
	goType, tags := mapTypeAndTags(typeInfo)

	// 保守增强：识别常见约束并加入 tags（不改变 goType 的映射）
	if strings.Contains(typeInfo, "unsigned") {
		tags["unsigned"] = "true"
	}
	if strings.Contains(typeInfo, "auto_increment") || strings.Contains(typeInfo, "autoincrement") {
		tags["autoIncrement"] = "true"
	}
	if strings.Contains(typeInfo, "primary key") || strings.Contains(typeInfo, "primary_key") {
		tags["primaryKey"] = "true"
	}
	if strings.Contains(typeInfo, "not null") {
		tags["notNull"] = "true"
	}
	// default 值提取（尽量简单提取）
	if idx := strings.Index(strings.ToUpper(typeInfo), "DEFAULT "); idx >= 0 {
		after := strings.TrimSpace(typeInfo[idx+8:])
		if after != "" {
			if strings.HasPrefix(after, "'") || strings.HasPrefix(after, "\"") {
				q := after[0]
				if j := strings.IndexByte(after[1:], q); j >= 0 {
					tags["default"] = after[1 : 1+j]
				} else {
					tags["default"] = strings.Trim(after, "'\"")
				}
			} else {
				parts := strings.Fields(after)
				if len(parts) > 0 {
					tags["default"] = strings.TrimRight(parts[0], ",")
				}
			}
		}
	}

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
