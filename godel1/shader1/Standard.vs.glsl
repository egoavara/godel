
//
// Reference > https://github.com/KhronosGroup/glTF-WebGL-PBR/blob/master/shaders/pbr-vert.glsl
//

// https://www.khronos.org/opengl/wiki/Layout_Qualifier_(GLSL)#Explicit_uniform_location
// ... GL_MAX_UNIFORM_LOCATIONS, which will be at least 1024 ...
//
uniform mat4 CameraMatrix;
uniform mat4 ModelMatrix;
uniform mat4 NormalMatrix;
#ifdef HAS_JOINT_0
#ifdef HAS_WEIGHT_0
uniform mat4 [64]JointMatrix;
#endif
#endif
#ifdef MORPH_SIZE
uniform float [MORPH_SIZE]MorphWeight;
#endif

layout (location = 0) in vec3 position;
#ifdef HAS_NORMAL
layout (location = 1) in vec3 normal;
#endif
#ifdef HAS_TANGENT
layout (location = 2) in vec4 tangent;
#endif
#ifdef HAS_COLOR
layout (location = 3) in vec4 color;
#endif
#ifdef HAS_COORD_0
layout (location = 4) in vec2 texCoord_0;
#endif
#ifdef HAS_JOINTS_0
layout (location = 5) in vec2 texCoord_1;
#endif
#ifdef HAS_JOINT_0
layout (location = 8) in ivec4 joint_0;
#endif
#ifdef HAS_WEIGHT_0
layout (location = 9) in vec4 weight_0;
#endif
// reserved, joint_1
// reserved, weight_1

// https://github.com/KhronosGroup/glTF/tree/master/specification/2.0#morph-targets
// gltf don't have limit for morph, but, at least 8 morph support
#ifdef HAS_MORPH_POSITION
    #if MORPH_SIZE >= 1
    layout (location = 10) in vec3 morph_0_position;
    #endif
    #if MORPH_SIZE >= 2
    layout (location = 13) in vec3 morph_1_position;
    #endif
    #if MORPH_SIZE >= 3
    layout (location = 16) in vec3 morph_2_position;
    #endif
    #if MORPH_SIZE >= 4
    layout (location = 19) in vec3 morph_3_position;
    #endif
    #if MORPH_SIZE >= 5
    layout (location = 22) in vec3 morph_4_position;
    #endif
    #if MORPH_SIZE >= 6
    layout (location = 25) in vec3 morph_5_position;
    #endif
    #if MORPH_SIZE >= 7
    layout (location = 28) in vec3 morph_6_position;
    #endif
    #if MORPH_SIZE >= 8
    layout (location = 31) in vec3 morph_7_position;
    #endif
#endif
#ifdef HAS_MORPH_NORMAL
    #if MORPH_SIZE >= 1
    layout (location = 11) in vec3 morph_0_normal;
    #endif
    #if MORPH_SIZE >= 2
    layout (location = 14) in vec3 morph_1_normal;
    #endif
    #if MORPH_SIZE >= 3
    layout (location = 17) in vec3 morph_2_normal;
    #endif
    #if MORPH_SIZE >= 4
    layout (location = 20) in vec3 morph_3_normal;
    #endif
    #if MORPH_SIZE >= 5
    layout (location = 23) in vec3 morph_4_normal;
    #endif
    #if MORPH_SIZE >= 6
    layout (location = 26) in vec3 morph_5_normal;
    #endif
    #if MORPH_SIZE >= 7
    layout (location = 39) in vec3 morph_6_normal;
    #endif
    #if MORPH_SIZE >= 8
    layout (location = 32) in vec3 morph_7_normal;
    #endif
#endif
#ifdef HAS_MORPH_TANGENT
    #if MORPH_SIZE >= 1
    layout (location = 12) in vec4 morph_0_tangent;
    #endif
    #if MORPH_SIZE >= 2
    layout (location = 15) in vec4 morph_1_tangent;
    #endif
    #if MORPH_SIZE >= 3
    layout (location = 18) in vec4 morph_2_tangent;
    #endif
    #if MORPH_SIZE >= 4
    layout (location = 21) in vec4 morph_3_tangent;
    #endif
    #if MORPH_SIZE >= 5
    layout (location = 24) in vec4 morph_4_tangent;
    #endif
    #if MORPH_SIZE >= 6
    layout (location = 27) in vec4 morph_5_tangent;
    #endif
    #if MORPH_SIZE >= 7
    layout (location = 30) in vec4 morph_6_tangent;
    #endif
    #if MORPH_SIZE >= 8
    layout (location = 33) in vec4 morph_7_tangent;
    #endif
#endif


out struct{
    vec3 position;
    vec2 texCoord_0;

    #ifdef HAS_COLOR
    vec4 color;
    #endif

    #ifdef HAS_COORD_1
    vec2 texCoord_1;
    #endif
    #ifdef HAS_NORMAL
        #ifdef HAS_TANGENT
            mat3 TBN;
        #else
            vec3 normal;
        #endif
    #endif
} fsout;


