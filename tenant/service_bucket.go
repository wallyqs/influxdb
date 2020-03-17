package tenant

import (
	"context"
	"errors"

	"github.com/influxdata/influxdb"
	"github.com/influxdata/influxdb/kv"
)

// FindBucketByID returns a single bucket by ID.
func (s *Service) FindBucketByID(ctx context.Context, id influxdb.ID) (*influxdb.Bucket, error) {
	var bucket *influxdb.Bucket
	err := s.store.View(ctx, func(tx kv.Tx) error {
		b, err := s.store.GetBucket(ctx, tx, id)
		if err != nil {
			return err
		}
		bucket = b
		return nil
	})

	if err != nil {
		return nil, err
	}

	return bucket, nil

}

func (s *Service) FindBucketByName(ctx context.Context, orgID influxdb.ID, name string) (*influxdb.Bucket, error) {
	var bucket *influxdb.Bucket
	err := s.store.View(ctx, func(tx kv.Tx) error {
		b, err := s.store.GetBucketByName(ctx, tx, orgID, name)
		if err != nil {
			return err
		}
		bucket = b
		return nil
	})

	if err != nil {
		return nil, err
	}

	return bucket, nil

}

// FindBucket returns the first bucket that matches filter.
func (s *Service) FindBucket(ctx context.Context, filter influxdb.BucketFilter) (*influxdb.Bucket, error) {
	buckets, _, err := s.FindBuckets(ctx, filter, influxdb.FindOptions{
		Limit: 1,
	})

	if err != nil {
		return nil, err
	}

	if len(buckets) != 1 {
		return nil, errors.New("unkown error")
	}

	return buckets[0], nil
}

// FindBuckets returns a list of buckets that match filter and the total count of matching buckets.
// Additional options provide pagination & sorting.
func (s *Service) FindBuckets(ctx context.Context, filter influxdb.BucketFilter, opt ...influxdb.FindOptions) ([]*influxdb.Bucket, int, error) {
	if filter.ID != nil {
		b, err := s.FindBucketByID(ctx, *filter.ID)
		if err != nil {
			return nil, 0, err
		}
		return []*influxdb.Bucket{b}, 1, nil
	}

	if filter.Name != nil && filter.OrganizationID != nil {
		b, err := s.FindBucketByName(ctx, *filter.OrganizationID, *filter.Name)
		if err != nil {
			return nil, 0, err
		}
		return []*influxdb.Bucket{b}, 1, nil
	}

	var buckets []*influxdb.Bucket
	err := s.store.View(ctx, func(tx kv.Tx) error {
		if filter.OrganizationID == nil && filter.Org != nil {
			org, err := s.store.GetOrgByName(ctx, tx, *filter.Org)
			if err != nil {
				return err
			}
			filter.OrganizationID = &org.ID
		}

		bs, err := s.store.ListBuckets(ctx, tx, BucketFilter{
			Name:           filter.Name,
			OrganizationID: filter.OrganizationID,
		}, opt...)
		if err != nil {
			return err
		}
		buckets = bs
		return nil
	})

	if err != nil {
		return nil, 0, err
	}

	return buckets, len(buckets), nil
}

// CreateBucket creates a new bucket and sets b.ID with the new identifier.
func (s *Service) CreateBucket(ctx context.Context, b *influxdb.Bucket) error {
	if !b.OrgID.Valid() {
		// we need a valid org id
		return ErrOrgNotFound
	}

	return s.store.Update(ctx, func(tx kv.Tx) error {
		// make sure the org exists
		if _, err := s.store.GetOrg(ctx, tx, b.OrgID); err != nil {
			return err
		}

		err := s.store.CreateBucket(ctx, tx, b)
		if err != nil {
			return err
		}

		return s.addOrgRelationToResource(ctx, tx, b.OrgID, b.ID, influxdb.BucketsResourceType)
	})
}

// UpdateBucket updates a single bucket with changeset.
// Returns the new bucket state after update.
func (s *Service) UpdateBucket(ctx context.Context, id influxdb.ID, upd influxdb.BucketUpdate) (*influxdb.Bucket, error) {
	var bucket *influxdb.Bucket
	err := s.store.Update(ctx, func(tx kv.Tx) error {
		b, err := s.store.UpdateBucket(ctx, tx, id, upd)
		if err != nil {
			return err
		}
		bucket = b
		return nil
	})

	if err != nil {
		return nil, err
	}

	return bucket, nil
}

// DeleteBucket removes a bucket by ID.
func (s *Service) DeleteBucket(ctx context.Context, id influxdb.ID) error {
	return s.store.Update(ctx, func(tx kv.Tx) error {
		bucket, err := s.store.GetBucket(ctx, tx, id)
		if err != nil {
			return err
		}
		if bucket.Type == influxdb.BucketTypeSystem {
			// TODO: I think we should allow bucket deletes but maybe im wrong.
			return errDeleteSystemBucket
		}

		if err := s.store.DeleteBucket(ctx, tx, id); err != nil {
			return err
		}
		return s.removeResourceRelations(ctx, tx, id)
	})
}