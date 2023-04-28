package bigmodelimpl

import (
	"github.com/opensourceways/xihe-server/async-server/domain/bigmodel"
	"github.com/opensourceways/xihe-server/async-server/domain/repository"
	bigmodelapp "github.com/opensourceways/xihe-server/bigmodel/app"
	"github.com/opensourceways/xihe-server/bigmodel/domain"
)

func NewBigModelImpl(s bigmodelapp.AsyncBigModelService) bigmodel.BigModel {
	return &bigmodelImpl{
		srv: s,
	}
}

type bigmodelImpl struct {
	srv bigmodelapp.AsyncBigModelService
}

func (impl *bigmodelImpl) GetIdleEndpoint(bid string) (
	c int, err error,
) {
	return impl.srv.GetIdleEndpoint(bid)
}

func (impl *bigmodelImpl) WuKong(d *repository.WuKongTask) (err error) {
	cmd := bigmodelapp.WuKongCmd{
		WuKongPictureMeta: domain.WuKongPictureMeta{
			Style: d.Style,
			Desc:  d.Desc,
		},

		ImgQuantity: 2,
	}

	return impl.srv.WuKong(d.Id, d.User, &cmd)
}