package preset

import "errors"

var (
    ErrPresetNotFound = errors.New("preset not found")
    ErrPresetAlreadyExists = errors.New("preset already exists")
    ErrEmptyName = errors.New("preset name must not be empty")
    ErrNameTooLong = errors.New("preset name is too long")
    ErrDescriptionTooLong = errors.New("preset description is too long")
    ErrNoItems = errors.New("preset must contain at least one item")
    ErrTotalPriceMismatch = errors.New("total price does not match sum of items")
)

var (
    ErrItemNotFound = errors.New("preset item not found")
    ErrInvalidProductID = errors.New("invalid product ID in preset item")
    ErrNilProductSummary = errors.New("product summary must not be nil")
)
