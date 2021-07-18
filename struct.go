package struct_generate

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type GStruct struct {
	Name         string
	FileName     string
	PkgName      string
	Comment      string
	EmbedStructs []string
	Fields       map[string]*GStructField
	Methods      map[string]*GStructMethod
}

type GStructField struct {
	Name    string
	Type    string
	Comment string
	Tags    string
	Index   int
}

func createGStruct() *GStruct {
	ret := &GStruct{}
	ret.Fields = make(map[string]*GStructField, 0)
	ret.Methods = make(map[string]*GStructMethod, 0)
	ret.EmbedStructs = make([]string, 0)
	return ret
}

type GStructMethod struct {
	Name     string
	Index    int
	FileName string
	StartPos token.Pos
	EndPos   token.Pos
	Content  string
}

type StructList []*GStruct

func ParseFile(pkgDir string) (ret StructList, err error) {
	ret = make(StructList, 0)

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, pkgDir, nil, 0)
	if err != nil {
		return
	}

	for pkgName, pkg := range pkgs {
		for fileName, pkgFile := range pkg.Files {
			// list scope structs
			for structName, structObj := range pkgFile.Scope.Objects {
				if structObj.Kind != ast.Typ {
					continue
				}
				gstruct := createGStruct()
				gstruct.FileName = fileName
				gstruct.Name = structName
				gstruct.PkgName = pkgName

				if structDecl, ok := structObj.Decl.(*ast.TypeSpec); ok {
					if structType, ok := structDecl.Type.(*ast.StructType); ok {
						for fieldIndex, field := range structType.Fields.List {
							if field.Names == nil {
								embedStructName := getTypeName(field.Type, nil)
								gstruct.EmbedStructs = append(gstruct.EmbedStructs, embedStructName)
								continue
							}

							gField := &GStructField{}
							gField.Name = field.Names[0].Name
							gField.Index = fieldIndex
							gField.Type = getTypeName(field.Type, nil)

							if field.Tag != nil {
								gField.Tags = field.Tag.Value
							}

							gstruct.Fields[gField.Name] = gField
						}
					}
				}

				ret = append(ret, gstruct)
			}
			// list structs methods
		}
	}

	fmt.Println(pkgs)
	return
}

func getTypeName(expr ast.Expr, parent ast.Expr) string {
	if fieldType, ok := expr.(*ast.Ident); ok {
		return fieldType.Name
	}
	// 指针类型
	if fieldType, ok := expr.(*ast.StarExpr); ok {
		return "*" + getTypeName(fieldType.X, fieldType)
	}

	if fieldType, ok := expr.(*ast.MapType); ok {
		return "map[" + getTypeName(fieldType.Key, fieldType) + "]" + getTypeName(fieldType.Value, fieldType)
	}

	if fieldType, ok := expr.(*ast.ArrayType); ok {
		return "[]" + getTypeName(fieldType.Elt, fieldType)
	}

	if fieldType, ok := expr.(*ast.InterfaceType); ok {
		if fieldType.Methods.NumFields() == 0 {
			return "interface{}"
		}

		funcList := make([]string, fieldType.Methods.NumFields())

		for i, f := range fieldType.Methods.List {
			funcList[i] = getNames(f.Names) + getTypeName(f.Type, fieldType)
		}

		return fmt.Sprintf("interface{\n%s\n}", strings.Join(funcList, "\n"))
	}

	if fieldType, ok := expr.(*ast.FuncType); ok {
		// 处理参数
		params := getFieldString(fieldType.Params)
		results := getFieldString(fieldType.Results)
		// 处理返回
		// 有多个返回值时
		if len(fieldType.Results.List) > 1 {
			format := "func(%s)(%s)"
			if parent != nil && isInterfaceType(parent) {
				format = "(%s)(%s)"
			}
			return fmt.Sprintf(format, params, results)
		}

		// 只有一个返回值, 并且那个返回值没有名字
		if len(fieldType.Results.List) == 1 &&
			fieldType.Results.List[0].Names == nil {
			format := "func(%s)%s"
			if parent != nil && isInterfaceType(parent) {
				format = "(%s)%s"
			}
			return fmt.Sprintf(format, params, results)
		}
		// 没有返回值

		format := "func (%s)"
		if parent != nil && isInterfaceType(parent) {
			format = "(%s)"
		}
		return fmt.Sprintf(format, params)
	}

	return ""
}

func getFieldString(fieldList *ast.FieldList) string {
	fields := make([]string, fieldList.NumFields())
	for i, p := range fieldList.List {
		if len(p.Names) > 0 {
			fields[i] = getNames(p.Names) + " " + getTypeName(p.Type, nil)
		} else {
			fields[i] = getTypeName(p.Type, nil)
		}
	}

	return strings.Join(fields, ", ")
}

func getNames(names []*ast.Ident) string {
	if len(names) == 0 {
		return ""
	}

	list := make([]string, len(names))

	for i, name := range names {
		list[i] = name.Name
	}

	return strings.Join(list, ", ")
}

func isInterfaceType(expr ast.Expr) bool {
	if _, ok := expr.(*ast.InterfaceType); ok {
		return true
	}
	return false
}
