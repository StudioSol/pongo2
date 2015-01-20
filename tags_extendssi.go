package pongo2

type tagExtendSSINode struct {
	filename string
}

func (node *tagExtendSSINode) Execute(ctx *ExecutionContext, writer TemplateWriter) *Error {
	return nil
}

func tagExtendSSIParser(doc *Parser, start *Token, arguments *Parser) (INodeTag, *Error) {
	extends_node := &tagExtendSSINode{}

	if doc.template.level > 1 {
		return nil, arguments.Error("The 'extends' tag can only defined on root level.", start)
	}

	if doc.template.parent != nil {
		// Already one parent
		return nil, arguments.Error("This template has already one parent.", start)
	}

	if filename_token := arguments.MatchType(TokenString); filename_token != nil {
		// prepared, static template

		// Get parent's filename
		parent_filename := doc.template.set.resolveFilename(doc.template, filename_token.Val)

		// Parse the parent
		parent_template, err := doc.template.set.FromFile(parent_filename)
		if err != nil {
			return nil, err.(*Error)
		}

		// Keep track of things
		parent_template.child = doc.template
		doc.template.parent = parent_template
		extends_node.filename = parent_filename
	} else {
		return nil, arguments.Error("Tag 'extendssi' requires a template filename as string.", nil)
	}

	//Change parent child
	if filename_token := arguments.MatchType(TokenString); filename_token != nil {
		ssi_template, err := doc.template.set.FromFile(doc.template.set.resolveFilename(doc.template, filename_token.Val))
		if err != nil {
			return nil, err.(*Error).updateFromTokenIfNeeded(doc.template, filename_token)
		}
		doc.template.parent.child = ssi_template
	} else {
		return nil, arguments.Error("Tag 'extendssi' requires a template ssi filename as string.", nil)
	}

	if arguments.Remaining() > 0 {
		return nil, arguments.Error("Tag 'extendssi' does only take 2 argument.", nil)
	}

	return extends_node, nil
}

func init() {
	RegisterTag("extendssi", tagExtendSSIParser)
}
