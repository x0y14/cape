package parse

import "cape/c/parse/tokenize"

var token *tokenize.Token

func isEof() bool {
	return token.Kind == tokenize.Eof
}

func peekKind(kind tokenize.TokenKind) *tokenize.Token {
	if token.Kind == kind {
		return token
	}
	return nil
}

func peekNextKind(kind tokenize.TokenKind) *tokenize.Token {
	if token.Next.Kind == kind {
		return token.Next
	}
	return nil
}

func peekNextNextKind(kind tokenize.TokenKind) *tokenize.Token {
	if token.Next.Next.Kind == kind {
		return token.Next.Next
	}
	return nil
}
