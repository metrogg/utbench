package snap

import (
	"bytes"
	"errors"
	"hash/crc32"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/gogo/protobuf/proto"
	"go.etcd.io/etcd/raft/raftpb"
	"go.etcd.io/etcd/snap/snappb"
	"go.uber.org/zap"
)

func TestRead(t *testing.T) {
	// Create a logger for testing
	logger := zap.NewNop()

	// Helper to create a valid snapshot file
	createValidSnapshotFile := func(t *testing.T, path string) {
		raftSnap := &raftpb.Snapshot{
			Data: []byte("test-data"),
		}
		raftData, err := proto.Marshal(raftSnap)
		if err != nil {
			t.Fatalf("failed to marshal raft snapshot: %v", err)
		}

		crc := crc32.Update(0, crcTable, raftData)
		serializedSnap := &snappb.Snapshot{
			Data: raftData,
			Crc:  crc,
		}
		snapData, err := proto.Marshal(serializedSnap)
		if err != nil {
			t.Fatalf("failed to marshal snappb snapshot: %v", err)
		}

		if err := ioutil.WriteFile(path, snapData, 0644); err != nil {
			t.Fatalf("failed to write snapshot file: %v", err)
		}
	}

	t.Run("Successfully read valid snapshot", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "snap-valid")
		createValidSnapshotFile(t, path)

		snap, err := Read(logger, path)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if snap == nil {
			t.Fatal("expected snapshot, got nil")
		}
		if !bytes.Equal(snap.Data, []byte("test-data")) {
			t.Errorf("unexpected data in snapshot, got %v", snap.Data)
		}
	})

	t.Run("Successfully read with nil logger (fallback to plog)", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "snap-nil-log")
		createValidSnapshotFile(t, path)

		snap, err := Read(nil, path)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if snap == nil {
			t.Fatal("expected snapshot, got nil")
		}
	})

	t.Run("Error when file does not exist", func(t *testing.T) {
		path := "/non/existent/path.snap"

		snap, err := Read(logger, path)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if snap != nil {
			t.Error("expected nil snapshot, got non-nil")
		}
		if !os.IsNotExist(err) {
			t.Errorf("expected file not found error, got %v", err)
		}
	})

	t.Run("Error when file is empty", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "snap-empty")
		if err := ioutil.WriteFile(path, []byte{}, 0644); err != nil {
			t.Fatal(err)
		}

		snap, err := Read(logger, path)
		if !errors.Is(err, ErrEmptySnapshot) {
			t.Errorf("expected ErrEmptySnapshot, got %v", err)
		}
		if snap != nil {
			t.Error("expected nil snapshot, got non-nil")
		}
	})

	t.Run("Error when file contains invalid snappb data", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "snap-invalid-snappb")
		if err := ioutil.WriteFile(path, []byte("invalid-data"), 0644); err != nil {
			t.Fatal(err)
		}

		snap, err := Read(logger, path)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if snap != nil {
			t.Error("expected nil snapshot, got non-nil")
		}
	})

	t.Run("Error when snappb data is empty (Data field)", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "snap-empty-data")

		// Create a snappb snapshot with empty data
		serializedSnap := &snappb.Snapshot{
			Data: []byte{},
			Crc:  0,
		}
		snapData, err := proto.Marshal(serializedSnap)
		if err != nil {
			t.Fatal(err)
		}
		if err := ioutil.WriteFile(path, snapData, 0644); err != nil {
			t.Fatal(err)
		}

		snap, err := Read(logger, path)
		if !errors.Is(err, ErrEmptySnapshot) {
			t.Errorf("expected ErrEmptySnapshot, got %v", err)
		}
		if snap != nil {
			t.Error("expected nil snapshot, got non-nil")
		}
	})

	t.Run("Error when CRC mismatch", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "snap-crc-mismatch")

		raftSnap := &raftpb.Snapshot{Data: []byte("data")}
		raftData, _ := proto.Marshal(raftSnap)

		// Set an incorrect CRC
		serializedSnap := &snappb.Snapshot{
			Data: raftData,
			Crc:  12345,
		}
		snapData, _ := proto.Marshal(serializedSnap)
		if err := ioutil.WriteFile(path, snapData, 0644); err != nil {
			t.Fatal(err)
		}

		snap, err := Read(logger, path)
		if !errors.Is(err, ErrCRCMismatch) {
			t.Errorf("expected ErrCRCMismatch, got %v", err)
		}
		if snap != nil {
			t.Error("expected nil snapshot, got non-nil")
		}
	})

	t.Run("Error when raftpb data is invalid", func(t *testing.T) {
		// This tests the unmarshalling of the inner raftpb.Snapshot
		dir := t.TempDir()
		path := filepath.Join(dir, "snap-invalid-raftpb")

		// Valid snappb wrapper, but invalid raftpb data inside
		invalidRaftData := []byte("not-a-raft-snapshot")
		crc := crc32.Update(0, crcTable, invalidRaftData)

		serializedSnap := &snappb.Snapshot{
			Data: invalidRaftData,
			Crc:  crc,
		}
		snapData, _ := proto.Marshal(serializedSnap)
		if err := ioutil.WriteFile(path, snapData, 0644); err != nil {
			t.Fatal(err)
		}

		snap, err := Read(logger, path)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if snap != nil {
			t.Error("expected nil snapshot, got non-nil")
		}
	})
}