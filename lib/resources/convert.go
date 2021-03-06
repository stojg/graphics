package resources

import (
	"math"
)

// ConvertToVertices takes an slice of float32 and turns them into nice Vertexes.
// It requires that the in data is packed in this order: [3] position, [3] normals, [2] texture coordinates. If the
// in data doesn't follow this convention, there will be tears and possibly your GPU will implode.
func ConvertToVertices(meshdata []float32, indices []uint32) []Vertex {
	const stride = 8

	if len(meshdata)%stride != 0 {
		panic("the mesh data is not a multiple of 8, want [3]Pos, [3]Normals, [2]TexCoords")
	}
	vertices := make([]Vertex, len(meshdata)/stride)

	// 1. Add Pos, Normal and TexCoords to all Vertices
	for i := 0; i < len(meshdata); i += stride {
		copy(vertices[i/stride].Pos[:], meshdata[i:i+3])
		copy(vertices[i/stride].Normal[:], meshdata[i+3:i+6])
		copy(vertices[i/stride].TexCoords[:], meshdata[i+6:i+8])
	}

	// 2. calculate tangents from the texture UVs so we can properly use bumpmaps texture on meshes (we can calculate the bi-tangents
	// in the vertex shader when we need it)

	var edge1 [3]float32
	var edge2 [3]float32
	var tangent [3]float32
	for indexPos := 0; indexPos < len(indices); indexPos += 3 {
		// check if we already have calculated the tangents
		if sqrLength(vertices[indices[indexPos]].Tangent) != 0 {
			continue
		}

		v0 := vertices[indices[indexPos]]
		v1 := vertices[indices[indexPos+1]]
		v2 := vertices[indices[indexPos+2]]

		deltaU1 := v1.TexCoords[0] - v0.TexCoords[0]
		deltaV1 := v1.TexCoords[1] - v0.TexCoords[1]
		deltaU2 := v2.TexCoords[0] - v0.TexCoords[0]
		deltaV2 := v2.TexCoords[1] - v0.TexCoords[1]

		f := 1.0 / (deltaU1*deltaV2 - deltaU2*deltaV1)

		edge1 = [3]float32{v1.Pos[0] - v0.Pos[0], v1.Pos[1] - v0.Pos[1], v1.Pos[2] - v0.Pos[2]}
		edge2 = [3]float32{v2.Pos[0] - v0.Pos[0], v2.Pos[1] - v0.Pos[1], v2.Pos[2] - v0.Pos[2]}

		tangent = [3]float32{
			f * (deltaV2*edge1[0] - deltaV1*edge2[0]),
			f * (deltaV2*edge1[1] - deltaV1*edge2[1]),
			f * (deltaV2*edge1[2] - deltaV1*edge2[2]),
		}

		// normalise tangent
		l := 1.0 / float32(math.Sqrt(float64(sqrLength(tangent))))
		tangent[0] *= l
		tangent[1] *= l
		tangent[2] *= l

		// check handiness with normal.Dot(tangent)
		if vertices[indices[indexPos]].Normal[0]*tangent[0]+vertices[indices[indexPos]].Normal[1]*tangent[1]+vertices[indices[indexPos]].Normal[2]*tangent[2] < 0.0 {
			tangent[0] *= -1
			tangent[1] *= -1
			tangent[2] *= -1
		}

		copy(vertices[indices[indexPos]].Tangent[:], tangent[:])
		copy(vertices[indices[indexPos+1]].Tangent[:], tangent[:])
		copy(vertices[indices[indexPos+2]].Tangent[:], tangent[:])
	}
	return vertices
}

func edge(a, b Vertex) [3]float32 {
	return [3]float32{a.Pos[0] - b.Pos[0], a.Pos[1] - b.Pos[1], a.Pos[2] - b.Pos[2]}
}

func normalise(vec [3]float32) [3]float32 {
	l := 1.0 / float32(math.Sqrt(float64(sqrLength(vec))))
	return [3]float32{vec[0] * l, vec[1] * l, vec[2] * l}
}

func sqrLength(vec [3]float32) float32 {
	return vec[0]*vec[0] + vec[1]*vec[1] + vec[2]*vec[2]
}

//func isZero(vec [3]float32) bool {
//	return vec[0] == 0 && vec[1] == 0 && vec[2] == 0
//}
