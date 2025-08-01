//go:build !js

package gl

// Enum is equivalent to GLenum, and is normally used with one of the
// constants defined in this package.
type Enum uint32

// Attrib identifies the location of a specific attribute variable.
type Attrib struct {
	Value uint
}

// Program identifies a compiled shader program.
type Program struct {
	Value uint32
}

// Shader identifies a GLSL shader.
type Shader struct {
	Value uint32
}

// Buffer identifies a GL buffer object.
type Buffer struct {
	Value uint32
}

// Framebuffer identifies a GL framebuffer.
type Framebuffer struct {
	Value uint32
}

// A Renderbuffer is a GL object that holds an image in an internal format.
type Renderbuffer struct {
	Value uint32
}

// A Texture identifies a GL texture unit.
type Texture struct {
	Value uint32
}

// Uniform identifies the location of a specific uniform variable.
type Uniform struct {
	Value int32
}

var (
	NoAttrib       Attrib
	NoProgram      Program
	NoShader       Shader
	NoBuffer       Buffer
	NoFramebuffer  Framebuffer
	NoRenderbuffer Renderbuffer
	NoTexture      Texture
	NoUniform      Uniform
)

// Object is a generic interface for OpenGL objects
type Object interface {
	Identifier() Enum
	Name() uint32
}

// Implement Name() for the Object interface
func (p Program) Name() uint32 {
	return p.Value
}

func (s Shader) Name() uint32 {
	return s.Value
}

func (b Buffer) Name() uint32 {
	return b.Value
}

func (fb Framebuffer) Name() uint32 {
	return fb.Value
}

func (rb Renderbuffer) Name() uint32 {
	return rb.Value
}

func (t Texture) Name() uint32 {
	return t.Value
}
