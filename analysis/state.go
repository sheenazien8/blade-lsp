package analysis

import (
	"github.com/sheenazien8/blade-lsp/lsp"
	"fmt"
	"regexp"
	"strings"
)

type State struct {
	Documents map[string]string
}

func NewState() State {
	return State{Documents: map[string]string{}}
}

func getDiagnosticsForFile(text string) []lsp.Diagnostic {
	diagnostics := []lsp.Diagnostic{}
	for row, line := range strings.Split(text, "\n") {
		if strings.Contains(line, "VS Code") {
			idx := strings.Index(line, "VS Code")
			diagnostics = append(diagnostics, lsp.Diagnostic{
				Range:    LineRange(row, idx, idx+len("VS Code")),
				Severity: 1,
				Source:   "Common Sense",
				Message:  "Please make sure we use good language in this video",
			})
		}

		if strings.Contains(line, "Neovim") {
			idx := strings.Index(line, "Neovim")
			diagnostics = append(diagnostics, lsp.Diagnostic{
				Range:    LineRange(row, idx, idx+len("Neovim")),
				Severity: 2,
				Source:   "Common Sense",
				Message:  "Great choice :)",
			})

		}
	}

	return diagnostics
}

func (s *State) OpenDocument(uri, text string) []lsp.Diagnostic {
	s.Documents[uri] = text

	return getDiagnosticsForFile(text)
}

func (s *State) UpdateDocument(uri, text string) []lsp.Diagnostic {
	s.Documents[uri] = text

	return getDiagnosticsForFile(text)
}

func (s *State) Hover(id int, uri string, position lsp.Position) lsp.HoverResponse {
	document := s.Documents[uri]

	return lsp.HoverResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: lsp.HoverResult{
			Contents: fmt.Sprintf("File: %s, Characters: %d", uri, len(document)),
		},
	}
}

func (s *State) Definition(id int, uri string, position lsp.Position) lsp.DefinitionResponse {
	return lsp.DefinitionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: lsp.Location{
			URI: uri,
			Range: lsp.Range{
				Start: lsp.Position{
					Line:      position.Line - 1,
					Character: 0,
				},
				End: lsp.Position{
					Line:      position.Line - 1,
					Character: 0,
				},
			},
		},
	}
}

func (s *State) TextDocumentCodeAction(id int, uri string) lsp.TextDocumentCodeActionResponse {
	text := s.Documents[uri]

	actions := []lsp.CodeAction{}
	for row, line := range strings.Split(text, "\n") {
		idx := strings.Index(line, "VS Code")
		if idx >= 0 {
			replaceChange := map[string][]lsp.TextEdit{}
			replaceChange[uri] = []lsp.TextEdit{
				{
					Range:   LineRange(row, idx, idx+len("VS Code")),
					NewText: "Neovim",
				},
			}

			actions = append(actions, lsp.CodeAction{
				Title: "Replace VS C*de with a superior editor",
				Edit:  &lsp.WorkspaceEdit{Changes: replaceChange},
			})

			censorChange := map[string][]lsp.TextEdit{}
			censorChange[uri] = []lsp.TextEdit{
				{
					Range:   LineRange(row, idx, idx+len("VS Code")),
					NewText: "VS C*de",
				},
			}

			actions = append(actions, lsp.CodeAction{
				Title: "Censor to VS C*de",
				Edit:  &lsp.WorkspaceEdit{Changes: censorChange},
			})
		}
	}

	response := lsp.TextDocumentCodeActionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: actions,
	}

	return response
}

func (s *State) TextDocumentCompletion(id int, uri string, position lsp.Position) lsp.CompletionResponse {
	document := s.Documents[uri]
	
	// Get the current line and character position
	lines := strings.Split(document, "\n")
	if position.Line >= len(lines) {
		return s.getDefaultCompletions(id)
	}
	
	currentLine := lines[position.Line]
	if position.Character > len(currentLine) {
		return s.getDefaultCompletions(id)
	}
	
	// Check if we're in a Blade variable context {{ $
	lineUpToCursor := currentLine[:position.Character]
	if strings.Contains(lineUpToCursor, "{{ $") || strings.HasSuffix(lineUpToCursor, "{{ $") {
		return s.getVariableCompletions(id, document)
	}
	
	// Check if we're typing @ for directives
	if strings.HasSuffix(lineUpToCursor, "@") || strings.Contains(lineUpToCursor, "@") {
		return s.getDirectiveCompletions(id)
	}
	
	return s.getDefaultCompletions(id)
}

func (s *State) getVariableCompletions(id int, document string) lsp.CompletionResponse {
	variables := s.extractVariablesFromDocument(document)
	
	items := []lsp.CompletionItem{}
	for _, variable := range variables {
		items = append(items, lsp.CompletionItem{
			Label:            variable,
			Kind:             6, // Variable
			Detail:           "Blade variable",
			Documentation:    fmt.Sprintf("Variable: $%s", variable),
			InsertText:       variable,
			InsertTextFormat: 1, // Plain text
		})
	}
	
	commonVars := []string{"user", "errors", "request", "session", "config", "app", "auth"}
	for _, variable := range commonVars {
		items = append(items, lsp.CompletionItem{
			Label:            variable,
			Kind:             6, // Variable
			Detail:           "Common Laravel variable",
			Documentation:    fmt.Sprintf("Common Laravel variable: $%s", variable),
			InsertText:       variable,
			InsertTextFormat: 1, // Plain text
		})
	}
	
	return lsp.CompletionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: items,
	}
}

