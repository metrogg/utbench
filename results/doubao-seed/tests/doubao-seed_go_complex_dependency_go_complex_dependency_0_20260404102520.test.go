package snap

import (
	"hash/crc32"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestRead(t *testing.T) {
	t.Parallel()
	testLogger := zaptest.NewLogger(t)
	tempDir := t.TempDir()

	t.Run("success with valid snapshot", func(t *testing.T) {
		t.Parallel()
		expectedSnap := &raftpb.Snapshot{
			Metadata: raftpb.SnapshotMetadata{
				Index: 1024,
				Term:  12,
			},
			Data: []byte("valid_test_snapshot_data_123"),
		}

		raftData, err := expectedSnap.Marshal()
		require.NoError(t, err)

		crc := crc32.Update(0, crcTable, raftData)
		serializedSnap := &snappb.Snapshot{
			Data: raftData,
			Crc:  crc,
		}

		snapBytes, err := serializedSnap.Marshal()
		require.NoError(t, err)

		tempFile, err := os.CreateTemp(tempDir, "valid_snap_*.bin")
		require.NoError(t, err)
		snapPath := tempFile.Name()
		_, err = tempFile.Write(snapBytes)
		require.NoError(t, err)
		require.NoError(t, tempFile.Close())

		result, err := Read(testLogger, snapPath)
		assert.NoError(t, err)
		assert.Equal(t, expectedSnap.Metadata.Index, result.Metadata.Index)
		assert.Equal(t, expectedSnap.Metadata.Term, result.Metadata.Term)
		assert.Equal(t, expectedSnap.Data, result.Data)
	})

	t.Run("file not found error", func(t *testing.T) {
		t.Parallel()
		nonExistentPath := tempDir + "/missing_snap_1234.bin"
		result, err := Read(nil, nonExistentPath)
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("empty file returns ErrEmptySnapshot", func(t *testing.T) {
		t.Parallel()
		tempFile, err := os.CreateTemp(tempDir, "empty_snap_*.bin")
		require.NoError(t, err)
		snapPath := tempFile.Name()
		require.NoError(t, tempFile.Close())

		result, err := Read(testLogger, snapPath)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, ErrEmptySnapshot)
	})

	t.Run("invalid snappb unmarshal returns error", func(t *testing.T) {
		t.Parallel()
		tempFile, err := os.CreateTemp(tempDir, "corrupt_snappb_*.bin")
		require.NoError(t, err)
		snapPath := tempFile.Name()
		_, err = tempFile.Write([]byte("invalid protobuf garbage data"))
		require.NoError(t, err)
		require.NoError(t, tempFile.Close())

		result, err := Read(nil, snapPath)
		assert.Nil(t, result)
		assert.Error(t, err)
	})

	t.Run("empty snappb data returns ErrEmptySnapshot", func(t *testing.T) {
		t.Parallel()
		serializedSnap := &snappb.Snapshot{
			Data: []byte{},
			Crc:  98765,
		}
		snapBytes, err := serializedSnap.Marshal()
		require.NoError(t, err)

		tempFile, err := os.CreateTemp(tempDir, "empty_data_snap_*.bin")
		require.NoError(t, err)
		snapPath := tempFile.Name()
		_, err = tempFile.Write(snapBytes)
		require.NoError(t, err)
		require.NoError(t, tempFile.Close())

		result, err := Read(testLogger, snapPath)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, ErrEmptySnapshot)
	})

	t.Run("zero CRC returns ErrEmptySnapshot", func(t *testing.T) {
		t.Parallel()
		serializedSnap := &snappb.Snapshot{
			Data: []byte("non-empty data with zero crc"),
			Crc:  0,
		}
		snapBytes, err := serializedSnap.Marshal()
		require.NoError(t, err)

		tempFile, err := os.CreateTemp(tempDir, "zero_crc_snap_*.bin")
		require.NoError(t, err)
		snapPath := tempFile.Name()
		_, err = tempFile.Write(snapBytes)
		require.NoError(t, err)
		require.NoError(t, tempFile.Close())

		result, err := Read(nil, snapPath)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, ErrEmptySnapshot)
	})

	t.Run("CRC mismatch returns ErrCRCMismatch", func(t *testing.T) {
		t.Parallel()
		testData := []byte("test data for crc check")
		wrongCrc := uint32(123123)
		serializedSnap := &snappb.Snapshot{
			Data: testData,
			Crc:  wrongCrc,
		}
		snapBytes, err := serializedSnap.Marshal()
		require.NoError(t, err)

		tempFile, err := os.CreateTemp(tempDir, "crc_mismatch_snap_*.bin")
		require.NoError(t, err)
		snapPath := tempFile.Name()
		_, err = tempFile.Write(snapBytes)
		require.NoError(t, err)
		require.NoError(t, tempFile.Close())

		result, err := Read(testLogger, snapPath)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, ErrCRCMismatch)
	})

	t.Run("invalid raftpb unmarshal returns error", func(t *testing.T) {
		t.Parallel()
		invalidRaftData := []byte("invalid raft protobuf content")
		crc := crc32.Update(0, crcTable, invalidRaftData)
		serializedSnap := &snappb.Snapshot{
			Data: invalidRaftData,
			Crc:  crc,
		}
		snapBytes, err := serializedSnap.Marshal()
		require.NoError(t, err)

		tempFile, err := os.CreateTemp(tempDir, "invalid_raft_snap_*.bin")
		require.NoError(t, err)
		snapPath := tempFile.Name()
		_, err = tempFile.Write(snapBytes)
		require.NoError(t, err)
		require.NoError(t, tempFile.Close())

		result, err := Read(nil, snapPath)
		assert.Nil(t, result)
		assert.Error(t, err)
	})
}