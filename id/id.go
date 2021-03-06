package id

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/sony/sonyflake"
	"github.com/sony/sonyflake/awsutil"
	"github.com/speps/go-hashids"
)

const (
	defaultSalt          = "z@mmik_orvyl"
	defaultSonyflakeTime = "2017-01-02T08:30:00"
)

//Generator wiil produce either string or numeric ID
type Generator interface {
	Next() (interface{}, error)
}

//Settings that will support id generator seed
type Settings struct {
	TimeSeed   time.Time
	Salt       string
	UseAWSData bool
}

type alphanum struct {
	sf *sonyflake.Sonyflake
	h  *hashids.HashID
}

func (a *alphanum) Next() (interface{}, error) {
	nid, err := a.sf.NextID()
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to execute sonyflake.NextID()")
	}

	id, err := a.h.Encode([]int{int(nid)})
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to encode ID using hashids.HashID")
	}

	return id, nil
}

type num struct {
	sf *sonyflake.Sonyflake
}

func (n *num) Next() (interface{}, error) {
	id, err := n.sf.NextID()
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to execute sonyflake.NextID()")
	}

	return id, nil
}

//NewGenerator produces generator. Either string or numeric ID generator
func NewGenerator(isAlphaNumeric bool, s Settings) (Generator, error) {
	sf, err := newSonyflake(s.TimeSeed, s.UseAWSData)
	if err != nil {
		return nil, err
	}

	if isAlphaNumeric {
		hd := hashids.NewData()

		hd.Salt = defaultSalt
		if len(s.Salt) > 0 {
			hd.Salt = s.Salt
		}

		h, err := hashids.NewWithData(hd)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to initialized hashids.HashID")
		}

		return &alphanum{h: h, sf: sf}, nil
	}

	return &num{sf: sf}, nil
}

func newSonyflake(tseed time.Time, UseAWSData bool) (*sonyflake.Sonyflake, error) {
	var s sonyflake.Settings

	s.StartTime = tseed
	if tseed.IsZero() {
		s.StartTime, _ = time.Parse("2006-01-02T15:04:05", defaultSonyflakeTime)
	}

	if UseAWSData {
		s.MachineID = awsutil.AmazonEC2MachineID
	}

	sf := sonyflake.NewSonyflake(s)
	if sf == nil {
		return nil, fmt.Errorf("Failed to initialize ID generator (sonyflake)")
	}

	return sf, nil
}
