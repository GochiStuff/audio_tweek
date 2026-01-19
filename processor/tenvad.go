package processor

type TenVADHandler struct {
	// TODO: plug a real detector here
	isActive bool
}

func NewTenVADHandler() (*TenVADHandler, error) {
	return &TenVADHandler{}, nil
}

func (h *TenVADHandler) IsSpeech(samples []int16) bool {
	h.isActive = false
	return h.isActive
}