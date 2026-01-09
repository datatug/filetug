package chroma2tcell

import (
	"fmt"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

func TestColorizeYAMLForTview(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s, err := ColorizeYAMLForTview("", lexers.Get)
		assert.NoError(t, err)
		assert.Equal(t, "", s)
	})

	t.Run("simple_yaml", func(t *testing.T) {
		s, err := ColorizeYAMLForTview("key: value", lexers.Get)
		assert.NoError(t, err)
		assert.Contains(t, s, "[")
		assert.Contains(t, s, "key")
		assert.Contains(t, s, "value")
	})

	t.Run("lexer_not_found", func(t *testing.T) {
		s, err := ColorizeYAMLForTview("key: value", func(string) chroma.Lexer { return nil })
		assert.NoError(t, err)
		assert.Contains(t, s, "key: value")
	})
}

func TestColorize(t *testing.T) {
	t.Run("invalid_lexer", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		_, _ = Colorize("text", "dracula", nil)
	})

	t.Run("with_lexer", func(t *testing.T) {
		lexer := lexers.Get("go")
		s, err := Colorize("package main", "dracula", lexer)
		assert.NoError(t, err)
		assert.Contains(t, s, "package")
	})

	t.Run("token_with_no_color", func(t *testing.T) {
		lexer := lexers.Get("go")
		// The style "swapoff" might not have colors for everything
		s, err := Colorize("package main", "swapoff", lexer)
		assert.NoError(t, err)
		assert.Contains(t, s, "package")
	})

	t.Run("tokenise_error", func(t *testing.T) {
		lexer := &mockLexer{err: fmt.Errorf("tokenise error")}
		_, err := Colorize("text", "dracula", lexer)
		assert.Error(t, err)
	})

	t.Run("zero_color", func(t *testing.T) {
		lexer := &mockLexer{
			tokens: []chroma.Token{
				{Type: chroma.TokenType(-1), Value: "plain text"},
			},
		}
		// Register a style that has No color for our custom token type
		builder := styles.Fallback.Builder()
		// TokenType(-1) should have no entry in this builder, thus it should be zero.
		customStyle, _ := builder.Build()
		// Let's force it by NOT adding any entries.
		// Actually, we need to make sure the style doesn't have a background color or anything that makes color.IsZero() false.

		s, err := Colorize("plain text", customStyle.Name, lexer)
		assert.NoError(t, err)
		t.Logf("Result: %q", s)
	})
}

type mockLexer struct {
	tokens []chroma.Token
	err    error
}

func (m *mockLexer) Tokenise(options *chroma.TokeniseOptions, text string) (chroma.Iterator, error) {
	_, _ = options, text
	if m.err != nil {
		return nil, m.err
	}
	return chroma.Literator(m.tokens...), nil
}

func (m *mockLexer) Config() *chroma.Config {
	return nil
}

func (m *mockLexer) SetRegistry(*chroma.LexerRegistry) chroma.Lexer {
	return m
}

func (m *mockLexer) SetAnalyser(analyser func(text string) float32) chroma.Lexer {
	_ = analyser
	return m
}

func (m *mockLexer) AnalyseText(string) float32 {
	return 0
}