vec3 getPosition(){
    vec3 temp = position;
    #ifdef HAS_MORPH_POSITION
        #if MORPH_SIZE >= 1
        temp += MorphWeight[0] * morph_0_position;
        #endif
        #if MORPH_SIZE >= 2
        temp += MorphWeight[1] * morph_1_position;
        #endif
        #if MORPH_SIZE >= 3
        temp += MorphWeight[2] * morph_2_position;
        #endif
        #if MORPH_SIZE >= 4
        temp += MorphWeight[3] * morph_3_position;
        #endif
        #if MORPH_SIZE >= 5
        temp += MorphWeight[4] * morph_4_position;
        #endif
        #if MORPH_SIZE >= 6
        temp += MorphWeight[5] * morph_5_position;
        #endif
        #if MORPH_SIZE >= 7
        temp += MorphWeight[6] * morph_6_position;
        #endif
        #if MORPH_SIZE >= 8
        temp += MorphWeight[7] * morph_7_position;
        #endif
    #endif
    return temp;
}

#ifdef HAS_NORMAL
vec3 getNormal(){
    vec3 temp = normal;
    #ifdef HAS_MORPH_NORMAL
        #if MORPH_SIZE >= 1
        temp += MorphWeight[0] * morph_0_normal;
        #endif
        #if MORPH_SIZE >= 2
        temp += MorphWeight[1] * morph_1_normal;
        #endif
        #if MORPH_SIZE >= 3
        temp += MorphWeight[2] * morph_2_normal;
        #endif
        #if MORPH_SIZE >= 4
        temp += MorphWeight[3] * morph_3_normal;
        #endif
        #if MORPH_SIZE >= 5
        temp += MorphWeight[4] * morph_4_normal;
        #endif
        #if MORPH_SIZE >= 6
        temp += MorphWeight[5] * morph_5_normal;
        #endif
        #if MORPH_SIZE >= 7
        temp += MorphWeight[6] * morph_6_normal;
        #endif
        #if MORPH_SIZE >= 8
        temp += MorphWeight[7] * morph_7_normal;
        #endif
    #endif
    return temp;
}
#endif


#ifdef HAS_TANGENT
vec4 getTangent(){
    vec4 temp = tangent;
    #ifdef HAS_MORPH_TANGENT
        #if MORPH_SIZE >= 1
        temp += MorphWeight[0] * morph_0_tangent;
        #endif
        #if MORPH_SIZE >= 2
        temp += MorphWeight[1] * morph_1_tangent;
        #endif
        #if MORPH_SIZE >= 3
        temp += MorphWeight[2] * morph_2_tangent;
        #endif
        #if MORPH_SIZE >= 4
        temp += MorphWeight[3] * morph_3_tangent;
        #endif
        #if MORPH_SIZE >= 5
        temp += MorphWeight[4] * morph_4_tangent;
        #endif
        #if MORPH_SIZE >= 6
        temp += MorphWeight[5] * morph_5_tangent;
        #endif
        #if MORPH_SIZE >= 7
        temp += MorphWeight[6] * morph_6_tangent;
        #endif
        #if MORPH_SIZE >= 8
        temp += MorphWeight[7] * morph_7_tangent;
        #endif
    #endif
    return temp;
}
#endif

void main() {
    vec4 pos = vec4(getPosition(), 1);
    mat4 mtx_model = ModelMatrix;
    mat4 mtx_normal = NormalMatrix;
    // Skining
    #ifdef HAS_JOINT_0
    #ifdef HAS_WEIGHT_0
        mat4 skinMatrix =
            weight_0.x * JointMatrix[int(joint_0.x)] +
            weight_0.y * JointMatrix[int(joint_0.y)] +
            weight_0.z * JointMatrix[int(joint_0.z)] +
            weight_0.w * JointMatrix[int(joint_0.w)];
        //pos = skinMatrix * pos;
        //
        mtx_model = mtx_model * skinMatrix;
        mtx_normal = transpose(inverse(mtx_model));
    #endif
    #endif
    pos = mtx_model * pos;
    //
    fsout.position = vec3(pos.xyz) / pos.w;
    // Normal,
    #ifdef HAS_NORMAL
        #ifdef HAS_TANGENT
            // if HAS_NORMAL and HAS_TANGENT
            vec3 n = getNormal();
            vec4 t = getTangent();
            vec3 normalW = normalize(vec3(mtx_normal * vec4(n.xyz, 0.0)));
            vec3 tangentW = normalize(vec3(mtx_model * vec4(t.xyz, 0.0)));
            vec3 bitangentW = cross(normalW, tangentW) * t.w;
            fsout.TBN = mat3(tangentW, bitangentW, normalW);
        #else
            // if HAS_NORMAL
            vec3 n = getNormal();
            fsout.normal = normalize(vec3(mtx_normal   * vec4(n.xyz, 0.0)));
        #endif
    #endif


    // TexCoord 1
    #ifdef HAS_COORD_0
        fsout.texCoord_0 = texCoord_0;
    #else
        fsout.texCoord_0 = vec2(0, 0);
//        fsout.texCoord_0 = vec2(float(joint_0.x), float(joint_0.y));
    #endif

    // TexCoord 1
    #ifdef HAS_COORD_1
        fsout.texCoord_1 = texCoord_1;
    #endif

    // Camera = Perspective * View
    gl_Position = CameraMatrix * pos; // needs w for proper perspective correction
}