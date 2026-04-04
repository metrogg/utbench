package main

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"go.etcd.io/etcd/raft/raftpb"
	"go.etcd.io/etcd/snap/snappb"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// 测试辅助函数：创建临时文件
func createTempFile(t *testing.T, content []byte) string {
	t.Helper()
	tmpfile, err := ioutil.TempFile("", "snapshot_test")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Remove(tmpfile.Name()) })

	if content != nil {
		if _, err := tmpfile.Write(content); err != nil {
			t.Fatal(err)
		}
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}
	return tmpfile.Name()
}

// 测试辅助函数：创建有效的序列化快照数据
func createValidSnapshotData(t *testing.T) []byte {
	t.Helper()
	snap := &raftpb.Snapshot{
		Data: []byte("test snapshot data"),
	}
	data, err := snap.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	
	serializedSnap := &snappb.Snapshot{
		Data: data,
		Crc:  crc32.Update(0, crcTable, data),
	}
	
	serializedData, err := serializedSnap.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	return serializedData
}

func TestRead(t *testing.T) {
	// 设置测试用的crcTable（假设这是包级变量）
	// 注意：实际测试中需要确保crcTable已正确初始化
	// 这里我们假设它已经在被测代码中定义

	t.Run("正常路径-使用zap日志器", func(t *testing.T) {
		// 创建观察者日志器
		core, recorded := observer.New(zapcore.InfoLevel)
		logger := zap.New(core)
		
		// 创建有效的快照文件
		validData := createValidSnapshotData(t)
		snapPath := createTempFile(t, validData)
		
		// 执行读取
		snap, err := Read(logger, snapPath)
		
		// 验证结果
		if err != nil {
			t.Errorf("期望无错误，但得到: %v", err)
		}
		if snap == nil {
			t.Error("期望返回非nil快照")
		}
		if string(snap.Data) != "test snapshot data" {
			t.Errorf("期望快照数据为'test snapshot data'，但得到: %s", snap.Data)
		}
		
		// 验证没有警告日志
		if logs := recorded.FilterMessage("failed to read a snap file").Len(); logs > 0 {
			t.Error("不期望看到失败读取日志")
		}
	})

	t.Run("正常路径-无日志器", func(t *testing.T) {
		// 创建有效的快照文件
		validData := createValidSnapshotData(t)
		snapPath := createTempFile(t, validData)
		
		// 执行读取（无日志器）
		snap, err := Read(nil, snapPath)
		
		// 验证结果
		if err != nil {
			t.Errorf("期望无错误，但得到: %v", err)
		}
		if snap == nil {
			t.Error("期望返回非nil快照")
		}
	})

	t.Run("文件不存在", func(t *testing.T) {
		core, recorded := observer.New(zapcore.WarnLevel)
		logger := zap.New(core)
		
		nonExistentPath := filepath.Join(t.TempDir(), "nonexistent.snap")
		snap, err := Read(logger, nonExistentPath)
		
		// 验证错误
		if !errors.Is(err, os.ErrNotExist) {
			t.Errorf("期望文件不存在错误，但得到: %v", err)
		}
		if snap != nil {
			t.Error("期望返回nil快照")
		}
		
		// 验证日志
		logs := recorded.All()
		if len(logs) != 1 {
			t.Errorf("期望1条警告日志，但得到%d条", len(logs))
		}
		if logs[0].Message != "failed to read a snap file" {
			t.Errorf("期望日志消息为'failed to read a snap file'，但得到: %s", logs[0].Message)
		}
	})

	t.Run("空文件", func(t *testing.T) {
		core, recorded := observer.New(zapcore.WarnLevel)
		logger := zap.New(core)
		
		emptyPath := createTempFile(t, []byte{})
		snap, err := Read(logger, emptyPath)
		
		// 验证错误
		if err != ErrEmptySnapshot {
			t.Errorf("期望ErrEmptySnapshot错误，但得到: %v", err)
		}
		if snap != nil {
			t.Error("期望返回nil快照")
		}
		
		// 验证日志
		logs := recorded.All()
		if len(logs) != 1 {
			t.Errorf("期望1条警告日志，但得到%d条", len(logs))
		}
		if logs[0].Message != "failed to read empty snapshot file" {
			t.Errorf("期望日志消息为'failed to read empty snapshot file'，但得到: %s", logs[0].Message)
		}
	})

	t.Run("无效的序列化数据", func(t *testing.T) {
		core, recorded := observer.New(zapcore.WarnLevel)
		logger := zap.New(core)
		
		// 创建无效的序列化数据
		invalidData := []byte("invalid protobuf data")
		invalidPath := createTempFile(t, invalidData)
		
		snap, err := Read(logger, invalidPath)
		
		// 验证错误
		if err == nil {
			t.Error("期望反序列化错误，但得到nil")
		}
		if snap != nil {
			t.Error("期望返回nil快照")
		}
		
		// 验证日志
		logs := recorded.All()
		if len(logs) != 1 {
			t.Errorf("期望1条警告日志，但得到%d条", len(logs))
		}
		if logs[0].Message != "failed to unmarshal snappb.Snapshot" {
			t.Errorf("期望日志消息为'failed to unmarshal snappb.Snapshot'，但得到: %s", logs[0].Message)
		}
	})

	t.Run("空的快照数据", func(t *testing.T) {
		core, recorded := observer.New(zapcore.WarnLevel)
		logger := zap.New(core)
		
		// 创建空的快照数据
		serializedSnap := &snappb.Snapshot{
			Data: []byte{},
			Crc:  123,
		}
		data, err := serializedSnap.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		
		emptyDataPath := createTempFile(t, data)
		snap, err := Read(logger, emptyDataPath)
		
		// 验证错误
		if err != ErrEmptySnapshot {
			t.Errorf("期望ErrEmptySnapshot错误，但得到: %v", err)
		}
		if snap != nil {
			t.Error("期望返回nil快照")
		}
		
		// 验证日志
		logs := recorded.All()
		if len(logs) != 1 {
			t.Errorf("期望1条警告日志，但得到%d条", len(logs))
		}
		if logs[0].Message != "failed to read empty snapshot data" {
			t.Errorf("期望日志消息为'failed to read empty snapshot data'，但得到: %s", logs[0].Message)
		}
	})

	t.Run("CRC校验失败", func(t *testing.T) {
		core, recorded := observer.New(zapcore.WarnLevel)
		logger := zap.New(core)
		
		// 创建有效的快照数据但错误的CRC
		snap := &raftpb.Snapshot{
			Data: []byte("test data"),
		}
		data, err := snap.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		
		serializedSnap := &snappb.Snapshot{
			Data: data,
			Crc:  999999, // 错误的CRC值
		}
		
		serializedData, err := serializedSnap.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		
		crcMismatchPath := createTempFile(t, serializedData)
		snapResult, err := Read(logger, crcMismatchPath)
		
		// 验证错误
		if err != ErrCRCMismatch {
			t.Errorf("期望ErrCRCMismatch错误，但得到: %v", err)
		}
		if snapResult != nil {
			t.Error("期望返回nil快照")
		}
		
		// 验证日志
		logs := recorded.All()
		if len(logs) != 1 {
			t.Errorf("期望1条警告日志，但得到%d条", len(logs))
		}
		if logs[0].Message != "snap file is corrupt" {
			t.Errorf("期望日志消息为'snap file is corrupt'，但得到: %s", logs[0].Message)
		}
	})

	t.Run("无效的raftpb.Snapshot数据", func(t *testing.T) {
		core, recorded := observer.New(zapcore.WarnLevel)
		logger := zap.New(core)
		
		// 创建有效的序列化快照但无效的内部数据
		invalidInnerData := []byte("invalid inner data")
		serializedSnap := &snappb.Snapshot{
			Data: invalidInnerData,
			Crc:  crc32.Update(0, crcTable, invalidInnerData),
		}
		
		serializedData, err := serializedSnap.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		
		invalidInnerPath := createTempFile(t, serializedData)
		snap, err := Read(logger, invalidInnerPath)
		
		// 验证错误
		if err == nil {
			t.Error("期望反序列化错误，但得到nil")
		}
		if snap != nil {
			t.Error("期望返回nil快照")
		}
		
		// 验证日志
		logs := recorded.All()
		if len(logs) != 1 {
			t.Errorf("期望1条警告日志，但得到%d条", len(logs))
		}
		if logs[0].Message != "failed to unmarshal raftpb.Snapshot" {
			t.Errorf("期望日志消息为'failed to unmarshal raftpb.Snapshot'，但得到: %s", logs[0].Message)
		}
	})

	t.Run("边界条件-零CRC值", func(t *testing.T) {
		core, recorded := observer.New(zapcore.WarnLevel)
		logger := zap.New(core)
		
		// 创建零CRC的快照数据
		snap := &raftpb.Snapshot{
			Data: []byte("test data"),
		}
		data, err := snap.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		
		serializedSnap := &snappb.Snapshot{
			Data: data,
			Crc:  0, // 零CRC值
		}
		
		serializedData, err := serializedSnap.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		
		zeroCrcPath := createTempFile(t, serializedData)
		snapResult, err := Read(logger, zeroCrcPath)
		
		// 验证错误（根据代码逻辑，零CRC应该返回ErrEmptySnapshot）
		if err != ErrEmptySnapshot {
			t.Errorf("期望ErrEmptySnapshot错误，但得到: %v", err)
		}
		if snapResult != nil {
			t.Error("期望返回nil快照")
		}
		
		// 验证日志
		logs := recorded.All()
		if len(logs) != 1 {
			t.Errorf("期望1条警告日志，但得到%d条", len(logs))
		}
		if logs[0].Message != "failed to read empty snapshot data" {
			t.Errorf("期望日志消息为'failed to read empty snapshot data'，但得到: %s", logs[0].Message)
		}
	})

	t.Run("文件读取权限错误", func(t *testing.T) {
		core, recorded := observer.New(zapcore.WarnLevel)
		logger := zap.New(core)
		
		// 创建临时目录
		tmpDir := t.TempDir()
		readOnlyPath := filepath.Join(tmpDir, "readonly.snap")
		
		// 创建文件但不给读取权限
		if err := ioutil.WriteFile(readOnlyPath, []byte("test"), 0222); err != nil {
			t.Fatal(err)
		}
		
		snap, err := Read(logger, readOnlyPath)
		
		// 验证错误
		if err == nil {
			t.Error("期望权限错误，但得到nil")
		}
		if snap != nil {
			t.Error("期望返回nil快照")
		}
		
		// 验证日志
		logs := recorded.All()
		if len(logs) != 1 {
			t.Errorf("期望1条警告日志，但得到%d条", len(logs))
		}
		if logs[0].Message != "failed to read a snap file" {
			t.Errorf("期望日志消息为'failed to read a snap file'，但得到: %s", logs[0].Message)
		}
	})
}