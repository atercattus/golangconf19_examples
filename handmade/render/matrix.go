package render

// Copy&paste from github.com/EngoEngine/engo/math.go

type (
	Matrix struct {
		Val [9]float32
		tmp [9]float32
	}
)

const (
	m00 = 0
	m01 = 3
	m02 = 6
	m10 = 1
	m11 = 4
	m12 = 7
	m20 = 2
	m21 = 5
	m22 = 8
)

func (m *Matrix) Identity() *Matrix {
	m.Val[m00] = 1
	m.Val[m10] = 0
	m.Val[m20] = 0
	m.Val[m01] = 0
	m.Val[m11] = 1
	m.Val[m21] = 0
	m.Val[m02] = 0
	m.Val[m12] = 0
	m.Val[m22] = 1
	return m
}

func (m *Matrix) Translate(x, y float32) *Matrix {
	m.tmp[m00] = 1
	m.tmp[m10] = 0
	m.tmp[m20] = 0

	m.tmp[m01] = 0
	m.tmp[m11] = 1
	m.tmp[m21] = 0

	m.tmp[m02] = x
	m.tmp[m12] = y
	m.tmp[m22] = 1

	multiplyMatrixes(m.Val[:], m.tmp[:])
	return m
}

func (m *Matrix) Scale(x, y float32) *Matrix {
	m.tmp[m00] = x
	m.tmp[m10] = 0
	m.tmp[m20] = 0

	m.tmp[m01] = 0
	m.tmp[m11] = y
	m.tmp[m21] = 0

	m.tmp[m02] = 0
	m.tmp[m12] = 0
	m.tmp[m22] = 1

	multiplyMatrixes(m.Val[:], m.tmp[:])
	return m
}

func multiplyMatrixes(m1, m2 []float32) {
	v00 := m1[m00]*m2[m00] + m1[m01]*m2[m10] + m1[m02]*m2[m20]
	v01 := m1[m00]*m2[m01] + m1[m01]*m2[m11] + m1[m02]*m2[m21]
	v02 := m1[m00]*m2[m02] + m1[m01]*m2[m12] + m1[m02]*m2[m22]

	v10 := m1[m10]*m2[m00] + m1[m11]*m2[m10] + m1[m12]*m2[m20]
	v11 := m1[m10]*m2[m01] + m1[m11]*m2[m11] + m1[m12]*m2[m21]
	v12 := m1[m10]*m2[m02] + m1[m11]*m2[m12] + m1[m12]*m2[m22]

	v20 := m1[m20]*m2[m00] + m1[m21]*m2[m10] + m1[m22]*m2[m20]
	v21 := m1[m20]*m2[m01] + m1[m21]*m2[m11] + m1[m22]*m2[m21]
	v22 := m1[m20]*m2[m02] + m1[m21]*m2[m12] + m1[m22]*m2[m22]
	m1[m00] = v00
	m1[m10] = v10
	m1[m20] = v20
	m1[m01] = v01
	m1[m11] = v11
	m1[m21] = v21
	m1[m02] = v02
	m1[m12] = v12
	m1[m22] = v22
}
