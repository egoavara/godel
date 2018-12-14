precision highp float;


const float PI = 3.141592653589793;
const float MINROUGHNESS = 0.04;


// Light
uniform vec3 LightDir;
uniform vec3 LightColor;

// Material

uniform vec4 BaseColorFactor;
uniform vec2 MetalRoughnessFactor;
uniform vec3 EmissiveFactor;

#ifdef HAS_BASECOLORTEX
uniform sampler2D BaseColorTex;
uniform int BaseColorTexCoord;
#endif

#ifdef HAS_METALROUGHNESSTEX
uniform sampler2D MetalRoughnessTex;
uniform int MetalRoughnessTexCoord;
#endif

#ifdef HAS_NORMALTEX
uniform sampler2D NormalTex;
uniform float NormalScale;
uniform int NormalTexCoord;
#endif

#ifdef HAS_EMISSIVETEX
uniform sampler2D EmissiveTex;
uniform int EmissiveTexCoord;
#endif

#ifdef HAS_OCCLUSIONTEX
uniform sampler2D OcculusionTex;
uniform float OcclusionStrength;
uniform int OcculusionTexCoord;
#endif

uniform vec3 Camera;
//
in struct{
   vec3 position;
    #ifdef HAS_COORD_0
    vec2 texCoord_0;
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



out vec4 outputColor;

//
struct PBRInfo
{
    float NdotL;                  // cos angle between normal and light direction
    float NdotV;                  // cos angle between normal and view direction
    float NdotH;                  // cos angle between normal and half vector
    float LdotH;                  // cos angle between light direction and half vector
    float VdotH;                  // cos angle between view direction and half vector
    float perceptualRoughness;    // roughness value, as authored by the model creator (input to shader)
    float metalness;              // metallic value at the surface
    vec3 reflectance0;            // full reflectance color (normal incidence angle)
    vec3 reflectance90;           // reflectance color at grazing angle
    float alphaRoughness;         // roughness mapped to a more linear change in the roughness (proposed by [2])
    vec3 diffuseColor;            // color contribution from diffuse lighting
    vec3 specularColor;           // color contribution from specular lighting
};
vec2 getTexCoord(int i){
    #ifdef HAS_COORD_0
        if(i == 0){
            return fsout.texCoord_0;
        }
    #endif
    #ifdef HAS_COORD_1
        if(i == 1){
            return fsout.texCoord_1;
        }
    #endif
    return vec2(0.f,0.f);
}
vec3 getNormal(){
    // Retrieve the tangent space matrix
    #ifndef HAS_TANGENT
        vec3 pos_dx = dFdx(fsout.position);
        vec3 pos_dy = dFdy(fsout.position);
        vec3 tex_dx = dFdx(vec3(getTexCoord(NormalTexCoord), 0.0));
        vec3 tex_dy = dFdy(vec3(getTexCoord(NormalTexCoord), 0.0));
        vec3 t = (tex_dy.t * pos_dx - tex_dx.t * pos_dy) / (tex_dx.s * tex_dy.t - tex_dy.s * tex_dx.t);

        #ifdef HAS_NORMAL
            vec3 ng = normalize(fsout.normal);
        #else
            vec3 ng = cross(pos_dx, pos_dy);
        #endif

        t = normalize(t - ng * dot(ng, t));
        vec3 b = normalize(cross(ng, t));
        mat3 tbn = mat3(t, b, ng);
    #else // HAS_TANGENTS
        mat3 tbn = fsout.TBN;
    #endif

    #ifdef HAS_NORMALTEX
        vec3 n = texture(NormalTex, getTexCoord(NormalTexCoord)).rgb;
        n = normalize(tbn * ((2.0 * n - 1.0) * vec3(NormalScale, NormalScale, 1.0)));
    #else
        // The tbn matrix is linearly interpolated, so we need to re-normalize
        vec3 n = normalize(tbn[2].xyz);
    #endif

    return n;
}
vec3 diffuse(PBRInfo pbrInputs){
    return pbrInputs.diffuseColor / PI;
}
vec3 specularReflection(PBRInfo pbrInputs){
    return pbrInputs.reflectance0 + (pbrInputs.reflectance90 - pbrInputs.reflectance0) * pow(clamp(1.0 - pbrInputs.VdotH, 0.0, 1.0), 5.0);
}
float geometricOcclusion(PBRInfo pbrInputs){
    float NdotL = pbrInputs.NdotL;
    float NdotV = pbrInputs.NdotV;
    float r = pbrInputs.alphaRoughness;

    float attenuationL = 2.0 * NdotL / (NdotL + sqrt(r * r + (1.0 - r * r) * (NdotL * NdotL)));
    float attenuationV = 2.0 * NdotV / (NdotV + sqrt(r * r + (1.0 - r * r) * (NdotV * NdotV)));
    return attenuationL * attenuationV;
}
float microfacetDistribution(PBRInfo pbrInputs){
    float roughnessSq = pbrInputs.alphaRoughness * pbrInputs.alphaRoughness;
    float f = (pbrInputs.NdotH * roughnessSq - pbrInputs.NdotH) * pbrInputs.NdotH + 1.0;
    return roughnessSq / (PI * f * f);
}

void main()
{
    float metallic = MetalRoughnessFactor.x;
    float perceptualRoughness = MetalRoughnessFactor.y;
    #ifdef HAS_METALROUGHNESSTEX
        vec4 mrSample = texture(MetalRoughnessTex, getTexCoord(MetalRoughnessTexCoord));
        perceptualRoughness = mrSample.g * perceptualRoughness;
        metallic = mrSample.b * metallic;
    #endif
    perceptualRoughness = clamp(perceptualRoughness, MINROUGHNESS, 1.0);
    metallic = clamp(metallic, 0.0, 1.0);
    float alphaRoughness = perceptualRoughness * perceptualRoughness;
    #ifdef HAS_BASECOLORTEX
        vec4 baseColor = texture(BaseColorTex, getTexCoord(BaseColorTexCoord)) * BaseColorFactor;
    #else
        vec4 baseColor = BaseColorFactor;
    #endif
    vec3 f0 = vec3(0.04);
    vec3 diffuseColor = baseColor.rgb * (vec3(1.0) - f0);
    diffuseColor *= 1.0 - metallic;
    vec3 specularColor = mix(f0, baseColor.rgb, metallic);
    float reflectance = max(max(specularColor.r, specularColor.g), specularColor.b);
    float reflectance90 = clamp(reflectance * 25.0, 0.0, 1.0);
    vec3 specularEnvironmentR0 = specularColor.rgb;
    vec3 specularEnvironmentR90 = vec3(1.0, 1.0, 1.0) * reflectance90;

    vec3 n = getNormal();
    vec3 v = normalize(Camera - fsout.position);
    vec3 l = normalize(LightDir);
    vec3 h = normalize(l+v);
    vec3 reflection = -normalize(reflect(v, n));

    float NdotL = clamp(dot(n, l), 0.001, 1.0);
    float NdotV = clamp(abs(dot(n, v)), 0.001, 1.0);
    float NdotH = clamp(dot(n, h), 0.0, 1.0);
    float LdotH = clamp(dot(l, h), 0.0, 1.0);
    float VdotH = clamp(dot(v, h), 0.0, 1.0);

    PBRInfo pbrInputs = PBRInfo(
        NdotL,
        NdotV,
        NdotH,
        LdotH,
        VdotH,
        perceptualRoughness,
        metallic,
        specularEnvironmentR0,
        specularEnvironmentR90,
        alphaRoughness,
        diffuseColor,
        specularColor
    );
    // ================================================================================================
    // PBRInfo pbrInputs : complete
    // ================================================================================================
    // Calculate the shading terms for the microfacet specular shading model
    vec3 F = specularReflection(pbrInputs);
    float G = geometricOcclusion(pbrInputs);
    float D = microfacetDistribution(pbrInputs);

    // Calculation of analytical lighting contribution
    vec3 diffuseContrib = (1.0 - F) * diffuse(pbrInputs);
    vec3 specContrib = F * G * D / (4.0 * NdotL * NdotV);
    // Obtain final intensity as reflectance (BRDF) scaled by the energy of the light (cosine law)
    vec3 color = NdotL * LightColor * (diffuseContrib + specContrib);

    // Apply optional PBR terms for additional (optional) shading
    #ifdef HAS_OCCLUSIONTEX
        float ao = texture(OcculusionTex, getTexCoord(OcculusionTexCoord)).r;
        color = mix(color, color * ao, OcclusionStrength);
    #endif

    #ifdef HAS_EMISSIVETEX
        vec3 emissive = texture(EmissiveTex, getTexCoord(EmissiveTexCoord)).rgb * EmissiveFactor;
        color += emissive;
    #endif
    outputColor = vec4(pow(color,vec3(1.0/2.2)), baseColor.a);
}
