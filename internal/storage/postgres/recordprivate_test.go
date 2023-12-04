package postgres

func (s *postgresSuite) TestUniqueID() {
	strA, strB := "foo", "bar"
	id1 := uniqueID(strA, strB)

	strA, strB = strB, strA
	id2 := uniqueID(strA, strB)

	s.Require().Equal(id1, id2)
}
