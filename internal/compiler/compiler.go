//nolint:goconst
package compiler

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/azin-lang/Azin/internal/codegen/c"
	"github.com/azin-lang/Azin/internal/fs"
	"github.com/azin-lang/Azin/internal/optimizer"
	"github.com/azin-lang/Azin/pkg/ast"
	"github.com/azin-lang/Azin/pkg/diagnostics"
	"github.com/azin-lang/Azin/pkg/lexer"
	"github.com/azin-lang/Azin/pkg/parser"
	"github.com/azin-lang/Azin/pkg/sema"
	"github.com/azin-lang/Azin/pkg/source"
)

func runCompilerCommand(name string, args []string, label, output string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("[%s] Compiling...\n", label)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s compilation failed: %w", label, err)
	}

	fmt.Printf("[Success] %s\n", output)
	return nil
}

func msvcOptimization(opt string) string {
	switch opt {
	case "1", "s", "z":
		return "/O1"
	case "2":
		return "/O2"
	case "3":
		return "/Ox"
	default:
		return "/Od"
	}
}

func runMSVC(cl string, sources []string, exeName string, opts Options) error {
	opt := msvcOptimization(opts.Optimization)
	args := append([]string{"/nologo", "/std:c11", opt, "/Fe:" + exeName}, sources...)

	return runCompilerCommand(
		cl,
		args,
		"MSVC",
		exeName,
	)
}

func runClang(clang string, sources []string, exeName string, opts Options) error {
	args := append([]string{"-std=c11", "-O" + opts.Optimization}, sources...)
	args = append(args, "-o", exeName, "-lm")

	return runCompilerCommand(
		clang,
		args,
		"Clang",
		exeName,
	)
}

func runGCC(gcc string, sources []string, exeName string, opts Options) error {
	args := append([]string{"-std=c11", "-O" + opts.Optimization}, sources...)
	args = append(args, "-o", exeName, "-lm")

	return runCompilerCommand(
		gcc,
		args,
		"GCC",
		exeName,
	)
}

type compiler struct {
	name string
	run  func(string, []string, string, Options) error
}

func runCompiler(sources []string, exeName string, opts Options) error {
	var compilers []compiler

	switch runtime.GOOS {
	case "windows":
		compilers = []compiler{
			{"cl.exe", runMSVC},
			{"clang", runClang},
			{"gcc", runGCC},
		}

	case "darwin":
		compilers = []compiler{
			{"clang", runClang},
			{"gcc", runGCC},
		}

	default:
		compilers = []compiler{
			{"gcc", runGCC},
			{"clang", runClang},
		}
	}

	for _, c := range compilers {
		if path, err := exec.LookPath(c.name); err == nil {
			return c.run(path, sources, exeName, opts)
		}
	}

	return fmt.Errorf("no supported C compiler found")
}

func writeCOutput(code, output string) error {
	if output == "" {
		output = "output.c"
	}
	if filepath.Ext(output) != ".c" {
		output += ".c"
	}

	if err := os.WriteFile(output, []byte(code), 0o600); err != nil {
		return fmt.Errorf("failed to write C source: %w", err)
	}

	fmt.Printf("[Success] Generated C source: %s\n", output)
	return nil
}

// Compile compiles the given source files to a C executable.
func Compile(files []*source.File, outputPath string, opts Options) error {
	if len(files) == 0 {
		return fmt.Errorf("no source files to compile")
	}

	return compileWithImports(files, outputPath, opts)
}

func compileWithImports(files []*source.File, outputPath string, opts Options) error {
	if len(files) == 0 {
		return fmt.Errorf("no source files to compile")
	}

	resolver := fs.NewResolver(opts.LibPaths)

	var allStmts []ast.Stmt
	var allFiles []*source.File
	var cumOffset uint32
	seen := map[string]bool{}

	queue := files
	for len(queue) > 0 {
		file := queue[0]
		queue = queue[1:]

		if seen[file.Name()] {
			continue
		}
		seen[file.Name()] = true

		diag := diagnostics.New(file)
		program, err := parseSource(file, diag)
		if err != nil {
			return fmt.Errorf("%s: %w", file.Name(), err)
		}

		if len(allFiles) > 0 {
			cumOffset++
			ast.AdjustPositions(program, cumOffset)
		}

		for _, stmt := range program.Statements {
			if imp, ok := stmt.(*ast.ImportStmt); ok {
				resolvedPath, err := resolver.Resolve(imp.Path.Value, filepath.Dir(file.Name()))
				if err != nil {
					return err
				}
				if !seen[resolvedPath] {
					data, err := os.ReadFile(resolvedPath)
					if err != nil {
						return fmt.Errorf("reading imported file %q: %w", resolvedPath, err)
					}
					queue = append(queue, source.New(resolvedPath, data))
				}
			}
		}

		allStmts = append(allStmts, program.Statements...)
		allFiles = append(allFiles, file)
		cumOffset += file.Len()
	}

	merged := mergeSources(allFiles)
	mergedProgram := &ast.Program{Statements: allStmts}
	diag := diagnostics.New(merged)

	analyzer := sema.New(diag)
	if err := analyzer.Analyze(mergedProgram); err != nil {
		return err
	}

	optimizer.Optimize(mergedProgram)
	cCode, err := transpileToC(mergedProgram)
	if err != nil {
		return err
	}

	if opts.EmitC {
		return writeCOutput(cCode, outputPath)
	}

	exeName := resolveExeName(outputPath)
	tmpPath, err := writeToTempFile(cCode)
	if err != nil {
		return err
	}
	defer func() {
		os.Remove(tmpPath)
	}()

	return runCompiler([]string{tmpPath}, exeName, opts)
}

