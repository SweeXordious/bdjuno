package utils

import (
	"context"

	"github.com/desmos-labs/juno/client"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/types"

	juno "github.com/desmos-labs/juno/types"
	iscntypes "github.com/likecoin/likechain/x/iscn/types"
)

// StoreIscnRecordFromMessage handles storing new iscn record inside the database
func StoreIscnRecordFromMessage(
	height int64, tx *juno.Tx, index int, msg *iscntypes.MsgCreateIscnRecord, iscnClient iscntypes.QueryClient, db *database.Db,
) error {

	event, err := tx.FindEventByType(index, iscntypes.EventTypeIscnRecord)
	if err != nil {
		return err
	}

	id, err := tx.FindAttributeByKey(event, iscntypes.AttributeKeyIscnId)
	if err != nil {
		return err
	}

	// Get the record
	res, err := iscnClient.RecordsById(
		context.Background(),
		&iscntypes.QueryRecordsByIdRequest{IscnId: id},
		client.GetHeightRequestHeader(height),
	)
	if err != nil {
		return err
	}
	iscnRecord := types.NewRecord(id, msg.Record.RecordNotes, msg.Record.ContentFingerprints, msg.Record.Stakeholders, msg.Record.ContentMetadata)
	newIscnRecord := types.NewIscnRecord(res.Owner, id, res.LatestVersion, res.Records[0].Ipld, iscnRecord, height)
	return db.SaveIscnRecord(newIscnRecord)
}

// UpdateIscnRecordFromMessage handles updating the existing iscn record inside the database
func UpdateIscnRecordFromMessage(
	height int64, tx *juno.Tx, index int, msg *iscntypes.MsgUpdateIscnRecord, iscnClient iscntypes.QueryClient, db *database.Db,
) error {

	id := msg.IscnId

	// Get the record
	res, err := iscnClient.RecordsById(
		context.Background(),
		&iscntypes.QueryRecordsByIdRequest{IscnId: id},
		client.GetHeightRequestHeader(height),
	)
	if err != nil {
		return err
	}
	iscnRecord := types.NewRecord(id, msg.Record.RecordNotes, msg.Record.ContentFingerprints, msg.Record.Stakeholders, msg.Record.ContentMetadata)
	newIscnRecord := types.NewIscnRecord(res.Owner, id, res.LatestVersion, res.Records[0].Ipld, iscnRecord, height)
	return db.UpdateIscnRecord(newIscnRecord)
}

// UpdateIscnRecordOwnershipFromMessage handles updating ownership of the existing iscn record inside the database
func UpdateIscnRecordOwnershipFromMessage(
	height int64, tx *juno.Tx, index int, msg *iscntypes.MsgChangeIscnRecordOwnership, iscnClient iscntypes.QueryClient, db *database.Db,
) error {
	id := msg.IscnId
	updatedIscnRecord := types.NewIscnChangeOwnership(msg.From, id, msg.NewOwner)
	return db.UpdateIscnRecordOwnership(updatedIscnRecord)
}
