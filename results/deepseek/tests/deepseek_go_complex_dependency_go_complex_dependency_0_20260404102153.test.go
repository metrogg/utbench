package main

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"go.etcd.io/etcd/raft/v3/raftpb"
	"go.etcd.io/etcd/server/v3/etcdserver/api/snap"
	snappb "go.etcd.io/etcd/server/v3/etcdserver/api/snap"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

// Mock ioutil.ReadFile for testing
type mockIoutil struct {
	readFileFunc func(string) ([]byte, error)
}

func (m *mockIoutil) ReadFile(filename string) ([]byte, error) {
	if m.readFileFunc != nil {
		return m.readFileFunc(filename)
	}
	return ioutil.ReadFile(filename)
}

func TestRead(t *testing.T) {
	// Setup test data
	validSnapData := []byte("valid snapshot data")
	validCRC := crc32.Update(0, crcTable, validSnapData)
	
	validSerializedSnap := snappb.Snapshot{
		Data: validSnapData,
		Crc:  validCRC,
	}
	validSerializedBytes, _ := validSerializedSnap.Marshal()
	
	validRaftSnap := raftpb.Snapshot{
		Data: []byte("raft snapshot data"),
	}
	validRaftSnapData, _ := validRaftSnap.Marshal()
	validSerializedSnapWithRaft := snappb.Snapshot{
		Data: validRaftSnapData,
		Crc:  crc32.Update(0, crcTable, validRaftSnapData),
	}
	validSerializedBytesWithRaft, _ := validSerializedSnapWithRaft.Marshal()

	tests := []struct {
		name           string
		snapname       string
		fileContent    []byte
		readFileErr    error
		lg             *zap.Logger
		wantErr        error
		wantSnap       *raftpb.Snapshot
		setupMock      func(*mockIoutil)
	}{
		// Normal path tests
		{
			name:        "success with logger",
			snapname:    "test.snap",
			fileContent: validSerializedBytesWithRaft,
			lg:          zaptest.NewLogger(t),
			wantSnap:    &validRaftSnap,
			wantErr:     nil,
		},
		{
			name:        "success without logger",
			snapname:    "test.snap",
			fileContent: validSerializedBytesWithRaft,
			lg:          nil,
			wantSnap:    &validRaftSnap,
			wantErr:     nil,
		},
		
		// Error path tests - file reading errors
		{
			name:        "read file error with logger",
			snapname:    "test.snap",
			readFileErr: errors.New("file not found"),
			lg:          zaptest.NewLogger(t),
			wantErr:     errors.New("file not found"),
		},
		{
			name:        "read file error without logger",
			snapname:    "test.snap",
			readFileErr: errors.New("file not found"),
			lg:          nil,
			wantErr:     errors.New("file not found"),
		},
		
		// Boundary condition tests - empty file
		{
			name:        "empty file with logger",
			snapname:    "test.snap",
			fileContent: []byte{},
			lg:          zaptest.NewLogger(t),
			wantErr:     snap.ErrEmptySnapshot,
		},
		{
			name:        "empty file without logger",
			snapname:    "test.snap",
			fileContent: []byte{},
			lg:          nil,
			wantErr:     snap.ErrEmptySnapshot,
		},
		
		// Error path tests - unmarshal errors
		{
			name:        "invalid serialized snap with logger",
			snapname:    "test.snap",
			fileContent: []byte("invalid data"),
			lg:          zaptest.NewLogger(t),
			wantErr:     errors.New("proto:"), // proto unmarshal error
		},
		{
			name:        "invalid serialized snap without logger",
			snapname:    "test.snap",
			fileContent: []byte("invalid data"),
			lg:          nil,
			wantErr:     errors.New("proto:"), // proto unmarshal error
		},
		
		// Boundary condition tests - empty snapshot data
		{
			name:     "empty snapshot data with logger",
			snapname: "test.snap",
			fileContent: func() []byte {
				emptySnap := snappb.Snapshot{
					Data: []byte{},
					Crc:  0,
				}
				b, _ := emptySnap.Marshal()
				return b
			}(),
			lg:      zaptest.NewLogger(t),
			wantErr: snap.ErrEmptySnapshot,
		},
		{
			name:     "empty snapshot data without logger",
			snapname: "test.snap",
			fileContent: func() []byte {
				emptySnap := snappb.Snapshot{
					Data: []byte{},
					Crc:  0,
				}
				b, _ := emptySnap.Marshal()
				return b
			}(),
			lg:      nil,
			wantErr: snap.ErrEmptySnapshot,
		},
		{
			name:     "zero crc with logger",
			snapname: "test.snap",
			fileContent: func() []byte {
				emptySnap := snappb.Snapshot{
					Data: []byte("data"),
					Crc:  0,
				}
				b, _ := emptySnap.Marshal()
				return b
			}(),
			lg:      zaptest.NewLogger(t),
			wantErr: snap.ErrEmptySnapshot,
		},
		
		// Error path tests - CRC mismatch
		{
			name:     "crc mismatch with logger",
			snapname: "test.snap",
			fileContent: func() []byte {
				corruptSnap := snappb.Snapshot{
					Data: validSnapData,
					Crc:  validCRC + 1, // Wrong CRC
				}
				b, _ := corruptSnap.Marshal()
				return b
			}(),
			lg:      zaptest.NewLogger(t),
			wantErr: snap.ErrCRCMismatch,
		},
		{
			name:     "crc mismatch without logger",
			snapname: "test.snap",
			fileContent: func() []byte {
				corruptSnap := snappb.Snapshot{
					Data: validSnapData,
					Crc:  validCRC + 1, // Wrong CRC
				}
				b, _ := corruptSnap.Marshal()
				return b
			}(),
			lg:      nil,
			wantErr: snap.ErrCRCMismatch,
		},
		
		// Error path tests - raftpb unmarshal error
		{
			name:     "invalid raftpb data with logger",
			snapname: "test.snap",
			fileContent: func() []byte {
				invalidSnap := snappb.Snapshot{
					Data: []byte("invalid raftpb data"),
					Crc:  crc32.Update(0, crcTable, []byte("invalid raftpb data")),
				}
				b, _ := invalidSnap.Marshal()
				return b
			}(),
			lg:      zaptest.NewLogger(t),
			wantErr: errors.New("proto:"), // proto unmarshal error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp file for testing
			tmpDir := t.TempDir()
			snapPath := filepath.Join(tmpDir, tt.snapname)
			
			// Write file content if provided
			if tt.fileContent != nil {
				err := ioutil.WriteFile(snapPath, tt.fileContent, 0644)
				if err != nil {
					t.Fatalf("failed to write test file: %v", err)
				}
			}
			
			// Mock ioutil.ReadFile
			originalReadFile := ioutilReadFile
			defer func() { ioutilReadFile = originalReadFile }()
			
			ioutilReadFile = func(filename string) ([]byte, error) {
				if tt.readFileErr != nil {
					return nil, tt.readFileErr
				}
				if tt.fileContent != nil && filename == snapPath {
					return tt.fileContent, nil
				}
				return ioutil.ReadFile(filename)
			}
			
			// Call the function under test
			gotSnap, gotErr := Read(tt.lg, snapPath)
			
			// Check error
			if tt.wantErr != nil {
				if gotErr == nil {
					t.Errorf("Read() error = nil, want %v", tt.wantErr)
				} else if tt.wantErr != snap.ErrEmptySnapshot && tt.wantErr != snap.ErrCRCMismatch {
					// For non-sentinel errors, just check that an error occurred
					if gotErr == nil {
						t.Errorf("Read() error = nil, want error")
					}
				} else if !errors.Is(gotErr, tt.wantErr) {
					t.Errorf("Read() error = %v, want %v", gotErr, tt.wantErr)
				}
			} else {
				if gotErr != nil {
					t.Errorf("Read() unexpected error = %v", gotErr)
				}
			}
			
			// Check snapshot
			if tt.wantSnap != nil {
				if gotSnap == nil {
					t.Errorf("Read() got nil snapshot, want %v", tt.wantSnap)
				} else if string(gotSnap.Data) != string(tt.wantSnap.Data) {
					t.Errorf("Read() snapshot data = %v, want %v", gotSnap.Data, tt.wantSnap.Data)
				}
			} else if gotSnap != nil {
				t.Errorf("Read() got unexpected snapshot = %v", gotSnap)
			}
		})
	}
}

