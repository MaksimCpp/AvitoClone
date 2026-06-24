package cached

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/MaksimCpp/AvitoClone/internal/domain/item"
	avitoredis "github.com/MaksimCpp/AvitoClone/internal/repository/redis"
)

type CachedItemRepository struct {
	repo item.ItemRepository
	cache *avitoredis.Cache
}

func NewCachedItemRepository(
	repo item.ItemRepository,
	cache *avitoredis.Cache,
) *CachedItemRepository {
	return &CachedItemRepository{
		repo: repo,
		cache: cache,
	}
}

func (сrepo *CachedItemRepository) Create(
	ctx context.Context, itemEntity *item.Item,
) error {
	err := сrepo.repo.Create(ctx, itemEntity)
	if err != nil {
		return err
	}

	data, err := json.Marshal(itemEntity)
	if err != nil {
		return err
	}

	key := fmt.Sprintf(
		"user_%d",
		itemEntity.ID,
	)
	_, err = сrepo.cache.Set(
		ctx, key, data, time.Minute * 10,
	)
	
	return err
}

func (crepo *CachedItemRepository) GetByID(
	ctx context.Context, id int,
) (*item.Item, error) {
	key := fmt.Sprintf(
		"user_%d",
		id,
	)

	data, err := crepo.cache.Get(ctx, key)

	if err == nil {
		var itemEntity item.Item
		err = json.Unmarshal(data, &itemEntity)

		if err != nil {
			return nil, err
		}

		return &itemEntity, nil
	}

	return crepo.repo.GetByID(ctx, id)
}

func (crepo *CachedItemRepository) List(
	ctx context.Context, limit int, offset int,
) ([]*item.Item, error) {
	return crepo.repo.List(ctx, limit, offset)
}

func (crepo *CachedItemRepository) ListByUserID(
	ctx context.Context, userID int,
) ([]*item.Item, error) {
	return crepo.repo.ListByUserID(ctx, userID)
}

func (crepo *CachedItemRepository) Delete(ctx context.Context, id int) error {
	err := crepo.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	key := fmt.Sprintf(
		"user_%d", id,
	)

	return crepo.cache.Delete(ctx, key)
}
