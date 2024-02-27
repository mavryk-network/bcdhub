package operations

import (
	"context"

	"github.com/mavryk-network/bcdhub/internal/bcd"
	"github.com/mavryk-network/bcdhub/internal/bcd/consts"
	"github.com/mavryk-network/bcdhub/internal/bcd/encoding"
	"github.com/mavryk-network/bcdhub/internal/helpers"
	"github.com/mavryk-network/bcdhub/internal/noderpc"
	"github.com/mavryk-network/bcdhub/internal/parsers"
)

// OperationParser -
type OperationParser interface {
	Parse(ctx context.Context, data noderpc.Operation, store parsers.Store) error
}

// Group -
type Group struct {
	*ParseParams
}

// NewGroup -
func NewGroup(params *ParseParams) Group {
	return Group{params}
}

// Parse -
func (opg Group) Parse(ctx context.Context, data noderpc.LightOperationGroup, store parsers.Store) error {
	helpers.SetTagSentry("hash", data.Hash)
	if data.Hash != "" {
		hash, err := encoding.DecodeBase58(data.Hash)
		if err != nil {
			return err
		}
		opg.hash = hash
	}

	for idx, item := range data.Contents {
		opg.contentIdx = int64(idx)

		if !opg.needParse(item) {
			continue
		}

		var operation noderpc.Operation
		if err := json.Unmarshal(item.Raw, &operation); err != nil {
			return err
		}

		contentParser := NewContent(opg.ParseParams)
		if err := contentParser.Parse(ctx, operation, store); err != nil {
			return err
		}
		contentParser.clear()
	}

	return nil
}

func (Group) needParse(item noderpc.LightOperation) bool {
	var destination string
	if item.Destination != nil {
		destination = *item.Destination
	}
	prefixCondition := bcd.IsContract(item.Source) || bcd.IsContract(destination)
	transactionCondition := item.Kind == consts.Transaction && prefixCondition
	originationCondition := (item.Kind == consts.Origination || item.Kind == consts.OriginationNew || item.Kind == consts.TxRollupOrigination)
	registerGlobalConstantCondition := item.Kind == consts.RegisterGlobalConstant
	eventCondition := item.Kind == consts.Event
	transferTicketCondition := item.Kind == consts.TransferTicket
	srCondition := item.Kind == consts.SrOriginate || item.Kind == consts.SrExecuteOutboxMessage
	return originationCondition || transactionCondition || srCondition ||
		registerGlobalConstantCondition || eventCondition || transferTicketCondition
}

// Content -
type Content struct {
	*ParseParams
}

// NewContent -
func NewContent(params *ParseParams) Content {
	return Content{params}
}

// Parse -
func (content Content) Parse(ctx context.Context, data noderpc.Operation, store parsers.Store) error {
	var operationParser OperationParser
	switch data.Kind {
	case consts.Origination, consts.OriginationNew:
		operationParser = NewOrigination(content.ParseParams)
	case consts.Transaction:
		operationParser = NewTransaction(content.ParseParams)
	case consts.RegisterGlobalConstant:
		operationParser = NewRegisterGlobalConstant(content.ParseParams)
	case consts.TxRollupOrigination:
		operationParser = NewTxRollupOrigination(content.ParseParams)
	case consts.Event:
		operationParser = NewEvent(content.ParseParams)
	case consts.TransferTicket:
		operationParser = NewTransferTicket(content.ParseParams)
	case consts.SrOriginate:
		operationParser = NewSrOriginate(content.ParseParams)
	case consts.SrExecuteOutboxMessage:
		operationParser = NewSrExecuteOutboxMessage(content.ParseParams)
	default:
		return nil
	}
	if err := operationParser.Parse(ctx, data, store); err != nil {
		return err
	}

	if err := content.parseInternal(ctx, data, store); err != nil {
		return err
	}

	return nil
}

func (content Content) parseInternal(ctx context.Context, data noderpc.Operation, store parsers.Store) error {
	if data.Metadata == nil {
		return nil
	}
	internals := data.Metadata.Internal
	if internals == nil {
		internals = data.Metadata.InternalOperations
		if internals == nil {
			return nil
		}
	}

	for i := range internals {
		if err := content.Parse(ctx, internals[i], store); err != nil {
			return err
		}
	}
	return nil
}

func (content *Content) clear() {
	content.main = nil
}
