#version 150

in vec3 Color;
in vec2 Texcoord;
out vec4 outColor;
uniform sampler2D tex0;
uniform sampler2D tex1;
uniform float time;

void main()
{
  float factor = (sin(time * 3.0) + 1.0) / 2.0;
  outColor = mix(texture(tex0, Texcoord), texture(tex1, Texcoord), factor) * vec4(Color, 1.0);
}
