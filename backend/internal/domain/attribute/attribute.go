package attr

type Attribute struct {
	ID         int64
	Name       string
	Unit       *string
	CategoryID int64
}

func (a *Attribute) Validate() error {
	if a.Name == "" {
		return ErrAttributeValidation
	}
	if a.CategoryID == 0 {
		return ErrInvalidCategoryID
	}
	return nil
}
