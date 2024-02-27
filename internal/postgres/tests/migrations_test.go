package tests

import (
	"context"
	"time"

	"github.com/mavryk-network/bcdhub/internal/models/types"
)

func (s *StorageTestSuite) TestMigrationGet() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	migrations, err := s.migrations.Get(ctx, 1)
	s.Require().NoError(err)
	s.Require().Len(migrations, 2)

	m := migrations[0]
	s.Require().EqualValues(4, m.ID)
	s.Require().EqualValues(3, m.ProtocolID)
	s.Require().EqualValues(1, m.PrevProtocolID)
	s.Require().EqualValues(1, m.ContractID)
	s.Require().EqualValues(2, m.Level)
	s.Require().EqualValues(types.MigrationKindUpdate, m.Kind)
}
