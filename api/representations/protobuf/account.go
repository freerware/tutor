package protobuf

import (
	"errors"

	"github.com/freerware/negotiator/representation"
	r "github.com/freerware/tutor/api/representations"
	"github.com/freerware/tutor/api/representations/protobuf/gen"
	"github.com/freerware/tutor/domain"
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	mediaTypeProtobuf  = "application/protobuf"
	mediaTypeXProtobuf = "application/x-protobuf"
)

type Account struct {
	gen.Account
	r.Representation
}

// NewAccount constructs a new account representation.
func NewAccount(a domain.Account) Account {
	marshaller := func(in interface{}) ([]byte, error) {
		message, ok := in.(proto.Message)
		if !ok {
			return []byte{}, errors.New("must provide Protobuf message to marshal successfully")
		}
		return proto.Marshal(message)
	}
	unmarshaller := func(b []byte, out interface{}) error {
		message, ok := out.(proto.Message)
		if !ok {
			return errors.New("must provide Protobuf message to unmarshal successfully")
		}
		return proto.Unmarshal(b, message)
	}
	acc := Account{}
	acc.UUID = a.UUID().String()
	acc.GivenName = a.GivenName()
	acc.Surname = a.Surname()
	acc.Username = a.Username()
	c := timestamppb.Timestamp{Seconds: a.CreatedAt().Unix()}
	acc.CreatedAt = &c
	u := timestamppb.Timestamp{Seconds: a.UpdatedAt().Unix()}
	acc.UpdatedAt = &u
	var d *timestamppb.Timestamp
	if a.DeletedAt() != nil {
		d = &timestamppb.Timestamp{Seconds: a.DeletedAt().Unix()}
	}
	acc.DeletedAt = d
	acc.SetContentCharset("ascii")
	acc.SetContentLanguage("en-US")
	acc.SetContentType(mediaTypeProtobuf)
	acc.SetSourceQuality(1.0)
	acc.SetContentEncoding([]string{"identity"})
	acc.SetMarshallers(map[string]representation.Marshaller{
		mediaTypeProtobuf:  marshaller,
		mediaTypeXProtobuf: marshaller,
	})
	acc.SetUnmarshallers(map[string]representation.Unmarshaller{
		mediaTypeProtobuf:  unmarshaller,
		mediaTypeXProtobuf: unmarshaller,
	})
	return acc
}

func (a Account) Bytes() ([]byte, error) {
	return a.Base.Bytes(&a)
}

func (a Account) FromBytes(b []byte) error {
	return a.Base.FromBytes(b, &a)
}
