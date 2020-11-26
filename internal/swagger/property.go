package swagger

// Property definats a property item in Definition object
type Property struct {
	BaseObject
	Item *BaseObject `json:"items,omitempty"`
}
