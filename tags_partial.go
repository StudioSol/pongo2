package pongo2

type tagPartialNode struct {
	filename string
}

func (node *tagPartialNode) Execute(ctx *ExecutionContext, writer TemplateWriter) *Error {
	return nil
}

func tagPartialParser(doc *Parser, start *Token, arguments *Parser) (INodeTag, *Error) {
	extends_node := &tagPartialNode{}

	if doc.template.level > 1 {
		return nil, arguments.Error("The 'partial' tag can only defined on root level.", start)
	}

	if doc.template.parent != nil {
		// Already one parent
		return nil, arguments.Error("This template has already one parent.", start)
	}
	var partial_template *Template
	var err error
	if filename_token := arguments.MatchType(TokenString); filename_token != nil {
		//partial template
		partial_template, err = doc.template.set.FromFile(doc.template.set.resolveFilename(doc.template, filename_token.Val))
		if err != nil {
			return nil, err.(*Error).updateFromTokenIfNeeded(doc.template, filename_token)
		}
	} else {
		return nil, arguments.Error("Tag 'partial' requires a template filename as string.", nil)
	}

	if arguments.Match(TokenIdentifier, "extends") != nil {
		if arguments.Match(TokenSymbol, "=") == nil {
			return nil, arguments.Error("Expected '='.", nil)
		}
		if filename_token := arguments.MatchType(TokenString); filename_token != nil {

			// Get parent's filename
			parent_filename := doc.template.set.resolveFilename(doc.template, filename_token.Val)

			// Parse the parent
			parent_template, err := doc.template.set.FromFile(parent_filename)
			if err != nil {
				return nil, err.(*Error)
			}

			// Keep track of things
			if partial_template == nil {
				parent_template.child = doc.template
			} else {
				parent_template.child = partial_template
			}

			doc.template.parent = parent_template
			extends_node.filename = parent_filename

		} else {
			return nil, arguments.Error("Tag 'partial' requires a template filename as string.", nil)
		}
		if arguments.Remaining() > 0 {
			return nil, arguments.Error("Tag 'partial' does only take 2 argument.", nil)
		}
	}

	return extends_node, nil
}

func init() {
	RegisterTag("partial", tagPartialParser)
}
