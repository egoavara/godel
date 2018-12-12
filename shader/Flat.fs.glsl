
uniform vec4 BaseColorFactor;

#ifdef HAS_BASECOLORTEX
uniform sampler2D BaseColorTex;
#endif

in struct{
   vec3 position;
   vec2 texCoord_0;
   #ifdef HAS_NORMAL
       #ifdef HAS_TANGENT
           mat3 TBN;
       #else
           vec3 normal;
       #endif
   #endif
} fsout;



out vec4 outputColor;

void main() {
    vec4 c = BaseColorFactor;
    #ifdef HAS_BASECOLORTEX
        c = c * vec4(texture(BaseColorTex, fsout.texCoord_0.xy));
    #endif
    outputColor = c;
}
