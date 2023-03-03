package repositoryimpl

import (
	"context"

	"github.com/opensourceways/xihe-server/course/domain"
	"github.com/opensourceways/xihe-server/course/domain/repository"
	"go.mongodb.org/mongo-driver/bson"
)

func NewPlayerRepo(m mongodbClient) repository.Player {
	return &playerRepoImpl{m}
}

type playerRepoImpl struct {
	cli mongodbClient
}

func (impl *playerRepoImpl) SavePlayer(p *domain.Player) (err error) {
	doc, err := impl.genPlayerDoc(p)
	if err != nil {
		return
	}
	f := func(ctx context.Context) error {
		_, err := impl.cli.NewDocIfNotExist(
			ctx, bson.M{
				fieldAccount: p.Account.Account(),
			}, doc,
		)
		return err
	}

	if err = withContext(f); err != nil {
		return
	}

	return
}

func (impl *playerRepoImpl) genPlayerDoc(p *domain.Player) (bson.M, error) {
	obj := DCoursePlayer{
		Id:        p.Id,
		CourseId:  p.CourseId,
		Name:      p.Account.Account(),
		CreatedAt: p.CreatedAt.CourseTime(),
	}

	return genDoc(obj)
}

// Player Count
func (impl *playerRepoImpl) PlayerCount(cid string) (int, error) {
	var v []struct {
		Total int `bson:"total"`
	}

	f := func(ctx context.Context) error {

		pipeline := bson.A{
			bson.M{"$match": bson.M{"$eq": cid}},
			bson.M{"$count": "total"},
		}

		cursor, err := impl.cli.Collection().Aggregate(ctx, pipeline)
		if err != nil {
			return err
		}

		return cursor.All(ctx, &v)
	}

	if err := withContext(f); err != nil || len(v) == 0 {
		return 0, err
	}

	return v[0].Total, nil
}
