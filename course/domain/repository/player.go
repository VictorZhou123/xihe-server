package repository

import "github.com/opensourceways/xihe-server/course/domain"

type Player interface {
	SavePlayer(*domain.Player) error
	PlayerCount(cid string) (int, error)
}
