package main

import (
	"testing"

	"rush/interpreter"
	"rush/lexer"
	"rush/parser"
)

func TestRegexpBasics(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`pattern = Regexp("\\d+"); pattern.pattern`,
			`\d+`,
		},
		{
			`pattern = Regexp("\\d+"); pattern`,
			`/\d+/`,
		},
		{
			`Regexp("invalid[")`,
			"RuntimeError: invalid regular expression: error parsing regexp: missing closing ]: `[`",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		env := interpreter.NewEnvironment()

		evaluated := interpreter.Eval(program, env)
		if evaluated == nil {
			t.Errorf("Eval returned nil")
			continue
		}

		if evaluated.Inspect() != tt.expected {
			t.Errorf("wrong result. got=%q, want=%q", evaluated.Inspect(), tt.expected)
		}
	}
}

func TestStringMethodsWithRegexp(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`text = "hello 123 world 456"; pattern = Regexp("\\d+"); text.match(pattern)`,
			`[123, 456]`,
		},
		{
			`text = "hello world"; pattern = Regexp("\\d+"); text.match(pattern)`,
			`[]`,
		},
		{
			`text = "hello 123 world"; pattern = Regexp("\\d+"); text.matches?(pattern)`,
			`true`,
		},
		{
			`text = "hello world"; pattern = Regexp("\\d+"); text.matches?(pattern)`,
			`false`,
		},
		{
			`text = "hello 123 world 456"; pattern = Regexp("\\d+"); text.replace(pattern, "X")`,
			`hello X world X`,
		},
		{
			`text = "a,b;c:d"; pattern = Regexp("[,;:]"); text.split(pattern)`,
			`[a, b, c, d]`,
		},
		{
			`text = "hello world"; text.match("\\w+")`,
			`[hello, world]`,
		},
		{
			`text = "hello world"; text.matches?("\\w+")`,
			`true`,
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		env := interpreter.NewEnvironment()

		evaluated := interpreter.Eval(program, env)
		if evaluated == nil {
			t.Errorf("Eval returned nil")
			continue
		}

		if evaluated.Inspect() != tt.expected {
			t.Errorf("wrong result for %q. got=%q, want=%q", tt.input, evaluated.Inspect(), tt.expected)
		}
	}
}

func TestRegexpObjectMethods(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`pattern = Regexp("\\d+"); pattern.matches?("hello 123")`,
			`true`,
		},
		{
			`pattern = Regexp("\\d+"); pattern.matches?("hello world")`,
			`false`,
		},
		{
			`pattern = Regexp("\\d+"); pattern.find_all("hello 123 world 456")`,
			`[123, 456]`,
		},
		{
			`pattern = Regexp("\\d+"); pattern.find_first("hello 123 world 456")`,
			`123`,
		},
		{
			`pattern = Regexp("\\d+"); pattern.find_first("hello world")`,
			`null`,
		},
		{
			`pattern = Regexp("\\d+"); pattern.replace("hello 123 world 456", "X")`,
			`hello X world X`,
		},
		{
			`pattern = Regexp("(?i)hello"); pattern.find_all("HELLO hello HeLLo")`,
			`[HELLO, hello, HeLLo]`,
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		env := interpreter.NewEnvironment()

		evaluated := interpreter.Eval(program, env)
		if evaluated == nil {
			t.Errorf("Eval returned nil")
			continue
		}

		if evaluated.Inspect() != tt.expected {
			t.Errorf("wrong result for %q. got=%q, want=%q", tt.input, evaluated.Inspect(), tt.expected)
		}
	}
}

func TestRegexpErrors(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`Regexp(123)`,
			"RuntimeError: argument to `Regexp` must be STRING, got INTEGER",
		},
		{
			`Regexp("\\d+", "extra")`,
			`RuntimeError: wrong number of arguments. got=2, want=1`,
		},
		{
			`pattern = Regexp("\\d+"); pattern.matches?()`,
			`RuntimeError: wrong number of arguments for matches?: want=1, got=0`,
		},
		{
			`pattern = Regexp("\\d+"); pattern.matches?(123)`,
			`RuntimeError: argument to matches? must be STRING, got INTEGER`,
		},
		{
			`pattern = Regexp("\\d+"); pattern.unknown_method("test")`,
			`RuntimeError: unknown property unknown_method for Regexp`,
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		env := interpreter.NewEnvironment()

		evaluated := interpreter.Eval(program, env)
		if evaluated == nil {
			t.Errorf("Eval returned nil")
			continue
		}

		if _, ok := evaluated.(*interpreter.Error); !ok {
			t.Errorf("expected error but got %T", evaluated)
			continue
		}

		if evaluated.Inspect() != tt.expected {
			t.Errorf("wrong error message. got=%q, want=%q", evaluated.Inspect(), tt.expected)
		}
	}
}