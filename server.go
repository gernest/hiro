package main

import (
	"context"

	"gitlab.com/gernest/hiro/hiro"
)

type Server struct{}

// Create creates a new qrcode.
func (s Server) Create(ctx context.Context, r *hiro.QRCodeRequest) (*hiro.QRCodeResponse, error) {

}
