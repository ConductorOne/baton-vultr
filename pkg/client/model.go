package client

type User struct {
	Id         string   `json:"id"`
	Name       string   `json:"name"`
	Email      string   `json:"email"`
	ApiEnabled bool     `json:"apiEnabled"`
	ACLs       []string `json:"acls"`
}

type UserResponse struct {
	Meta   *PaginationMeta `json:"meta"`
	Result []User          `json:"users"`
}

// HasPaginationData satisfies uhttp.PaginatedResponse. Vultr sends meta on every list
// response and empties meta.links.next on the last page, so meta going missing means the
// cursor was dropped - which would both end the sync early and nil-deref meta.links.next.
func (r *UserResponse) HasPaginationData() bool {
	return r.Meta != nil
}

type PaginationMeta struct {
	Total int `json:"total"`
	Links struct {
		Next string `json:"next"`
		Prev string `json:"prev"`
	} `json:"links"`
}

type UserSingleResponse struct {
	Result User `json:"user"`
}

type Account struct {
	Balance           float64  `json:"balance"`
	PendingCharges    float64  `json:"pending_charges"`
	Name              string   `json:"name"`
	Email             string   `json:"email"`
	ACLs              []string `json:"acls"`
	LastPaymentDate   string   `json:"last_payment_date"`
	LastPaymentAmount float64  `json:"last_payment_amount"`
}

type AccountResponse struct {
	Account Account `json:"account"`
}

type ACL string
