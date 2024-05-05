package utils

import (
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

// GenerateQRCode 生成图形验证码
func GenerateQRCode(text string) (barcode.Barcode, error) {
	qrCode, err := qr.Encode(text, qr.M, qr.Auto)
	if err != nil {
		return nil, err
	}
	qrCode, err = barcode.Scale(qrCode, 200, 200)
	if err != nil {
		return nil, err
	}
	return qrCode, nil
}