func compileSingle(file *source.File, outputPath string, opts Options) error {
	diag := diagnostics.New(file)

	program, err := parseSource(file, diag)
	if err != nil {
		return err
	}

	analyzer := sema.New(diag)

	if err := analyzer.Analyze(program); err != nil {
		return err
	}

	optimizer.Optimize(program)
	cCode, err := transpileToC(program)
	if err != nil {
		return err
	}

	if opts.EmitC {
		return writeCOutput(cCode, outputPath)
	}

	exeName := resolveExeName(outputPath)

	tmpPath, err := writeToTempFile(cCode)
	if err != nil {
		return err
	}
	defer func(name string) {
		if err := os.Remove(name); err != nil {
			fmt.Printf("warning: failed to remove temp file %s: %v\n", name, err)
		}
	}(tmpPath)

	return runCompiler([]string{tmpPath}, exeName, opts)
}

func compileMulti(files []*source.File, outputPath string, opts Options) error {
	resolver := fs.NewResolver(opts.LibPaths)

	var allStmts []ast.Stmt
	var allFiles []*source.File
	var cumOffset uint32
	seen := map[string]bool{}

	queue := files
	for len(queue) > 0 {
		file := queue[0]
		queue = queue[1:]

		if seen[file.Name()] {
			continue
		}
		seen[file.Name()] = true

		diag := diagnostics.New(file)
		program, err := parseSource(file, diag)
		if err != nil {
			return fmt.Errorf("%s: %w", file.Name(), err)
		}

		if len(allFiles) > 0 {
			cumOffset++
			ast.AdjustPositions(program, cumOffset)
		}

		for _, stmt := range program.Statements {
			if imp, ok := stmt.(*ast.ImportStmt); ok {
				resolvedPath, err := resolver.Resolve(imp.Path.Value, filepath.Dir(file.Name()))
				if err != nil {
					return err
				}
				if !seen[resolvedPath] {
					data, err := os.ReadFile(resolvedPath)
					if err != nil {
						return fmt.Errorf("reading imported file %q: %w", resolvedPath, err)
					}
					queue = append(queue, source.New(resolvedPath, data))
				}
			}
		}

		allStmts = append(allStmts, program.Statements...)
		allFiles = append(allFiles, file)
		cumOffset += file.Len()
	}

	merged := mergeSources(allFiles)
	mergedProgram := &ast.Program{Statements: allStmts}
	diag := diagnostics.New(merged)

	analyzer := sema.New(diag)
	if err := analyzer.Analyze(mergedProgram); err != nil {
		return err
	}

	optimizer.Optimize(mergedProgram)
	cCode, err := transpileToC(mergedProgram)
	if err != nil {
		return err
	}

	if opts.EmitC {
		return writeCOutput(cCode, outputPath)
	}

	exeName := resolveExeName(outputPath)
	tmpPath, err := writeToTempFile(cCode)
	if err != nil {
		return err
	}
	defer func() {
		os.Remove(tmpPath)
	}()

	return runCompiler([]string{tmpPath}, exeName, opts)
}

func mergeSources(files []*source.File) *source.File {
	if len(files) == 1 {
		return files[0]
	}

	var names []string
	var totalLen int
	for _, f := range files {
		names = append(names, f.Name())
		totalLen += int(f.Len())
	}

	text := make([]byte, 0, totalLen+len(files))
	for i, f := range files {
		if i > 0 {
			text = append(text, '\n')
		}
		text = append(text, f.Slice(0, f.Len())...)
	}

	displayName := strings.Join(names, " + ")
	return source.New(displayName, text)
}

func parseSource(file *source.File, diag *diagnostics.Engine) (*ast.Program, error) {
	tokens := lexer.New(file, diag).Tokenize()
	if err := diag.Err(); err != nil {
		return nil, err
	}

	sourceParser := parser.New(string(file.Slice(0, file.Len())), tokens, diag)

	program := sourceParser.ParseProgram()

	if err := sourceParser.Err(); err != nil {
		return nil, err
	}

	return program, diag.Err()
}

func transpileToC(program *ast.Program) (string, error) {
	tx := c.New()
	return tx.Transpile(program)
}

func resolveExeName(output string) string {
	if output == "" {
		if runtime.GOOS == "windows" {
			return "output.exe"
		}
		return "output"
	}

	output = strings.TrimSuffix(output, ".c")

	if runtime.GOOS == "windows" && filepath.Ext(output) != ".exe" {
		output += ".exe"
	}

	return output
}

func writeToTempFile(content string) (string, error) {
	f, err := os.CreateTemp("", "azin_*.c")
	if err != nil {
		return "", fmt.Errorf("failed to create temp source: %w", err)
	}

	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("warning: failed to close temp source: %v\n", err)
		}
	}()

	if _, err := f.WriteString(content); err != nil {
		return "", fmt.Errorf("failed to write temp source: %w", err)
	}

	return f.Name(), nil
}
