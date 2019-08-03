package tests

import (
	"testing"

	"github.com/gernest/hiro/models"
	uuid "github.com/satori/go.uuid"
)

func TestQRCodeCacheRequest(t *testing.T) {
	id, err := uuid.FromString("b51ab8d0-b122-4d0e-b235-97448acc0f7f")
	if err != nil {
		t.Fatal(err)
	}
	o := &models.QRCodeCacheRequest{
		Owner:  id,
		UUID:   id,
		Width:  150,
		Height: 150,
	}
	expect := `usr%3Db51ab8d0-b122-4d0e-b235-97448acc0f7f%26id%3Db51ab8d0-b122-4d0e-b235-97448acc0f7f%26w%3D150%26h%3D150`
	key := o.Key()
	if key != expect {
		t.Errorf("expected %s got %s", expect, key)
	}

	n := &models.QRCodeCacheRequest{}
	err = n.FromKey(key)
	if err != nil {
		t.Fatal(err)
	}

	if !uuid.Equal(n.UUID, o.UUID) {
		t.Errorf("expected %s got %s", o.UUID, n.UUID)
	}
	if !uuid.Equal(n.Owner, o.Owner) {
		t.Errorf("expected %s got %s", o.Owner, n.Owner)
	}
	if n.Width != o.Width {
		t.Errorf("expected %d got %d", o.Width, n.Width)
	}
	if n.Height != o.Height {
		t.Errorf("expected %d got %d", o.Height, n.Height)
	}
}
