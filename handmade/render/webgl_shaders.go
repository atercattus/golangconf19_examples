// +build js

package render

const (
	vertexShaderCode = `
attribute vec2 in_Position;
attribute vec2 in_TexCoords;
attribute vec4 in_Color;

uniform mat3 PV;

varying vec4 var_Color;
varying vec2 var_TexCoords;

void main(void) {
	vec3 matr = PV * vec3(in_Position, 1.0);
	gl_Position = vec4(matr.xy, 0, 1);
	var_Color = in_Color;
	var_TexCoords = in_TexCoords;
}
`

	fragmentShaderCode = `
#ifdef GL_ES
#define LOWP lowp
precision mediump float;
#else
#define LOWP
#endif

varying vec4 var_Color;
varying vec2 var_TexCoords;

uniform sampler2D uSampler;

void main(void) {
	gl_FragColor = texture2D(uSampler, var_TexCoords) * var_Color;
}
`
)
