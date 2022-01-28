package rule

type CreateMemberRequest struct {
	Nickname string   `json:"nickname,omitempty"`  // required，不小于 4 位 不超过 20 位
	Username string   `json:"username,omitempty"`  // required，只支持大小写，长度不小于 8 位 不超过 20 位
	Password string   `json:"password,omitempty"`  // required，同时包括大小写、数字，长度不少于 8 位 不超过 20 位
	UserType UserType `json:"user_type,omitempty"` // required, 枚举值
}

type CreateMemberResponse struct {
	Code ErrNo
	Data struct {
		UserID string // int64 范围
	}
}
