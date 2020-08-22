package allegro

// #include <allegro5/allegro.h>
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

/* Shader variable names */
const (
	SHADER_VAR_COLOR           = "al_color"
	SHADER_VAR_POS             = "al_pos"
	SHADER_VAR_PROJVIEW_MATRIX = "al_projview_matrix"
	SHADER_VAR_TEX             = "al_tex"
	SHADER_VAR_TEXCOORD        = "al_texcoord"
	SHADER_VAR_TEX_MATRIX      = "al_tex_matrix"
	SHADER_VAR_USER_ATTR       = "al_user_attr_"
	SHADER_VAR_USE_TEX         = "al_use_tex"
	SHADER_VAR_USE_TEX_MATRIX  = "al_use_tex_matrix"
)

var (
	ShaderIsNull = errors.New("shader is null")
)

type ShaderType int

const (
	VERTEX_SHADER ShaderType = C.ALLEGRO_VERTEX_SHADER
	PIXEL_SHADER             = C.ALLEGRO_PIXEL_SHADER
)

type ShaderPlatform int

const (
	SHADER_AUTO ShaderPlatform = C.ALLEGRO_SHADER_AUTO
	SHADER_GLSL                = C.ALLEGRO_SHADER_GLSL
	SHADER_HLSL                = C.ALLEGRO_SHADER_HLSL
)

type Shader C.struct_ALLEGRO_SHADER

func CreateShader(platform ShaderPlatform) (*Shader, error) {
	s := C.al_create_shader(C.ALLEGRO_SHADER_PLATFORM(platform))
	if s == nil {
		return nil, errors.New("failed to create shader")
	}
	return (*Shader)(s), nil
}

func (s *Shader) AttachSource(stype ShaderType, source string) error {
	if s == nil {
		return ShaderIsNull
	}

	source_ := C.CString(source)
	defer freeString(source_)

	ok := C.al_attach_shader_source((*C.ALLEGRO_SHADER)(s), C.ALLEGRO_SHADER_TYPE(stype), source_)
	if !ok {
		return errors.New("failed to attach shader source")
	}

	return nil
}

func (s *Shader) AttachSourceFile(stype ShaderType, filename string) error {
	if s == nil {
		return ShaderIsNull
	}

	filename_ := C.CString(filename)
	defer freeString(filename_)

	ok := C.al_attach_shader_source_file((*C.ALLEGRO_SHADER)(s), C.ALLEGRO_SHADER_TYPE(stype), filename_)
	if !ok {
		return fmt.Errorf("failed to attach shader source file \"%s\"", filename)
	}

	return nil
}

func (s *Shader) Build() error {
	if s == nil {
		return ShaderIsNull
	}

	ok := C.al_build_shader((*C.ALLEGRO_SHADER)(s))
	if !ok {
		return errors.New("failed to build shader")
	}

	return nil
}

func (s *Shader) Log() (string, error) {
	if s == nil {
		return "", ShaderIsNull
	}

	log := C.al_get_shader_log((*C.ALLEGRO_SHADER)(s))
	return C.GoString(log), nil
}

func (s *Shader) Platform() (ShaderPlatform, error) {
	if s == nil {
		return 0, ShaderIsNull
	}

	p := C.al_get_shader_platform((*C.ALLEGRO_SHADER)(s))
	return ShaderPlatform(p), nil
}

func UseShader(s *Shader) error {
	ok := C.al_use_shader((*C.ALLEGRO_SHADER)(s))
	if !ok {
		return errors.New("failed to use shader")
	}

	return nil
}

func (s *Shader) Destroy() {
	C.al_destroy_shader((*C.ALLEGRO_SHADER)(s))
}

func SetShaderSampler(name string, bmp *Bitmap, unit int) error {
	if bmp == nil {
		return BitmapIsNull
	}

	name_ := C.CString(name)
	defer freeString(name_)

	ok := C.al_set_shader_sampler(name_, (*C.ALLEGRO_BITMAP)(bmp), C.int(unit))
	if !ok {
		return fmt.Errorf("failed to set shader sampler for \"%s\"", name)
	}

	return nil
}

func SetShaderMatrix(name string, matrix *Transform) error {
	name_ := C.CString(name)
	defer freeString(name_)

	ok := C.al_set_shader_matrix(name_, (*C.ALLEGRO_TRANSFORM)(matrix))
	if !ok {
		return fmt.Errorf("failed to set shader matrix for \"%s\"", name)
	}

	return nil
}

func SetShaderInt(name string, i int) error {
	name_ := C.CString(name)
	defer freeString(name_)

	ok := C.al_set_shader_int(name_, C.int(i))
	if !ok {
		return fmt.Errorf("failed to set shader int for \"%s\"", name)
	}

	return nil
}

func SetShaderFloat(name string, f float32) error {
	name_ := C.CString(name)
	defer freeString(name_)

	ok := C.al_set_shader_float(name_, C.float(f))
	if !ok {
		return fmt.Errorf("failed to set shader float for \"%s\"", name)
	}

	return nil
}

func SetShaderIntVector(name string, i [][]int) error {
	name_ := C.CString(name)
	defer freeString(name_)

	var ok C.bool

	if len(i) == 0 {
		ok = C.al_set_shader_int_vector(name_, 0, nil, 0)
	} else {
		elems := len(i)
		components := len(i[0])

		cmem := malloc(C.size_t(elems*components) * C.size_t(unsafe.Sizeof(C.int(0))))
		if cmem == nil {
			return errors.New("failed to allocate int vector")
		}
		defer free(cmem)

		garr := (*[1<<30 - 1]C.int)(cmem)
		idx := 0
		for _, slice := range i {
			for _, elem := range slice {
				garr[idx] = C.int(elem)
				idx++
			}
		}

		ok = C.al_set_shader_int_vector(name_, C.int(components), (*C.int)(cmem), C.int(elems))
	}

	if !ok {
		return fmt.Errorf("failed to set int vector for \"%s\"", name)
	}

	return nil
}

func SetShaderFloatVector(name string, f [][]float32) error {
	name_ := C.CString(name)
	defer freeString(name_)

	var ok C.bool

	if len(f) == 0 {
		ok = C.al_set_shader_float_vector(name_, 0, nil, 0)
	} else {
		elems := len(f)
		components := len(f[0])

		cmem := malloc(C.size_t(elems*components) * C.size_t(unsafe.Sizeof(C.float(0.0))))
		if cmem == nil {
			return errors.New("failed to allocate float vector")
		}
		defer free(cmem)

		garr := (*[1<<30 - 1]C.float)(cmem)
		idx := 0
		for _, slice := range f {
			for _, elem := range slice {
				garr[idx] = C.float(elem)
				idx++
			}
		}

		ok = C.al_set_shader_float_vector(name_, C.int(components), (*C.float)(cmem), C.int(elems))
	}

	if !ok {
		return fmt.Errorf("failed to set float vector for \"%s\"", name)
	}

	return nil
}

func SetShaderBool(name string, b bool) error {
	name_ := C.CString(name)
	defer freeString(name_)

	ok := C.al_set_shader_bool(name_, C.bool(b))
	if !ok {
		return fmt.Errorf("failed to set shader bool for \"%s\"", name)
	}

	return nil
}

func DefaultShaderSource(platform ShaderPlatform, typ ShaderType) string {
	cstr := C.al_get_default_shader_source(C.ALLEGRO_SHADER_PLATFORM(platform), C.ALLEGRO_SHADER_TYPE(typ))
	return C.GoString(cstr)
}
