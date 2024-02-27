package operations

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/mavryk-network/bcdhub/internal/bcd"
	"github.com/mavryk-network/bcdhub/internal/bcd/consts"
	astContract "github.com/mavryk-network/bcdhub/internal/bcd/contract"
	"github.com/mavryk-network/bcdhub/internal/models/account"
	"github.com/mavryk-network/bcdhub/internal/models/bigmapaction"
	"github.com/mavryk-network/bcdhub/internal/models/bigmapdiff"
	"github.com/mavryk-network/bcdhub/internal/models/contract"
	"github.com/mavryk-network/bcdhub/internal/models/operation"
	"github.com/mavryk-network/bcdhub/internal/models/ticket"
	"github.com/mavryk-network/bcdhub/internal/models/types"
	"github.com/mavryk-network/bcdhub/internal/noderpc"
	"github.com/mavryk-network/bcdhub/internal/parsers"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func readJSONFile(name string, response interface{}) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewDecoder(f).Decode(response)
}

func readTestScript(address, symLink string) ([]byte, error) {
	path := filepath.Join("./test/contracts", fmt.Sprintf("%s_%s.json", address, symLink))
	return os.ReadFile(path)
}

func readRPCScript(_ context.Context, address string, _ int64) (noderpc.Script, error) {
	var script noderpc.Script
	storageFile := fmt.Sprintf("./data/rpc/script/script/%s.json", address)
	if _, err := os.Lstat(storageFile); !os.IsNotExist(err) {
		f, err := os.Open(storageFile)
		if err != nil {
			return script, err
		}
		defer f.Close()

		err = json.NewDecoder(f).Decode(&script)
		return script, err
	}
	return script, errors.Errorf("unknown RPC script: %s", address)
}

func readTestScriptModel(_ context.Context, address, symLink string) (contract.Script, error) {
	data, err := readTestScript(address, bcd.SymLinkBabylon)
	if err != nil {
		return contract.Script{}, err
	}
	var buffer bytes.Buffer
	buffer.WriteString(`{"code":`)
	buffer.Write(data)
	buffer.WriteString(`,"storage":{}}`)
	script, err := astContract.NewParser(buffer.Bytes())
	if err != nil {
		return contract.Script{}, errors.Wrap(err, "astContract.NewParser")
	}
	if err := script.Parse(); err != nil {
		return contract.Script{}, err
	}
	var s bcd.RawScript
	if err := json.Unmarshal(data, &s); err != nil {
		return contract.Script{}, err
	}
	return contract.Script{
		Code:        s.Code,
		Parameter:   s.Parameter,
		Storage:     s.Storage,
		Hash:        script.Hash,
		FailStrings: script.FailStrings.Values(),
		Annotations: script.Annotations.Values(),
		Tags:        types.NewTags(script.Tags.Values()),
		Hardcoded:   script.HardcodedAddresses.Values(),
	}, nil
}

func readTestScriptPart(_ context.Context, address, symLink, part string) ([]byte, error) {
	data, err := readTestScript(address, bcd.SymLinkBabylon)
	if err != nil {
		return nil, err
	}
	var s bcd.RawScript
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}

	switch part {
	case consts.CODE:
		return s.Code, nil
	case consts.PARAMETER:
		return s.Parameter, nil
	case consts.STORAGE:
		return s.Storage, nil
	}
	return nil, nil
}

func readTestContractModel(_ context.Context, address string) (contract.Contract, error) {
	var c contract.Contract
	f, err := os.Open(fmt.Sprintf("./data/models/contract/%s.json", address))
	if err != nil {
		return c, err
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&c)
	return c, err
}

func compareParserResponse(t *testing.T, got, want *parsers.TestStore) {
	require.Len(t, got.BigMapState, len(want.BigMapState))
	require.Len(t, got.Contracts, len(want.Contracts))
	require.Len(t, got.Migrations, len(want.Migrations))
	require.Len(t, got.Operations, len(want.Operations))
	require.Len(t, got.GlobalConstants, len(want.GlobalConstants))
	require.Len(t, got.Tickets, len(want.Tickets))

	for i := range got.Contracts {
		compareContract(t, want.Contracts[i], got.Contracts[i])
	}
	for i := range got.Migrations {
		require.Equal(t, want.Migrations[i], got.Migrations[i])
	}
	for i := range got.Operations {
		compareOperations(t, want.Operations[i], got.Operations[i])
	}
	for i := range got.BigMapState {
		require.Equal(t, want.BigMapState[i], got.BigMapState[i])
	}
	for i := range got.GlobalConstants {
		require.Equal(t, want.GlobalConstants[i], got.GlobalConstants[i])
	}
	for hash := range got.Tickets {
		require.Equal(t, want.Tickets[hash], got.Tickets[hash])
	}
	for key, wantAddress := range want.Accounts {
		gotAddress, ok := got.Accounts[key]
		require.True(t, ok)
		require.Equal(t, wantAddress, gotAddress)
	}
}

func compareAddress(t *testing.T, want, got account.Account) {
	if want.Address == "" && got.Address == "" {
		return
	}

	require.Equal(t, want, got)
}

