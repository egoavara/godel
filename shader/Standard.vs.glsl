

//
// Reference > https://github.com/KhronosGroup/glTF-WebGL-PBR/blob/master/shaders/pbr-vert.glsl
//

uniform mat4 CameraMatrix;
uniform mat4 ModelMatrix;
uniform mat4 NormalMatrix;

layout (location = 0) in vec3 position;
#ifdef HAS_NORMAL
layout (location = 1) in vec3 normal;
#endif
#ifdef HAS_TANGENT
layout (location = 2) in vec3 tangent;
#endif
//
#ifdef HAS_COORD_0
layout (location = 4) in vec2 texCoord_0;
#endif
//#ifdef HAS_COORD_1
//layout (location = 5) in vec2 texCoord_1;
//#endif


out struct{
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

void main() {
    vec4 pos = ModelMatrix * vec4(position, 1);
    fsout.position = vec3(pos.xyz) / pos.w;
    // Normal,
    #ifdef HAS_NORMAL
        #ifdef HAS_TANGENT
            // if HAS_NORMAL and HAS_TANGENT
            vec3 normalW = normalize(vec3(NormalMatrix * vec4(normal.xyz, 0.0)));
            vec3 tangentW = normalize(vec3(ModelMatrix * vec4(tangent.xyz, 0.0)));
            vec3 bitangentW = cross(normalW, tangentW) * tangent.w;
            fsout.TBN = mat3(tangentW, bitangentW, normalW);
        #else
            // if HAS_NORMAL
            fsout.normal = normalize(vec3(ModelMatrix * vec4(normal.xyz, 0.0)));
        #endif
    #endif
    // TexCoord
    #ifdef HAS_COORD_0
        fsout.texCoord_0 = texCoord_0;
    #else
        fsout.texCoord_0 = vec2(0.0,0.0);
    #endif
    // Camera = Perspective * View
    gl_Position = CameraMatrix * ModelMatrix * vec4(position, 1); // needs w for proper perspective correction
}