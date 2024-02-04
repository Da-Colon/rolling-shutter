// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: keyper.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/jackc/pgconn"
)

const countBatchConfigs = `-- name: CountBatchConfigs :one
SELECT count(*) FROM tendermint_batch_config
`

func (q *Queries) CountBatchConfigs(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countBatchConfigs)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countBatchConfigsInBlockRange = `-- name: CountBatchConfigsInBlockRange :one
SELECT COUNT(*)
FROM tendermint_batch_config
WHERE $1 <= activation_block_number AND activation_block_number < $2
`

type CountBatchConfigsInBlockRangeParams struct {
	StartBlock int64
	EndBlock   int64
}

func (q *Queries) CountBatchConfigsInBlockRange(ctx context.Context, arg CountBatchConfigsInBlockRangeParams) (int64, error) {
	row := q.db.QueryRow(ctx, countBatchConfigsInBlockRange, arg.StartBlock, arg.EndBlock)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countDecryptionKeyShares = `-- name: CountDecryptionKeyShares :one
SELECT count(*) FROM decryption_key_share
WHERE eon = $1 AND epoch_id = $2
`

type CountDecryptionKeySharesParams struct {
	Eon     int64
	EpochID []byte
}

func (q *Queries) CountDecryptionKeyShares(ctx context.Context, arg CountDecryptionKeySharesParams) (int64, error) {
	row := q.db.QueryRow(ctx, countDecryptionKeyShares, arg.Eon, arg.EpochID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const deletePolyEval = `-- name: DeletePolyEval :exec

DELETE FROM poly_evals ev WHERE ev.eon=$1 AND ev.receiver_address=$2
`

type DeletePolyEvalParams struct {
	Eon             int64
	ReceiverAddress string
}

// PolyEvalsWithEncryptionKeys could probably already delete the entries from the poly_evals table.
// I wasn't able to make this work, because of bugs in sqlc
func (q *Queries) DeletePolyEval(ctx context.Context, arg DeletePolyEvalParams) error {
	_, err := q.db.Exec(ctx, deletePolyEval, arg.Eon, arg.ReceiverAddress)
	return err
}

const deletePolyEvalByEon = `-- name: DeletePolyEvalByEon :execresult
DELETE FROM poly_evals ev WHERE ev.eon=$1
`

func (q *Queries) DeletePolyEvalByEon(ctx context.Context, eon int64) (pgconn.CommandTag, error) {
	return q.db.Exec(ctx, deletePolyEvalByEon, eon)
}

const deletePureDKG = `-- name: DeletePureDKG :exec
DELETE FROM puredkg WHERE eon=$1
`

func (q *Queries) DeletePureDKG(ctx context.Context, eon int64) error {
	_, err := q.db.Exec(ctx, deletePureDKG, eon)
	return err
}

const deleteShutterMessage = `-- name: DeleteShutterMessage :exec
DELETE FROM tendermint_outgoing_messages WHERE id=$1
`

func (q *Queries) DeleteShutterMessage(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteShutterMessage, id)
	return err
}

const existsDecryptionKey = `-- name: ExistsDecryptionKey :one
SELECT EXISTS (
    SELECT 1
    FROM decryption_key
    WHERE eon = $1 AND epoch_id = $2
)
`

type ExistsDecryptionKeyParams struct {
	Eon     int64
	EpochID []byte
}

func (q *Queries) ExistsDecryptionKey(ctx context.Context, arg ExistsDecryptionKeyParams) (bool, error) {
	row := q.db.QueryRow(ctx, existsDecryptionKey, arg.Eon, arg.EpochID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const existsDecryptionKeyShare = `-- name: ExistsDecryptionKeyShare :one
SELECT EXISTS (
    SELECT 1
    FROM decryption_key_share
    WHERE eon = $1 AND epoch_id = $2 AND keyper_index = $3
)
`

type ExistsDecryptionKeyShareParams struct {
	Eon         int64
	EpochID     []byte
	KeyperIndex int64
}

func (q *Queries) ExistsDecryptionKeyShare(ctx context.Context, arg ExistsDecryptionKeyShareParams) (bool, error) {
	row := q.db.QueryRow(ctx, existsDecryptionKeyShare, arg.Eon, arg.EpochID, arg.KeyperIndex)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const getAllEons = `-- name: GetAllEons :many
SELECT eon, height, activation_block_number, keyper_config_index FROM eons ORDER BY eon
`

func (q *Queries) GetAllEons(ctx context.Context) ([]Eon, error) {
	rows, err := q.db.Query(ctx, getAllEons)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Eon
	for rows.Next() {
		var i Eon
		if err := rows.Scan(
			&i.Eon,
			&i.Height,
			&i.ActivationBlockNumber,
			&i.KeyperConfigIndex,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAndDeleteEonPublicKeys = `-- name: GetAndDeleteEonPublicKeys :many
WITH t1 AS (DELETE FROM outgoing_eon_keys RETURNING eon_public_key, eon)
SELECT t1.eon_public_key, t1.eon, eons.activation_block_number, tbc.keypers, tbc.keyper_config_index
FROM t1
INNER JOIN eons
      ON t1.eon = eons.eon
INNER JOIN tendermint_batch_config tbc
      ON eons.keyper_config_index = tbc.keyper_config_index
`

type GetAndDeleteEonPublicKeysRow struct {
	EonPublicKey          []byte
	Eon                   int64
	ActivationBlockNumber int64
	Keypers               []string
	KeyperConfigIndex     int32
}

func (q *Queries) GetAndDeleteEonPublicKeys(ctx context.Context) ([]GetAndDeleteEonPublicKeysRow, error) {
	rows, err := q.db.Query(ctx, getAndDeleteEonPublicKeys)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAndDeleteEonPublicKeysRow
	for rows.Next() {
		var i GetAndDeleteEonPublicKeysRow
		if err := rows.Scan(
			&i.EonPublicKey,
			&i.Eon,
			&i.ActivationBlockNumber,
			&i.Keypers,
			&i.KeyperConfigIndex,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getBatchConfig = `-- name: GetBatchConfig :one
SELECT keyper_config_index, height, keypers, threshold, started, activation_block_number
FROM tendermint_batch_config
WHERE keyper_config_index = $1
`

func (q *Queries) GetBatchConfig(ctx context.Context, keyperConfigIndex int32) (TendermintBatchConfig, error) {
	row := q.db.QueryRow(ctx, getBatchConfig, keyperConfigIndex)
	var i TendermintBatchConfig
	err := row.Scan(
		&i.KeyperConfigIndex,
		&i.Height,
		&i.Keypers,
		&i.Threshold,
		&i.Started,
		&i.ActivationBlockNumber,
	)
	return i, err
}

const getBatchConfigs = `-- name: GetBatchConfigs :many
SELECT keyper_config_index, height, keypers, threshold, started, activation_block_number
FROM tendermint_batch_config
ORDER BY keyper_config_index
`

func (q *Queries) GetBatchConfigs(ctx context.Context) ([]TendermintBatchConfig, error) {
	rows, err := q.db.Query(ctx, getBatchConfigs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []TendermintBatchConfig
	for rows.Next() {
		var i TendermintBatchConfig
		if err := rows.Scan(
			&i.KeyperConfigIndex,
			&i.Height,
			&i.Keypers,
			&i.Threshold,
			&i.Started,
			&i.ActivationBlockNumber,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getDKGResult = `-- name: GetDKGResult :one
SELECT eon, success, error, pure_result FROM dkg_result
WHERE eon = $1
`

func (q *Queries) GetDKGResult(ctx context.Context, eon int64) (DkgResult, error) {
	row := q.db.QueryRow(ctx, getDKGResult, eon)
	var i DkgResult
	err := row.Scan(
		&i.Eon,
		&i.Success,
		&i.Error,
		&i.PureResult,
	)
	return i, err
}

const getDKGResultForBlockNumber = `-- name: GetDKGResultForBlockNumber :one
SELECT eon, success, error, pure_result FROM dkg_result
WHERE eon = (SELECT eon FROM eons WHERE activation_block_number <= $1
ORDER BY activation_block_number DESC, height DESC
LIMIT 1)
`

func (q *Queries) GetDKGResultForBlockNumber(ctx context.Context, blockNumber int64) (DkgResult, error) {
	row := q.db.QueryRow(ctx, getDKGResultForBlockNumber, blockNumber)
	var i DkgResult
	err := row.Scan(
		&i.Eon,
		&i.Success,
		&i.Error,
		&i.PureResult,
	)
	return i, err
}

const getDecryptionKey = `-- name: GetDecryptionKey :one
SELECT eon, epoch_id, decryption_key FROM decryption_key
WHERE eon = $1 AND epoch_id = $2
`

type GetDecryptionKeyParams struct {
	Eon     int64
	EpochID []byte
}

func (q *Queries) GetDecryptionKey(ctx context.Context, arg GetDecryptionKeyParams) (DecryptionKey, error) {
	row := q.db.QueryRow(ctx, getDecryptionKey, arg.Eon, arg.EpochID)
	var i DecryptionKey
	err := row.Scan(&i.Eon, &i.EpochID, &i.DecryptionKey)
	return i, err
}

const getDecryptionKeyShare = `-- name: GetDecryptionKeyShare :one
SELECT eon, epoch_id, keyper_index, decryption_key_share FROM decryption_key_share
WHERE eon = $1 AND epoch_id = $2 AND keyper_index = $3
`

type GetDecryptionKeyShareParams struct {
	Eon         int64
	EpochID     []byte
	KeyperIndex int64
}

func (q *Queries) GetDecryptionKeyShare(ctx context.Context, arg GetDecryptionKeyShareParams) (DecryptionKeyShare, error) {
	row := q.db.QueryRow(ctx, getDecryptionKeyShare, arg.Eon, arg.EpochID, arg.KeyperIndex)
	var i DecryptionKeyShare
	err := row.Scan(
		&i.Eon,
		&i.EpochID,
		&i.KeyperIndex,
		&i.DecryptionKeyShare,
	)
	return i, err
}

const getEncryptionKeys = `-- name: GetEncryptionKeys :many
SELECT address, encryption_public_key FROM tendermint_encryption_key
`

func (q *Queries) GetEncryptionKeys(ctx context.Context) ([]TendermintEncryptionKey, error) {
	rows, err := q.db.Query(ctx, getEncryptionKeys)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []TendermintEncryptionKey
	for rows.Next() {
		var i TendermintEncryptionKey
		if err := rows.Scan(&i.Address, &i.EncryptionPublicKey); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEon = `-- name: GetEon :one
SELECT eon, height, activation_block_number, keyper_config_index FROM eons WHERE eon=$1
`

func (q *Queries) GetEon(ctx context.Context, eon int64) (Eon, error) {
	row := q.db.QueryRow(ctx, getEon, eon)
	var i Eon
	err := row.Scan(
		&i.Eon,
		&i.Height,
		&i.ActivationBlockNumber,
		&i.KeyperConfigIndex,
	)
	return i, err
}

const getEonForBlockNumber = `-- name: GetEonForBlockNumber :one
SELECT eon, height, activation_block_number, keyper_config_index FROM eons
WHERE activation_block_number <= $1
ORDER BY activation_block_number DESC, height DESC
LIMIT 1
`

func (q *Queries) GetEonForBlockNumber(ctx context.Context, blockNumber int64) (Eon, error) {
	row := q.db.QueryRow(ctx, getEonForBlockNumber, blockNumber)
	var i Eon
	err := row.Scan(
		&i.Eon,
		&i.Height,
		&i.ActivationBlockNumber,
		&i.KeyperConfigIndex,
	)
	return i, err
}

const getLastBatchConfigSent = `-- name: GetLastBatchConfigSent :one
SELECT keyper_config_index FROM last_batch_config_sent LIMIT 1
`

func (q *Queries) GetLastBatchConfigSent(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, getLastBatchConfigSent)
	var keyper_config_index int64
	err := row.Scan(&keyper_config_index)
	return keyper_config_index, err
}

const getLastBlockSeen = `-- name: GetLastBlockSeen :one
SELECT block_number FROM last_block_seen LIMIT 1
`

func (q *Queries) GetLastBlockSeen(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, getLastBlockSeen)
	var block_number int64
	err := row.Scan(&block_number)
	return block_number, err
}

const getLastCommittedHeight = `-- name: GetLastCommittedHeight :one
SELECT last_committed_height
FROM tendermint_sync_meta
ORDER BY current_block DESC, last_committed_height DESC
LIMIT 1
`

func (q *Queries) GetLastCommittedHeight(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, getLastCommittedHeight)
	var last_committed_height int64
	err := row.Scan(&last_committed_height)
	return last_committed_height, err
}

const getLatestBatchConfig = `-- name: GetLatestBatchConfig :one
SELECT keyper_config_index, height, keypers, threshold, started, activation_block_number
FROM tendermint_batch_config
ORDER BY keyper_config_index DESC
LIMIT 1
`

func (q *Queries) GetLatestBatchConfig(ctx context.Context) (TendermintBatchConfig, error) {
	row := q.db.QueryRow(ctx, getLatestBatchConfig)
	var i TendermintBatchConfig
	err := row.Scan(
		&i.KeyperConfigIndex,
		&i.Height,
		&i.Keypers,
		&i.Threshold,
		&i.Started,
		&i.ActivationBlockNumber,
	)
	return i, err
}

const getNextShutterMessage = `-- name: GetNextShutterMessage :one
SELECT id, description, msg from tendermint_outgoing_messages
ORDER BY id
LIMIT 1
`

func (q *Queries) GetNextShutterMessage(ctx context.Context) (TendermintOutgoingMessage, error) {
	row := q.db.QueryRow(ctx, getNextShutterMessage)
	var i TendermintOutgoingMessage
	err := row.Scan(&i.ID, &i.Description, &i.Msg)
	return i, err
}

const insertBatchConfig = `-- name: InsertBatchConfig :exec
INSERT INTO tendermint_batch_config (keyper_config_index, height, keypers, threshold, started, activation_block_number)
VALUES ($1, $2, $3, $4, $5, $6)
`

type InsertBatchConfigParams struct {
	KeyperConfigIndex     int32
	Height                int64
	Keypers               []string
	Threshold             int32
	Started               bool
	ActivationBlockNumber int64
}

func (q *Queries) InsertBatchConfig(ctx context.Context, arg InsertBatchConfigParams) error {
	_, err := q.db.Exec(ctx, insertBatchConfig,
		arg.KeyperConfigIndex,
		arg.Height,
		arg.Keypers,
		arg.Threshold,
		arg.Started,
		arg.ActivationBlockNumber,
	)
	return err
}

const insertDKGResult = `-- name: InsertDKGResult :exec
INSERT INTO dkg_result (eon,success,error,pure_result)
VALUES ($1,$2,$3,$4)
`

type InsertDKGResultParams struct {
	Eon        int64
	Success    bool
	Error      sql.NullString
	PureResult []byte
}

func (q *Queries) InsertDKGResult(ctx context.Context, arg InsertDKGResultParams) error {
	_, err := q.db.Exec(ctx, insertDKGResult,
		arg.Eon,
		arg.Success,
		arg.Error,
		arg.PureResult,
	)
	return err
}

const insertDecryptionKey = `-- name: InsertDecryptionKey :execresult
INSERT INTO decryption_key (eon, epoch_id, decryption_key)
VALUES ($1, $2, $3)
ON CONFLICT DO NOTHING
`

type InsertDecryptionKeyParams struct {
	Eon           int64
	EpochID       []byte
	DecryptionKey []byte
}

func (q *Queries) InsertDecryptionKey(ctx context.Context, arg InsertDecryptionKeyParams) (pgconn.CommandTag, error) {
	return q.db.Exec(ctx, insertDecryptionKey, arg.Eon, arg.EpochID, arg.DecryptionKey)
}

const insertDecryptionKeyShare = `-- name: InsertDecryptionKeyShare :exec
INSERT INTO decryption_key_share (eon, epoch_id, keyper_index, decryption_key_share)
VALUES ($1, $2, $3, $4)
`

type InsertDecryptionKeyShareParams struct {
	Eon                int64
	EpochID            []byte
	KeyperIndex        int64
	DecryptionKeyShare []byte
}

func (q *Queries) InsertDecryptionKeyShare(ctx context.Context, arg InsertDecryptionKeyShareParams) error {
	_, err := q.db.Exec(ctx, insertDecryptionKeyShare,
		arg.Eon,
		arg.EpochID,
		arg.KeyperIndex,
		arg.DecryptionKeyShare,
	)
	return err
}

const insertEncryptionKey = `-- name: InsertEncryptionKey :exec
INSERT INTO tendermint_encryption_key (address, encryption_public_key) VALUES ($1, $2)
`

type InsertEncryptionKeyParams struct {
	Address             string
	EncryptionPublicKey []byte
}

func (q *Queries) InsertEncryptionKey(ctx context.Context, arg InsertEncryptionKeyParams) error {
	_, err := q.db.Exec(ctx, insertEncryptionKey, arg.Address, arg.EncryptionPublicKey)
	return err
}

const insertEon = `-- name: InsertEon :exec
INSERT INTO eons (eon, height, activation_block_number, keyper_config_index)
VALUES ($1, $2, $3, $4)
`

type InsertEonParams struct {
	Eon                   int64
	Height                int64
	ActivationBlockNumber int64
	KeyperConfigIndex     int64
}

func (q *Queries) InsertEon(ctx context.Context, arg InsertEonParams) error {
	_, err := q.db.Exec(ctx, insertEon,
		arg.Eon,
		arg.Height,
		arg.ActivationBlockNumber,
		arg.KeyperConfigIndex,
	)
	return err
}

const insertEonPublicKey = `-- name: InsertEonPublicKey :exec
INSERT INTO outgoing_eon_keys (eon_public_key, eon)
VALUES ($1, $2)
`

type InsertEonPublicKeyParams struct {
	EonPublicKey []byte
	Eon          int64
}

func (q *Queries) InsertEonPublicKey(ctx context.Context, arg InsertEonPublicKeyParams) error {
	_, err := q.db.Exec(ctx, insertEonPublicKey, arg.EonPublicKey, arg.Eon)
	return err
}

const insertPolyEval = `-- name: InsertPolyEval :exec
INSERT INTO poly_evals (eon, receiver_address, eval)
VALUES ($1, $2, $3)
`

type InsertPolyEvalParams struct {
	Eon             int64
	ReceiverAddress string
	Eval            []byte
}

func (q *Queries) InsertPolyEval(ctx context.Context, arg InsertPolyEvalParams) error {
	_, err := q.db.Exec(ctx, insertPolyEval, arg.Eon, arg.ReceiverAddress, arg.Eval)
	return err
}

const insertPureDKG = `-- name: InsertPureDKG :exec
INSERT INTO puredkg (eon, puredkg) VALUES ($1, $2)
ON CONFLICT (eon) DO UPDATE SET puredkg=EXCLUDED.puredkg
`

type InsertPureDKGParams struct {
	Eon     int64
	Puredkg []byte
}

func (q *Queries) InsertPureDKG(ctx context.Context, arg InsertPureDKGParams) error {
	_, err := q.db.Exec(ctx, insertPureDKG, arg.Eon, arg.Puredkg)
	return err
}

const polyEvalsWithEncryptionKeys = `-- name: PolyEvalsWithEncryptionKeys :many
SELECT ev.eon, ev.receiver_address, ev.eval,
       k.encryption_public_key,
       eons.height
FROM poly_evals ev
INNER JOIN tendermint_encryption_key k
      ON ev.receiver_address = k.address
INNER JOIN eons eons
      ON ev.eon = eons.eon
ORDER BY ev.eon
`

type PolyEvalsWithEncryptionKeysRow struct {
	Eon                 int64
	ReceiverAddress     string
	Eval                []byte
	EncryptionPublicKey []byte
	Height              int64
}

func (q *Queries) PolyEvalsWithEncryptionKeys(ctx context.Context) ([]PolyEvalsWithEncryptionKeysRow, error) {
	rows, err := q.db.Query(ctx, polyEvalsWithEncryptionKeys)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PolyEvalsWithEncryptionKeysRow
	for rows.Next() {
		var i PolyEvalsWithEncryptionKeysRow
		if err := rows.Scan(
			&i.Eon,
			&i.ReceiverAddress,
			&i.Eval,
			&i.EncryptionPublicKey,
			&i.Height,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const scheduleSerializedShutterMessage = `-- name: ScheduleSerializedShutterMessage :one
INSERT INTO tendermint_outgoing_messages (description, msg)
VALUES ($1, $2)
RETURNING id
`

type ScheduleSerializedShutterMessageParams struct {
	Description string
	Msg         []byte
}

func (q *Queries) ScheduleSerializedShutterMessage(ctx context.Context, arg ScheduleSerializedShutterMessageParams) (int32, error) {
	row := q.db.QueryRow(ctx, scheduleSerializedShutterMessage, arg.Description, arg.Msg)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const selectDecryptionKeyShares = `-- name: SelectDecryptionKeyShares :many
SELECT eon, epoch_id, keyper_index, decryption_key_share FROM decryption_key_share
WHERE eon = $1 AND epoch_id = $2
`

type SelectDecryptionKeySharesParams struct {
	Eon     int64
	EpochID []byte
}

func (q *Queries) SelectDecryptionKeyShares(ctx context.Context, arg SelectDecryptionKeySharesParams) ([]DecryptionKeyShare, error) {
	rows, err := q.db.Query(ctx, selectDecryptionKeyShares, arg.Eon, arg.EpochID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DecryptionKeyShare
	for rows.Next() {
		var i DecryptionKeyShare
		if err := rows.Scan(
			&i.Eon,
			&i.EpochID,
			&i.KeyperIndex,
			&i.DecryptionKeyShare,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const selectPureDKG = `-- name: SelectPureDKG :many
SELECT eon, puredkg FROM puredkg
`

func (q *Queries) SelectPureDKG(ctx context.Context) ([]Puredkg, error) {
	rows, err := q.db.Query(ctx, selectPureDKG)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Puredkg
	for rows.Next() {
		var i Puredkg
		if err := rows.Scan(&i.Eon, &i.Puredkg); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const setBatchConfigStarted = `-- name: SetBatchConfigStarted :exec
UPDATE tendermint_batch_config SET started = TRUE
WHERE keyper_config_index = $1
`

func (q *Queries) SetBatchConfigStarted(ctx context.Context, keyperConfigIndex int32) error {
	_, err := q.db.Exec(ctx, setBatchConfigStarted, keyperConfigIndex)
	return err
}

const setLastBatchConfigSent = `-- name: SetLastBatchConfigSent :exec
INSERT INTO last_batch_config_sent (keyper_config_index) VALUES ($1)
ON CONFLICT (enforce_one_row) DO UPDATE
SET keyper_config_index = $1
`

func (q *Queries) SetLastBatchConfigSent(ctx context.Context, keyperConfigIndex int64) error {
	_, err := q.db.Exec(ctx, setLastBatchConfigSent, keyperConfigIndex)
	return err
}

const setLastBlockSeen = `-- name: SetLastBlockSeen :exec
INSERT INTO last_block_seen (block_number) VALUES ($1)
ON CONFLICT (enforce_one_row) DO UPDATE
SET block_number = $1
`

func (q *Queries) SetLastBlockSeen(ctx context.Context, blockNumber int64) error {
	_, err := q.db.Exec(ctx, setLastBlockSeen, blockNumber)
	return err
}

const tMGetSyncMeta = `-- name: TMGetSyncMeta :one
SELECT current_block, last_committed_height, sync_timestamp
FROM tendermint_sync_meta
ORDER BY current_block DESC, last_committed_height DESC
LIMIT 1
`

func (q *Queries) TMGetSyncMeta(ctx context.Context) (TendermintSyncMetum, error) {
	row := q.db.QueryRow(ctx, tMGetSyncMeta)
	var i TendermintSyncMetum
	err := row.Scan(&i.CurrentBlock, &i.LastCommittedHeight, &i.SyncTimestamp)
	return i, err
}

const tMSetSyncMeta = `-- name: TMSetSyncMeta :exec
INSERT INTO tendermint_sync_meta (current_block, last_committed_height, sync_timestamp)
VALUES ($1, $2, $3)
`

type TMSetSyncMetaParams struct {
	CurrentBlock        int64
	LastCommittedHeight int64
	SyncTimestamp       time.Time
}

func (q *Queries) TMSetSyncMeta(ctx context.Context, arg TMSetSyncMetaParams) error {
	_, err := q.db.Exec(ctx, tMSetSyncMeta, arg.CurrentBlock, arg.LastCommittedHeight, arg.SyncTimestamp)
	return err
}