package id

import (
	"errors"
	"github.com/sony/sonyflake"
	"github.com/sony/sonyflake/awsutil"
	"github.com/speps/go-hashids"
	"log"
	"time"
)

const (
	DEFAULT_SALT           = "z@mmik_orvyl"
	DEFAULT_SONYFLAKE_TIME = "2017-01-02T08:30:00"
)

type Generator interface {
	Next() (interface{}, error)
}

type Settings struct {
	TimeSeed       time.Time
	Salt           string
	IsAlphaNumeric bool
	InDocker       bool
}

type alphanum struct {
	h *hashids.HashID
}

func (a *alphanum) Next() (interface{}, error) {
	var s sonyflake.Settings
	t, _ := time.Parse("2006-01-02T15:04:05", DEFAULT_SONYFLAKE_TIME)
	s.StartTime = t
	idGen := sonyflake.NewSonyflake(s)

	nid, err := idGen.NextID()
	if err != nil {
		log.Panic("Failed to generate ID: " + err.Error())
		return nil, err
	}

	id, err := a.h.Encode([]int{int(nid)})
	if err != nil {
		log.Panic("Failed to generate ID: " + err.Error())
		return nil, err
	}

	return id, nil
}

type num struct {
	sf *sonyflake.Sonyflake
}

func (n *num) Next() (interface{}, error) {
	id, err := n.sf.NextID()
	if err != nil {
		log.Println("Failed to generate ID: " + err.Error())
		return nil, err
	}

	return id, nil
}

func NewGenerator(s Settings) (Generator, error) {
	if s.IsAlphaNumeric {
		hd := hashids.NewData()

		salt := DEFAULT_SALT
		if len(s.Salt) > 0 {
			salt = s.Salt
		}
		hd.Salt = salt
		h, err := hashids.NewWithData(hd)
		if err != nil {
			log.Panic("Failed to initialized ID generator: ", err.Error())
			return nil, err
		}

		return &alphanum{h: h}, nil
	}

	var st sonyflake.Settings
	t, _ := time.Parse("2006-01-02T15:04:05", DEFAULT_SONYFLAKE_TIME)
	if !s.TimeSeed.IsZero() {
		t = s.TimeSeed
	}
	st.StartTime = t

	if s.InDocker {
		st.MachineID = awsutil.AmazonEC2MachineID
	}

	sf := sonyflake.NewSonyflake(st)
	if sf == nil {
		log.Panic("Failed to initialize ID generator (sonyflake)")
		return nil, errors.New("Failed to initialize ID generator (sonyflake)")
	}

	return &num{sf: sf}, nil
}