func (s *State) extractVariablesFromDocument(document string) []string {
	variables := make(map[string]bool)
	
	re := regexp.MustCompile(`\{\{\s*\$(\w+)`)
	matches := re.FindAllStringSubmatch(document, -1)
	for _, match := range matches {
		if len(match) > 1 {
			variables[match[1]] = true
		}
	}
	
	foreachRe := regexp.MustCompile(`@foreach\s*\(\s*\$(\w+)\s+as\s+\$(\w+)\s*\)`)
	foreachMatches := foreachRe.FindAllStringSubmatch(document, -1)
	for _, match := range foreachMatches {
		if len(match) > 2 {
			variables[match[1]] = true
			variables[match[2]] = true
		}
	}
	
	ifRe := regexp.MustCompile(`@if\s*\(\s*\$(\w+)`)
	ifMatches := ifRe.FindAllStringSubmatch(document, -1)
	for _, match := range ifMatches {
		if len(match) > 1 {
			variables[match[1]] = true
		}
	}
	
	result := make([]string, 0, len(variables))
	for variable := range variables {
		result = append(result, variable)
	}
	
	return result
}

func (s *State) getDirectiveCompletions(id int) lsp.CompletionResponse {
	// Laravel Blade completion items
	items := []lsp.CompletionItem{
		// Blade directives
		{
			Label:         "@if",
			Detail:        "Blade conditional directive",
			Documentation: "Conditional statement: @if($condition) ... @endif",
		},
		{
			Label:         "@foreach",
			Detail:        "Blade loop directive",
			Documentation: "Loop through arrays: @foreach($items as $item) ... @endforeach",
		},
		{
			Label:         "@for",
			Detail:        "Blade for loop directive",
			Documentation: "For loop: @for($i = 0; $i < 10; $i++) ... @endfor",
		},
		{
			Label:         "@while",
			Detail:        "Blade while loop directive",
			Documentation: "While loop: @while($condition) ... @endwhile",
		},
		{
			Label:         "@extends",
			Detail:        "Blade template inheritance",
			Documentation: "Extend a parent template: @extends('layouts.app')",
		},
		{
			Label:         "@section",
			Detail:        "Blade template section",
			Documentation: "Define a section: @section('content') ... @endsection",
		},
		{
			Label:         "@yield",
			Detail:        "Blade template yield",
			Documentation: "Output a section: @yield('content')",
		},
		{
			Label:         "@include",
			Detail:        "Blade template include",
			Documentation: "Include another template: @include('partials.header')",
		},
		{
			Label:         "@csrf",
			Detail:        "Laravel CSRF token",
			Documentation: "Generate CSRF token field for forms",
		},
		{
			Label:         "@method",
			Detail:        "Laravel HTTP method spoofing",
			Documentation: "Spoof HTTP methods: @method('PUT')",
		},
		{
			Label:         "@auth",
			Detail:        "Laravel authentication check",
			Documentation: "Check if user is authenticated: @auth ... @endauth",
		},
		{
			Label:         "@guest",
			Detail:        "Laravel guest check",
			Documentation: "Check if user is guest: @guest ... @endguest",
		},
		{
			Label:         "@can",
			Detail:        "Laravel authorization check",
			Documentation: "Check user permissions: @can('update', $post) ... @endcan",
		},
		{
			Label:         "@error",
			Detail:        "Laravel validation error",
			Documentation: "Display validation errors: @error('field') ... @enderror",
		},
		{
			Label:         "@empty",
			Detail:        "Blade empty check",
			Documentation: "Check if variable is empty: @empty($variable) ... @endempty",
		},
		{
			Label:         "@isset",
			Detail:        "Blade isset check",
			Documentation: "Check if variable is set: @isset($variable) ... @endisset",
		},
		{
			Label:         "@switch",
			Detail:        "Blade switch statement",
			Documentation: "Switch statement: @switch($variable) @case(1) ... @endswitch",
		},
		{
			Label:         "@php",
			Detail:        "Blade PHP block",
			Documentation: "Raw PHP code: @php ... @endphp",
		},
		{
			Label:         "@json",
			Detail:        "Blade JSON output",
			Documentation: "Output JSON: @json($array)",
		},
		{
			Label:         "@component",
			Detail:        "Blade component",
			Documentation: "Use a component: @component('alert') ... @endcomponent",
		},
	}

	return lsp.CompletionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: items,
	}
}

func (s *State) getDefaultCompletions(id int) lsp.CompletionResponse {
	return s.getDirectiveCompletions(id)
}

func LineRange(line, start, end int) lsp.Range {
	return lsp.Range{
		Start: lsp.Position{
			Line:      line,
			Character: start,
		},
		End: lsp.Position{
			Line:      line,
			Character: end,
		},
	}
}
