package format

import (
    "bytes"
    "fmt"
    "strings"

    "github.com/burnlang/burn/pkg/ast"
    "github.com/burnlang/burn/pkg/lexer"
    "github.com/burnlang/burn/pkg/parser"
)


func Format(source string) (string, error) {
    
    lex := lexer.New(source)
    tokens, err := lex.Tokenize()
    if err != nil {
        return "", fmt.Errorf("lexical error: %v", err)
    }

    
    p := parser.New(tokens)
    program, err := p.Parse()
    if err != nil {
        return "", fmt.Errorf("parse error: %v", err)
    }

    
    var buf bytes.Buffer
    formatter := newAstFormatter(&buf)
    formatter.formatProgram(program)

    return buf.String(), nil
}


type astFormatter struct {
    buf       *bytes.Buffer
    indent    int
    newlineOk bool
}

func newAstFormatter(buf *bytes.Buffer) *astFormatter {
    return &astFormatter{
        buf:       buf,
        indent:    0,
        newlineOk: true,
    }
}


func (f *astFormatter) write(s string) {
    f.buf.WriteString(s)
}

func (f *astFormatter) writeIndent() {
    f.buf.WriteString(strings.Repeat("  ", f.indent))
}

func (f *astFormatter) writeLine(s string) {
    f.writeIndent()
    f.buf.WriteString(s)
    f.buf.WriteString("\n")
    f.newlineOk = true
}

func (f *astFormatter) writeNewLine() {
    if f.newlineOk {
        f.buf.WriteString("\n")
    }
}


func (f *astFormatter) formatProgram(program *ast.Program) {
    
    for i, decl := range program.Declarations {
        if i > 0 && needsExtraNewLine(program.Declarations[i-1], decl) {
            f.writeNewLine()
        }
        f.formatDeclaration(decl)
        f.writeNewLine()
    }
}


func needsExtraNewLine(prev, curr ast.Node) bool {
    
    if _, ok := prev.(*ast.ImportDeclaration); ok {
        if _, ok := curr.(*ast.ImportDeclaration); ok {
            return false
        }
        return true
    }

    
    if _, ok := prev.(*ast.MultiImportDeclaration); ok {
        if _, ok := curr.(*ast.ImportDeclaration); ok {
            return false
        }
        if _, ok := curr.(*ast.MultiImportDeclaration); ok {
            return false
        }
        return true
    }

    
    if _, ok := prev.(*ast.FunctionDeclaration); ok {
        return true
    }

    
    if _, ok := prev.(*ast.ClassDeclaration); ok {
        return true
    }

    return false
}


func (f *astFormatter) formatDeclaration(decl ast.Node) {
    switch d := decl.(type) {
    case *ast.ImportDeclaration:
        f.formatImport(d)
    case *ast.MultiImportDeclaration:
        f.formatMultiImport(d)
    case *ast.FunctionDeclaration:
        f.formatFunction(d)
    case *ast.ClassDeclaration:
        f.formatClass(d)
    case *ast.VariableDeclaration:
        f.formatVariable(d)
    case *ast.TypeDeclaration:
        f.formatType(d)
    case *ast.ExpressionStatement:
        f.formatExpressionInline(d.Expression)
    default:
        f.writeLine(fmt.Sprintf("// Unhandled declaration type: %T", d))
    }
}


func (f *astFormatter) formatImport(imp *ast.ImportDeclaration) {
    f.writeLine(fmt.Sprintf("import %q", imp.Path))
}


func (f *astFormatter) formatMultiImport(imp *ast.MultiImportDeclaration) {
    f.writeLine("import {")
    f.indent++
    for _, i := range imp.Imports {
        f.writeLine(fmt.Sprintf("%q,", i.Path))
    }
    f.indent--
    f.writeLine("}")
}


func (f *astFormatter) formatFunction(fn *ast.FunctionDeclaration) {
    
    params := formatParams(fn.Parameters)
    returnType := ""
    if fn.ReturnType != "" {
        returnType = ": " + fn.ReturnType
    }
    f.writeLine(fmt.Sprintf("fun %s(%s)%s {", fn.Name, params, returnType))
    
    
    f.indent++
    for i, stmt := range fn.Body {
        if i > 0 {
            
            if needsBlankLine(fn.Body[i-1], stmt) {
                f.writeNewLine()
            }
        }
        f.formatStatement(stmt)
    }
    f.indent--
    
    f.writeLine("}")
}


func (f *astFormatter) formatClass(class *ast.ClassDeclaration) {
    f.writeLine(fmt.Sprintf("class %s {", class.Name))
    f.indent++
    
    for i, method := range class.Methods {
        if i > 0 {
            f.writeNewLine()
        }
        f.formatFunction(method)
    }
    
    f.indent--
    f.writeLine("}")
}


