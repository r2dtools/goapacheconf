package rawparser

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

const (
	whitespace         = `[^\S\r\n]+`
	ident              = `[\w]+`
	newLine            = `[\r\n]+`
	comment            = `(?:#)[^\n]*`
	blockStartOpen     = `<`
	blockStartClose    = `>`
	blockEnd           = `</\s*` + ident + `\s*>`
	expression         = `[^#\s<>]+`
	stringDoubleQuoted = `"(?:\\"|[^"])*"`
	stringSingleQuoted = `'(?:\\'|[^'])*'`
)

type Config struct {
	Entries []*Entry `@@*`
}

type Entry struct {
	StartNewLines  []string        `@NewLine*`
	Comment        *Comment        `( @@`
	Directive      *Directive      `| @@`
	BlockDirective *BlockDirective `| @@ )`
	EndNewLines    []string        `@NewLine*`
}

type Comment struct {
	Value string `@Comment`
}

type Directive struct {
	Identifier string   `@Ident`
	Values     []*Value `@@*`
}

type BlockDirective struct {
	Identifier string        `"<"@Ident`
	Parameters []*Value      `@@*">"`
	Content    *BlockContent `@@`
}

type BlockContent struct {
	Entries []*Entry `@@*`
}

type Value struct {
	Expression string `@Expression | @StringDoubleQuoted | @StringSingleQuoted`
}

type RawParser struct {
	participleParser *participle.Parser[Config]
}

func (p *RawParser) Parse(content string) (*Config, error) {
	return p.participleParser.ParseString("", content)
}

func GetRawParser() (*RawParser, error) {
	def := lexer.MustStateful(lexer.Rules{
		"Root": {
			{Name: `NewLine`, Pattern: newLine, Action: nil},
			{Name: `whitespace`, Pattern: whitespace, Action: nil},
			{Name: `Comment`, Pattern: comment, Action: nil},
			{Name: `Ident`, Pattern: ident, Action: lexer.Push("Directive")},
			{Name: "BlockStart", Pattern: blockStartOpen, Action: lexer.Push("BlockIdent")},
		},
		"Directive": {
			{Name: `whitespace`, Pattern: whitespace, Action: nil},
			{Name: `StringDoubleQuoted`, Pattern: stringDoubleQuoted, Action: nil},
			{Name: `StringSingleQuoted`, Pattern: stringSingleQuoted, Action: nil},
			{Name: "Expression", Pattern: expression, Action: nil},
			{Name: `Comment`, Pattern: comment, Action: nil},
			{Name: `NewLine`, Pattern: newLine, Action: lexer.Pop()},
			lexer.Return(),
		},
		"BlockIdent": {
			{Name: `whitespace`, Pattern: whitespace, Action: nil},
			{Name: `Ident`, Pattern: ident, Action: lexer.Push("BlockParams")},
			{Name: `NewLine`, Pattern: newLine, Action: lexer.Push("BlockContent")},
			lexer.Return(),
		},
		"BlockParams": {
			{Name: `whitespace`, Pattern: whitespace, Action: nil},
			{Name: `StringDoubleQuoted`, Pattern: stringDoubleQuoted, Action: nil},
			{Name: `StringSingleQuoted`, Pattern: stringSingleQuoted, Action: nil},
			{Name: "Expression", Pattern: expression, Action: nil},
			{Name: "BlockStartClose", Pattern: blockStartClose, Action: lexer.Pop()},
		},
		"BlockContent": {
			{Name: `whitespace`, Pattern: whitespace, Action: nil},
			{Name: `NewLine`, Pattern: newLine, Action: nil},
			{Name: `Comment`, Pattern: comment, Action: nil},
			{Name: `Ident`, Pattern: ident, Action: lexer.Push("Directive")},
			{Name: "blockEnd", Pattern: blockEnd, Action: lexer.Pop()},
			{Name: "BlockStart", Pattern: blockStartOpen, Action: lexer.Push("BlockIdent")},
		},
	})

	participleParser, err := participle.Build[Config](
		participle.Lexer(def),
		participle.UseLookahead(50),
	)

	if err != nil {
		return nil, err
	}

	parser := RawParser{
		participleParser: participleParser,
	}

	return &parser, nil
}
