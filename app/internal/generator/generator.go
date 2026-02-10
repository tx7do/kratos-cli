package generator

type Generator struct {
	options GeneratorOptions
}

func NewGenerator() *Generator {
	return &Generator{
		options: GeneratorOptions{},
	}
}

func (g *Generator) GetOptions() GeneratorOptions {
	return g.options
}

func (g *Generator) SetOptions(options GeneratorOptions) {
	g.options = options
}

func (g *Generator) EditOption(o *Option) {
	if o == nil {
		return
	}

	for i, opt := range g.options {
		if opt.TableName == o.TableName {
			g.options[i] = o
			return
		}
	}
}

func (g *Generator) AddOption(o *Option) {
	if o == nil {
		return
	}

	if o.TableName == "" {
		return
	}

	o.ID = uint32(len(g.options) + 1)

	g.options = append(g.options, o)
}

func (g *Generator) CleanOptions() {
	g.options = GeneratorOptions{}
}

func (g *Generator) ValidateOptions() string {
	if len(g.options) == 0 {
		return "no tables selected"
	}

	for _, opt := range g.options {
		if opt.TableName == "" {
			return "table name cannot be empty"
		}
		if opt.Service == "" {
			return "service name cannot be empty"
		}
	}

	return ""
}

func (g *Generator) GetValidateOptions() GeneratorOptions {
	var options GeneratorOptions
	for _, opt := range g.options {
		if opt.TableName != "" &&
			opt.Service != "" &&
			!opt.Exclude {
			options = append(options, opt)
		}
	}
	return options
}
