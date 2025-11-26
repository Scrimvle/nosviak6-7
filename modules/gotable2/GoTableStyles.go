package gotable2

type Style struct {
	BorderValues 	int 	`json:"border_values"`
	*Header 			`json:"header"`
	*Body 			`json:"body"`
}

// Header is the styling for a table header
type Header struct {
	AboveHeaderLeft 		*string
	AboveHeaderRight 		*string
	AboveHeaderIntersection 	*string
	AboveHeaderHorizontal 	        *string
	HeaderLeft 		        string
	HeaderIntersection 		string
	HeaderRight 		        string
	BelowHeaderLeft 		*string
	BelowHeaderRight 		*string
	BelowHeaderIntersection 	*string
	BelowHeaderHorizontal 	        *string
}

type Body struct {
	ValueLeft 		string
	ValueIntersection 	string
	ValueRight 		string
	BelowBodyLeft 		*string
	BelowBodyRight 		*string
	BelowBodyIntersection 	*string
	BelowBodyHorizontal 	*string

}

var (
	// DEFAULTBOLD is the default style for the header
	DEFAULTBOLD *Style = &Style{
		BorderValues: 0,
		Header: &Header{
			AboveHeaderLeft: stringPointer("╔"),
			AboveHeaderRight: stringPointer("╗"),
			AboveHeaderIntersection: stringPointer("╦"),
			AboveHeaderHorizontal: stringPointer("═"),

			HeaderLeft: "║",
			HeaderRight: "║",
			HeaderIntersection: "║",

			BelowHeaderLeft: stringPointer("╠"),
			BelowHeaderRight: stringPointer("╣"),
			BelowHeaderIntersection: stringPointer("╬"),
			BelowHeaderHorizontal: stringPointer("═"),
		},

		Body: &Body{
			ValueLeft: "║",
			ValueRight: "║",
			ValueIntersection: "║",

			BelowBodyLeft: stringPointer("╚"),
			BelowBodyRight: stringPointer("╝"),
			BelowBodyHorizontal: stringPointer("═"),
			BelowBodyIntersection: stringPointer("╩"),
		},
	}

	// DEFAULT is the default style for the header
	DEFAULT *Style = &Style{
		BorderValues: 0,
		Header: &Header{
			HeaderLeft: "│",
			HeaderRight: "│",
			HeaderIntersection: "│",

			BelowHeaderLeft: stringPointer("├"),
			BelowHeaderRight: stringPointer("┤"),
			BelowHeaderIntersection: stringPointer("┼"),
			BelowHeaderHorizontal: stringPointer("─"),
		},

		Body: &Body{
			ValueLeft: "│",
			ValueRight: "│",
			ValueIntersection: "│",
		},
	}
	
	SLICK *Style = &Style{
		BorderValues: 0,
		Header: &Header{
			HeaderLeft: "",
			HeaderRight: "",
			HeaderIntersection: "",
		},

		Body: &Body{
			ValueLeft: "",
			ValueRight: "",
			ValueIntersection: "",
		},
	}
)

func (h *Header) hasAboveHeader() bool {
	return h != nil && h.AboveHeaderLeft != nil && h.AboveHeaderRight != nil && h.AboveHeaderIntersection != nil && h.AboveHeaderHorizontal != nil
}

func (h *Header) hasBelowHeader() bool {
	return h != nil && h.BelowHeaderLeft != nil && h.BelowHeaderRight != nil && h.BelowHeaderIntersection != nil && h.BelowHeaderHorizontal != nil
}

func (b *Body) hasBelowBody() bool {
	return b != nil && b.BelowBodyLeft != nil && b.BelowBodyRight != nil && b.BelowBodyIntersection != nil && b.BelowBodyHorizontal != nil
}