func compareOperations(t *testing.T, want, got *operation.Operation) {
	require.EqualValues(t, want.Internal, got.Internal)
	compareInt64Ptr(t, want.Nonce, got.Nonce)
	require.EqualValues(t, want.Timestamp, got.Timestamp)
	require.EqualValues(t, want.Level, got.Level)
	require.EqualValues(t, want.ContentIndex, got.ContentIndex)
	require.EqualValues(t, want.Counter, got.Counter)
	require.EqualValues(t, want.GasLimit, got.GasLimit)
	require.EqualValues(t, want.StorageLimit, got.StorageLimit)
	require.EqualValues(t, want.Fee, got.Fee)
	require.EqualValues(t, want.Amount, got.Amount)
	require.EqualValues(t, want.Burned, got.Burned)
	require.EqualValues(t, want.AllocatedDestinationContractBurned, got.AllocatedDestinationContractBurned)
	require.EqualValues(t, want.ProtocolID, got.ProtocolID)
	require.Equal(t, want.Hash, got.Hash)
	require.EqualValues(t, want.Status, got.Status)
	require.EqualValues(t, want.Kind, got.Kind)
	compareAddress(t, want.Initiator, got.Initiator)
	compareAddress(t, want.Source, got.Source)
	compareAddress(t, want.Destination, got.Destination)
	compareAddress(t, want.Delegate, got.Delegate)
	require.Equal(t, want.Entrypoint, got.Entrypoint)
	compareBytesArray(t, want.Parameters, got.Parameters)
	compareBytesArray(t, want.DeffatedStorage, got.DeffatedStorage)
	require.EqualValues(t, want.Tags, got.Tags)
	require.Len(t, got.BigMapDiffs, len(want.BigMapDiffs))
	require.Len(t, got.BigMapActions, len(want.BigMapActions))
	require.Len(t, got.TicketUpdates, len(want.TicketUpdates))

	for i := range want.BigMapDiffs {
		compareBigMapDiff(t, want.BigMapDiffs[i], got.BigMapDiffs[i])
	}

	for i := range want.BigMapActions {
		compareBigMapAction(t, want.BigMapActions[i], got.BigMapActions[i])
	}

	for i := range want.TicketUpdates {
		compareTicketUpdates(t, want.TicketUpdates[i], got.TicketUpdates[i])
	}
}

func compareTicketUpdates(t *testing.T, want, got *ticket.TicketUpdate) {
	require.EqualValues(t, want.Account, got.Account)
	require.EqualValues(t, want.Ticket, got.Ticket)
	require.EqualValues(t, want.Level, got.Level)
	require.EqualValues(t, want.Timestamp, got.Timestamp)
}

func compareBigMapDiff(t *testing.T, want, got *bigmapdiff.BigMapDiff) {
	require.EqualValues(t, want.Contract, got.Contract)
	require.EqualValues(t, want.KeyHash, got.KeyHash)
	require.EqualValues(t, want.Level, got.Level)
	require.EqualValues(t, want.Timestamp, got.Timestamp)
	require.EqualValues(t, want.ProtocolID, got.ProtocolID)
	require.EqualValues(t, want.Ptr, got.Ptr)
	compareBytesArray(t, want.KeyBytes(), got.KeyBytes())
	compareBytesArray(t, want.ValueBytes(), got.ValueBytes())
}

func compareBytesArray(t *testing.T, want, got []byte) {
	if len(want) > 0 {
		require.JSONEq(t, string(want), string(got))
	}
}

func compareBigMapAction(t *testing.T, want, got *bigmapaction.BigMapAction) {
	require.EqualValues(t, want.Action, got.Action)
	compareInt64Ptr(t, want.SourcePtr, got.SourcePtr)
	compareInt64Ptr(t, want.DestinationPtr, got.DestinationPtr)
	require.EqualValues(t, want.Level, got.Level)
	require.EqualValues(t, want.Address, got.Address)
	require.EqualValues(t, want.Timestamp, got.Timestamp)
}

func compareContract(t *testing.T, want, got *contract.Contract) {
	require.Equal(t, want.Account, got.Account)
	require.Equal(t, want.Manager, got.Manager)
	require.Equal(t, want.Level, got.Level)
	require.Equal(t, want.Timestamp, got.Timestamp)
	require.Equal(t, want.Tags, got.Tags)
	compareScript(t, want.Alpha, got.Alpha)
	compareScript(t, want.Babylon, got.Babylon)
}

func compareScript(t *testing.T, want, got contract.Script) {
	require.Equal(t, want.Hash, got.Hash)
	require.ElementsMatch(t, want.Entrypoints, got.Entrypoints)
	require.ElementsMatch(t, want.Annotations, got.Annotations)
	require.ElementsMatch(t, want.FailStrings, got.FailStrings)
	require.ElementsMatch(t, want.Hardcoded, got.Hardcoded)
	require.ElementsMatch(t, want.Code, got.Code)
}

func compareInt64Ptr(t *testing.T, want, got *int64) {
	require.Condition(t, func() (success bool) {
		return (want != nil && got != nil && *want == *got) || (want == nil && got == nil)
	})
}