func (f *astFormatter) formatVariable(varDecl *ast.VariableDeclaration) {
    keyword := "var"
    if varDecl.IsConst {
        keyword = "const"
    }
    
    typeAnnotation := ""
    if varDecl.Type != "" {
        typeAnnotation = ": " + varDecl.Type
    }
    
    if varDecl.Value == nil {
        f.writeLine(fmt.Sprintf("%s %s%s", keyword, varDecl.Name, typeAnnotation))
        return
    }
    
    f.writeIndent()
    f.write(fmt.Sprintf("%s %s%s = ", keyword, varDecl.Name, typeAnnotation))
    f.formatExpressionInline(varDecl.Value)
    f.write("\n")
}


func (f *astFormatter) formatType(typeDecl *ast.TypeDeclaration) {
    f.writeLine(fmt.Sprintf("type %s {", typeDecl.Name))
    f.indent++
    
    for i, field := range typeDecl.Fields {
        comma := ","
        if i == len(typeDecl.Fields)-1 {
            comma = ""
        }
        f.writeLine(fmt.Sprintf("%s: %s%s", field.Name, field.Type, comma))
    }
    
    f.indent--
    f.writeLine("}")
}


func (f *astFormatter) formatStatement(stmt ast.Node) {
    switch s := stmt.(type) {
    case *ast.ExpressionStatement:
        f.writeIndent()
        f.formatExpressionInline(s.Expression)
        f.write("\n")
    case *ast.ReturnStatement:
        f.writeIndent()
        f.write("return")
        if s.Value != nil {
            f.write(" ")
            f.formatExpressionInline(s.Value)
        }
        f.write("\n")
    case *ast.IfStatement:
        f.formatIfStatement(s)
    case *ast.WhileStatement:
        f.formatWhileStatement(s)
    case *ast.ForStatement:
        f.formatForStatement(s)
    case *ast.VariableDeclaration:
        f.formatVariable(s)
    case *ast.BlockStatement:
        f.formatBlockStatement(s)
    default:
        f.writeLine(fmt.Sprintf("// Unhandled statement type: %T", s))
    }
}


func (f *astFormatter) formatIfStatement(ifStmt *ast.IfStatement) {
    f.writeIndent()
    f.write("if (")
    f.formatExpressionInline(ifStmt.Condition)
    f.write(") {\n")
    
    f.indent++
    for _, stmt := range ifStmt.Consequence {
        f.formatStatement(stmt)
    }
    f.indent--
    
    if len(ifStmt.Alternative) > 0 {
        f.writeLine("} else {")
        
        f.indent++
        for _, stmt := range ifStmt.Alternative {
            f.formatStatement(stmt)
        }
        f.indent--
        
        f.writeLine("}")
    } else {
        f.writeLine("}")
    }
}


func (f *astFormatter) formatWhileStatement(whileStmt *ast.WhileStatement) {
    f.writeIndent()
    f.write("while (")
    f.formatExpressionInline(whileStmt.Condition)
    f.write(") {\n")
    
    f.indent++
    for _, stmt := range whileStmt.Body {
        f.formatStatement(stmt)
    }
    f.indent--
    
    f.writeLine("}")
}


func (f *astFormatter) formatForStatement(forStmt *ast.ForStatement) {
    f.writeIndent()
    f.write("for (")
    
    
    if forStmt.Init != nil {
        switch init := forStmt.Init.(type) {
        case *ast.VariableDeclaration:
            keyword := "var"
            if init.IsConst {
                keyword = "const"
            }
            f.write(fmt.Sprintf("%s %s", keyword, init.Name))
            if init.Type != "" {
                f.write(fmt.Sprintf(": %s", init.Type))
            }
            f.write(" = ")
            f.formatExpressionInline(init.Value)
        default:
            f.formatExpressionInline(forStmt.Init)
        }
    }
    
    f.write("; ")
    
    
    f.formatExpressionInline(forStmt.Condition)
    f.write("; ")
    
    
    f.formatExpressionInline(forStmt.Update)
    f.write(") {\n")
    
    
    f.indent++
    for _, stmt := range forStmt.Body {
        f.formatStatement(stmt)
    }
    f.indent--
    
    f.writeLine("}")
}


func (f *astFormatter) formatBlockStatement(block *ast.BlockStatement) {
    if len(block.Statements) == 0 {
        f.writeLine("{}")
        return
    }
    
    f.writeLine("{")
    f.indent++
    
    for _, stmt := range block.Statements {
        f.formatStatement(stmt)
    }
    
    f.indent--
    f.writeLine("}")
}


