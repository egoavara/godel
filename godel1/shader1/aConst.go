package shader1

const (
	VertexAttributePosition           = 0
	VertexAttributeNormal             = 1
	VertexAttributeTangent            = 2
	VertexAttributeColor              = 3
	VertexAttributeTexCoord0          = 4
	VertexAttributeTexCoord1          = 5
	VertexAttributeJoint0             = 8
	VertexAttributeWeight0            = 9
)

var (
	VertexAttributeMorphPosition = [8]int{10, 13, 16, 19, 22, 25, 28, 31}
	VertexAttributeMorphNormal   = [8]int{11, 14, 17, 20, 23, 26, 29, 32}
	VertexAttributeMorphTangent  = [8]int{12, 15, 18, 21, 24, 27, 30, 33}
)
