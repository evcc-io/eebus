package server

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"log"

	"github.com/denisbrodbeck/machineid"
)

type UniqueID struct {
	id     string
	Prefix string
}

func (u UniqueID) String() string {
	if len(u.id) == 0 {
		u.id = u.generate()
	}
	return u.id
}

func (u UniqueID) generate() string {
	id, err := uniqueID()
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%s-%0x", u.Prefix, id)
}

func uniqueID() ([]byte, error) {
	id, err := machineid.ID()
	if err != nil {
		return nil, err
	}

	mac := hmac.New(sha1.New, []byte(id))
	_, err = mac.Write([]byte(id))
	sum := mac.Sum(nil)

	return sum[:8], err
}
