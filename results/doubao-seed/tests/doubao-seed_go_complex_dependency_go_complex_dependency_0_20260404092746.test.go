package snap

import (
	"errors"
	"hash/crc32"
	"os"
	"testing"

	"go.uber.org/zap"

	raftpb "your/project/path/raftpb"
	snappb "your/project/path/snappb"
)

var crcTable = crc32.MakeTable(crc32.Castagnoli)

func TestRead(t *testing.T) {
	t.Parallel()
	lg := zap.NewNop()

	t.Run("file_not_exist", func(t *testing.T) {
		t.Parallel()
		nonExistPath := "non_exist_snap_1234.tmp"
		snap, err := Read(lg, nonExistPath)
		if snap != nil {
			t.Error("expected nil snapshot for non-existent file")
		}
		if !errors.Is(err, os.ErrNotExist) {
			t.Errorf("expected os.ErrNotExist, got %v", err)
		}

		snap, err = Read(nil, nonExistPath)
		if snap != nil {
			t.Error("expected nil snapshot for non-existent file with nil logger")
		}
		if !errors.Is(err, os.ErrNotExist) {
			t.Errorf("expected os.ErrNotExist with nil logger, got %v", err)
		}
	})

	t.Run("empty_file", func(t *testing.T) {
		t.Parallel()
		f, err := os.CreateTemp("", "empty_snap_*.tmp")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		defer os.Remove(f.Name())
		_ = f.Close()

		snap, err := Read(lg, f.Name())
		if snap != nil {
			t.Error("expected nil snapshot for empty file")
		}
		if !errors.Is(err, ErrEmptySnapshot) {
			t.Errorf("expected ErrEmptySnapshot, got %v", err)
		}

		snap, err = Read(nil, f.Name())
		if snap != nil {
			t.Error("expected nil snapshot for empty file with nil logger")
		}
		if !errors.Is(err, ErrEmptySnapshot) {
			t.Errorf("expected ErrEmptySnapshot with nil logger, got %v", err)
		}
	})

	t.Run("invalid_snappb_snapshot", func(t *testing.T) {
		t.Parallel()
		f, err := os.CreateTemp("", "invalid_snappb_*.tmp")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		defer os.Remove(f.Name())
		_, _ = f.Write([]byte("garbage non-proto content"))
		_ = f.Close()

		snap, err := Read(lg, f.Name())
		if snap != nil {
			t.Error("expected nil snapshot for invalid snappb")
		}
		if err == nil {
			t.Error("expected error for invalid snappb, got nil")
		}
	})

	t.Run("snappb_empty_data", func(t *testing.T) {
		t.Parallel()
		serialSnap := snappb.Snapshot{
			Data: []byte{},
			Crc:  12345,
		}
		b, err := serialSnap.Marshal()
		if err != nil {
			t.Fatalf("failed to marshal empty data snappb: %v", err)
		}

		f, err := os.CreateTemp("", "empty_data_snap_*.tmp")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		defer os.Remove(f.Name())
		_, _ = f.Write(b)
		_ = f.Close()

		snap, err := Read(lg, f.Name())
		if snap != nil {
			t.Error("expected nil snapshot for empty data")
		}
		if !errors.Is(err, ErrEmptySnapshot) {
			t.Errorf("expected ErrEmptySnapshot, got %v", err)
		}
	})

	t.Run("snappb_zero_crc", func(t *testing.T) {
		t.Parallel()
		serialSnap := snappb.Snapshot{
			Data: []byte("non-empty test data"),
			Crc:  0,
		}
		b, err := serialSnap.Marshal()
		if err != nil {
			t.Fatalf("failed to marshal zero crc snappb: %v", err)
		}

		f, err := os.CreateTemp("", "zero_crc_snap_*.tmp")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		defer os.Remove(f.Name())
		_, _ = f.Write(b)
		_ = f.Close()

		snap, err := Read(lg, f.Name())
		if snap != nil {
			t.Error("expected nil snapshot for zero crc")
		}
		if !errors.Is(err, ErrEmptySnapshot) {
			t.Errorf("expected ErrEmptySnapshot, got %v", err)
		}
	})

	t.Run("crc_mismatch", func(t *testing.T) {
		t.Parallel()
		data := []byte("test data for crc check")
		correctCrc := crc32.Update(0, crcTable, data)
		wrongCrc := correctCrc + 1

		serialSnap := snappb.Snapshot{
			Data: data,
			Crc:  wrongCrc,
		}
		b, err := serialSnap.Marshal()
		if err != nil {
			t.Fatalf("failed to marshal crc mismatch snappb: %v", err)
		}

		f, err := os.CreateTemp("", "crc_mismatch_snap_*.tmp")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		defer os.Remove(f.Name())
		_, _ = f.Write(b)
		_ = f.Close()

		snap, err := Read(lg, f.Name())
		if snap != nil {
			t.Error("expected nil snapshot for crc mismatch")
		}
		if !errors.Is(err, ErrCRCMismatch) {
			t.Errorf("expected ErrCRCMismatch, got %v", err)
		}
	})

	t.Run("invalid_raftpb_snapshot", func(t *testing.T) {
		t.Parallel()
		invalidRaftData := []byte("invalid raft proto content")
		crc := crc32.Update(0, crcTable, invalidRaftData)

		serialSnap := snappb.Snapshot{
			Data: invalidRaftData,
			Crc:  crc,
		}
		b, err := serialSnap.Marshal()
		if err != nil {
			t.Fatalf("failed to marshal invalid raft snappb: %v", err)
		}

		f, err := os.CreateTemp("", "invalid_raft_snap_*.tmp")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		defer os.Remove(f.Name())
		_, _ = f.Write(b)
		_ = f.Close()

		snap, err := Read(lg, f.Name())
		if snap != nil {
			t.Error("expected nil snapshot for invalid raftpb")
		}
		if err == nil {
			t.Error("expected error for invalid raftpb, got nil")
		}
	})

	t.Run("success_valid_snapshot", func(t *testing.T) {
		t.Parallel()
		expectedSnap := &raftpb.Snapshot{
			Term:  1024,
			Index: 2048,
		}
		raftData, err := expectedSnap.Marshal()
		if err != nil {
			t.Fatalf("failed to marshal valid raftpb: %v", err)
		}

		crc := crc32.Update(0, crcTable, raftData)
		serialSnap := snappb.Snapshot{
			Data: raftData,
			Crc:  crc,
		}
		snapBytes, err := serialSnap.Marshal()
		if err != nil {
			t.Fatalf("failed to marshal valid snappb: %v", err)
		}

		f, err := os.CreateTemp("", "valid_snap_*.tmp")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		defer os.Remove(f.Name())
		_, _ = f.Write(snapBytes)
		_ = f.Close()

		snap, err := Read(lg, f.Name())
		if err != nil {
			t.Fatalf("unexpected error reading valid snapshot: %v", err)
		}
		if snap == nil {
			t.Fatal("expected non-nil snapshot for valid input")
		}
		if snap.Term != expectedSnap.Term {
			t.Errorf("expected term %d, got %d", expectedSnap.Term, snap.Term)
		}
		if snap.Index != expectedSnap.Index {
			t.Errorf("expected index %d, got %d", expectedSnap.Index, snap.Index)
		}

		snap, err = Read(nil, f.Name())
		if err != nil {
			t.Fatalf("unexpected error reading valid snapshot with nil logger: %v", err)
		}
		if snap == nil {
			t.Fatal("expected non-nil snapshot for valid input with nil logger")
		}
		if snap.Term != expectedSnap.Term {
			t.Errorf("expected term %d with nil logger, got %d", expectedSnap.Term, snap.Term)
		}
		if snap.Index != expectedSnap.Index {
			t.Errorf("expected index %d with nil logger, got %d", expectedSnap.Index, snap.Index)
		}
	})
}