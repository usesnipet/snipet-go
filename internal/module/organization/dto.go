package organization

type CreateOrganizationDTO struct {
	Name string `json:"name" validate:"required,max=255"`
	Slug string `json:"slug" validate:"required,max=255"`
}