// Global variable to mock ioutil.ReadFile
var ioutilReadFile = ioutil.ReadFile

func init() {
	// Override ioutil.ReadFile in the snap package
	ioutil.ReadFile = func(filename string) ([]byte, error) {
		return ioutilReadFile(filename)
	}
}

// Helper to create test snapshot files
func createTestSnapshotFile(t *testing.T, data []byte, crc uint32) []byte {
	snap := snappb.Snapshot{
		Data: data,
		Crc:  crc,
	}
	b, err := snap.Marshal()
	if err != nil {
		t.Fatalf("failed to marshal test snapshot: %v", err)
	}
	return b
}

// Test edge cases
func TestRead_EdgeCases(t *testing.T) {
	t.Run("nil logger with all error paths", func(t *testing.T) {
		tmpDir := t.TempDir()
		snapPath := filepath.Join(tmpDir, "test.snap")
		
		// Test file read error
		originalReadFile := ioutilReadFile
		ioutilReadFile = func(filename string) ([]byte, error) {
			return nil, errors.New("permission denied")
		}
		defer func() { ioutilReadFile = originalReadFile }()
		
		_, err := Read(nil, snapPath)
		if err == nil {
			t.Error("expected error for permission denied")
		}
	})
	
	t.Run("valid logger with all success paths", func(t *testing.T) {
		tmpDir := t.TempDir()
		snapPath := filepath.Join(tmpDir, "test.snap")
		
		// Create valid snapshot
		raftSnap := raftpb.Snapshot{
			Data: []byte("test data"),
		}
		raftData, _ := raftSnap.Marshal()
		snap := snappb.Snapshot{
			Data: raftData,
			Crc:  crc32.Update(0, crcTable, raftData),
		}
		snapData, _ := snap.Marshal()
		
		err := ioutil.WriteFile(snapPath, snapData, 0644)
		if err != nil {
			t.Fatalf("failed to write test file: %v", err)
		}
		
		gotSnap, err := Read(zaptest.NewLogger(t), snapPath)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if gotSnap == nil {
			t.Error("expected non-nil snapshot")
		}
	})
}

// Test with mock controller for more complex scenarios
func TestRead_WithMockController(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	
	t.Run("file system error simulation", func(t *testing.T) {
		tmpDir := t.TempDir()
		snapPath := filepath.Join(tmpDir, "test.snap")
		
		// Simulate various file system errors
		testErrors := []error{
			&os.PathError{Op: "open", Path: snapPath, Err: os.ErrNotExist},
			&os.PathError{Op: "open", Path: snapPath, Err: os.ErrPermission},
			errors.New("disk full"),
		}
		
		for _, testErr := range testErrors {
			originalReadFile := ioutilReadFile
			ioutilReadFile = func(filename string) ([]byte, error) {
				return nil, testErr
			}
			
			_, err := Read(zaptest.NewLogger(t), snapPath)
			if err == nil {
				t.Errorf("expected error for %v", testErr)
			}
			
			ioutilReadFile = originalReadFile
		}
	})
}