func (f *astFormatter) formatExpressionInline(expr ast.Node) {
    switch e := expr.(type) {
    case *ast.LiteralExpression:
        switch val := e.Value.(type) {
        case string:
            f.write(fmt.Sprintf("%q", val))
        default:
            f.write(fmt.Sprintf("%v", val))
        }
    case *ast.VariableExpression:
        f.write(e.Name)
    case *ast.CallExpression:
        f.formatCallExpression(e)
    case *ast.BinaryExpression:
        f.formatBinaryExpression(e)
    case *ast.UnaryExpression:
        f.formatUnaryExpression(e)
    case *ast.ObjectLiteral:
        f.formatObjectLiteral(e)
    case *ast.ArrayLiteral:
        f.formatArrayLiteral(e)
    case *ast.ClassMethodCallExpression:
        f.formatClassMethodCall(e)
    case *ast.MemberExpression:
        f.formatMemberExpression(e)
    case *ast.IndexExpression:
        f.formatIndexExpression(e)
    default:
        f.write(fmt.Sprintf("/* Unhandled expr type: %T */", e))
    }
}


func (f *astFormatter) formatBinaryExpression(expr *ast.BinaryExpression) {
    f.write("(")
    f.formatExpressionInline(expr.Left)
    f.write(fmt.Sprintf(" %s ", expr.Operator))
    f.formatExpressionInline(expr.Right)
    f.write(")")
}


func (f *astFormatter) formatUnaryExpression(expr *ast.UnaryExpression) {
    f.write(expr.Operator)
    f.formatExpressionInline(expr.Right)
}


func (f *astFormatter) formatCallExpression(expr *ast.CallExpression) {
    f.formatExpressionInline(expr.Callee)
    f.write("(")
    
    for i, arg := range expr.Arguments {
        if i > 0 {
            f.write(", ")
        }
        f.formatExpressionInline(arg)
    }
    
    f.write(")")
}


func (f *astFormatter) formatObjectLiteral(obj *ast.ObjectLiteral) {
    if len(obj.Fields) == 0 {
        f.write("{}")
        return
    }
    
    f.write("{\n")
    f.indent++
    
    i := 0
    for key, value := range obj.Fields {
        f.writeIndent()
        f.write(fmt.Sprintf("%s: ", key))
        f.formatExpressionInline(value)
        
        if i < len(obj.Fields)-1 {
            f.write(",")
        }
        f.write("\n")
        i++
    }
    
    f.indent--
    f.writeIndent()
    f.write("}")
}


func (f *astFormatter) formatArrayLiteral(arr *ast.ArrayLiteral) {
    if len(arr.Elements) == 0 {
        f.write("[]")
        return
    }
    
    
    if len(arr.Elements) <= 3 && isSimpleExpressions(arr.Elements) {
        f.write("[")
        for i, elem := range arr.Elements {
            if i > 0 {
                f.write(", ")
            }
            f.formatExpressionInline(elem)
        }
        f.write("]")
        return
    }
    
    
    f.write("[\n")
    f.indent++
    
    for i, elem := range arr.Elements {
        f.writeIndent()
        f.formatExpressionInline(elem)
        if i < len(arr.Elements)-1 {
            f.write(",")
        }
        f.write("\n")
    }
    
    f.indent--
    f.writeIndent()
    f.write("]")
}


func (f *astFormatter) formatClassMethodCall(expr *ast.ClassMethodCallExpression) {
    f.write(fmt.Sprintf("%s.%s(", expr.ClassName, expr.MethodName))
    
    for i, arg := range expr.Arguments {
        if i > 0 {
            f.write(", ")
        }
        f.formatExpressionInline(arg)
    }
    
    f.write(")")
}


func (f *astFormatter) formatMemberExpression(expr *ast.MemberExpression) {
    f.formatExpressionInline(expr.Object)
    f.write(".")
    f.write(expr.Property)
}


func (f *astFormatter) formatIndexExpression(expr *ast.IndexExpression) {
    f.formatExpressionInline(expr.Array)
    f.write("[")
    f.formatExpressionInline(expr.Index)
    f.write("]")
}


func formatParams(params []*ast.Parameter) string {
    var parts []string
    for _, param := range params {
        if param.Type != "" {
            parts = append(parts, fmt.Sprintf("%s: %s", param.Name, param.Type))
        } else {
            parts = append(parts, param.Name)
        }
    }
    return strings.Join(parts, ", ")
}


func needsBlankLine(prev, curr ast.Node) bool {
    
    switch curr.(type) {
    case *ast.IfStatement, *ast.WhileStatement, *ast.ForStatement:
        return true
    }
    
    
    switch prev.(type) {
    case *ast.IfStatement, *ast.WhileStatement, *ast.ForStatement, *ast.BlockStatement:
        return true
    }
    
    
    if _, ok := prev.(*ast.VariableDeclaration); ok {
        if _, ok := curr.(*ast.VariableDeclaration); ok {
            return false
        }
    }
    
    return false
}


func isSimpleExpressions(exprs []ast.Node) bool {
    for _, expr := range exprs {
        switch expr.(type) {
        case *ast.LiteralExpression, *ast.VariableExpression:
            continue
        default:
            return false
        }
    }
    return true
}