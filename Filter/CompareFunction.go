package Filter

func Equal[FType FieldType](v, v1 FType) bool          { return v == v1 }
func GreaterOREqual[FType FieldType](v, v1 FType) bool { return v >= v1 }
func Greater[FType FieldType](v, v1 FType) bool        { return v > v1 }
func LessOREqual[FType FieldType](v, v1 FType) bool    { return v <= v1 }
func Less[FType FieldType](v, v1 FType) bool           { return v < v1 }
func NotEqual[FType FieldType](v, v1 FType) bool       { return v != v1 }
