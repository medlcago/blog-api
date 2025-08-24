package cache

type KeyBuilder struct {
	prefix    string
	separator string
}

func (b *KeyBuilder) SetPrefix(prefix string) *KeyBuilder {
	b.prefix = prefix
	return b
}

func (b *KeyBuilder) SetSeparator(separator string) *KeyBuilder {
	b.separator = separator
	return b
}

func NewKeyBuilder(prefix, separator string) *KeyBuilder {
	if separator == "" {
		separator = ":"
	}
	return &KeyBuilder{prefix: prefix, separator: separator}
}

func (b *KeyBuilder) Build(key string) string {
	if key == "" {
		return ""
	}

	if b.prefix == "" {
		return key
	}

	return b.prefix + "_" + key
}
