package preset

import (
	"strings"
	"time"

	"github.com/Neimess/zorkin-store-project/internal/domain/product"
)

type Preset struct {
	ID          int64
	Name        string
	Description *string
	TotalPrice  float64
	ImageURL    *string
	CreatedAt   time.Time
	Items       []PresetItem
}

type PresetItem struct {
	ID        int64
	PresetID  int64
	ProductID int64
	Product   *product.ProductSummary
}

func (p *Preset) Validate() error {
	name := strings.TrimSpace(p.Name)
	if name == "" {
		return ErrEmptyName
	}
	if len(name) > 100 {
		return ErrNameTooLong
	}
	if len(*p.Description) > 500 {
		return ErrDescriptionTooLong
	}
	if len(p.Items) == 0 {
		return ErrNoItems
	}
	return nil
}
