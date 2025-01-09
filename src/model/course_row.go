package model

type CourseRow struct {
	CourseCode       string
	CourseName       string
	EvaluationAmount int
	EmptyNoteAmount  int
	CourseLink       string
}

func (cr *CourseRow) Equal(other *CourseRow) bool {
	return cr.CourseCode == other.CourseCode &&
		cr.CourseName == other.CourseName &&
		cr.EvaluationAmount == other.EvaluationAmount &&
		cr.EmptyNoteAmount == other.EmptyNoteAmount
}
