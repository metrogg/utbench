import pytest
import random
import statistics
from unittest.mock import patch
from collections import OrderedDict

def task_func(LETTERS):
    random_dict = {k: [random.randint(0, 100) for _ in range(random.randint(1, 10))] for k in LETTERS}
    sorted_dict = dict(sorted(random_dict.items(), key=lambda item: statistics.mean(item[1]), reverse=True))
    return sorted_dict

class TestTaskFunc:
    def test_normal_case_with_multiple_letters(self):
        """测试正常情况：多个字母"""
        letters = ['A', 'B', 'C', 'D']
        
        with patch('random.randint') as mock_randint, \
             patch('random.randint') as mock_randint_len:
            
            # 设置随机数生成器返回固定值
            mock_randint.side_effect = [10, 20, 30, 40, 50, 60, 70, 80]
            mock_randint_len.side_effect = [3, 3, 3, 3]
            
            result = task_func(letters)
            
            # 验证结果包含所有字母
            assert set(result.keys()) == set(letters)
            # 验证排序正确（均值降序）
            means = [statistics.mean(v) for v in result.values()]
            assert means == sorted(means, reverse=True)
    
    def test_single_letter(self):
        """测试边界情况：单个字母"""
        letters = ['X']
        
        with patch('random.randint') as mock_randint, \
             patch('random.randint') as mock_randint_len:
            
            mock_randint.side_effect = [25, 50, 75]
            mock_randint_len.side_effect = [3]
            
            result = task_func(letters)
            
            assert len(result) == 1
            assert 'X' in result
            assert result['X'] == [25, 50, 75]
    
    def test_empty_input(self):
        """测试边界情况：空列表输入"""
        letters = []
        result = task_func(letters)
        
        assert result == {}
        assert len(result) == 0
    
    def test_duplicate_letters(self):
        """测试重复字母的情况"""
        letters = ['A', 'A', 'B']
        
        with patch('random.randint') as mock_randint, \
             patch('random.randint') as mock_randint_len:
            
            # 注意：字典键不能重复，所以第二个'A'会被覆盖
            mock_randint.side_effect = [10, 20, 30, 40, 50, 60]
            mock_randint_len.side_effect = [2, 2, 2]
            
            result = task_func(letters)
            
            # 由于字典键唯一，结果应该只有2个键
            assert len(result) == 2
            assert set(result.keys()) == {'A', 'B'}
    
    def test_sorting_by_mean_descending(self):
        """测试按均值降序排序的正确性"""
        letters = ['A', 'B', 'C']
        
        # 创建固定的测试数据
        test_data = {
            'A': [10, 20, 30],  # 均值: 20
            'B': [40, 50, 60],  # 均值: 50
            'C': [70, 80, 90],  # 均值: 80
        }
        
        def mock_randint(a, b):
            # 模拟返回固定值
            if hasattr(mock_randint, 'counter'):
                mock_randint.counter += 1
            else:
                mock_randint.counter = 0
            
            values = [10, 20, 30, 40, 50, 60, 70, 80, 90]
            return values[mock_randint.counter % len(values)]
        
        def mock_randint_len(a, b):
            return 3  # 固定列表长度为3
        
        with patch('random.randint') as mock_rand:
            mock_rand.side_effect = mock_randint
            
            with patch('random.randint') as mock_len:
                mock_len.side_effect = mock_randint_len
                
                result = task_func(letters)
                
                # 验证排序顺序：C(80) > B(50) > A(20)
                keys = list(result.keys())
                assert keys == ['C', 'B', 'A']
    
    def test_list_length_variation(self):
        """测试不同列表长度的情况"""
        letters = ['A', 'B']
        
        with patch('random.randint') as mock_randint, \
             patch('random.randint') as mock_randint_len:
            
            # A: 长度1，值[50]
            # B: 长度5，值[10, 20, 30, 40, 50]
            mock_randint.side_effect = [50, 10, 20, 30, 40, 50]
            mock_randint_len.side_effect = [1, 5]
            
            result = task_func(letters)
            
            # 验证列表长度
            assert len(result['A']) == 1
            assert len(result['B']) == 5
    
    def test_edge_case_all_same_mean(self):
        """测试边界情况：所有列表均值相同"""
        letters = ['A', 'B', 'C']
        
        with patch('random.randint') as mock_randint, \
             patch('random.randint') as mock_randint_len:
            
            # 所有列表都有相同的均值30
            mock_randint.side_effect = [20, 40, 25, 35, 30, 30]
            mock_randint_len.side_effect = [2, 2, 2]
            
            result = task_func(letters)
            
            # 当均值相同时，顺序可能保持原顺序（稳定排序）
            # 至少验证所有键都存在
            assert set(result.keys()) == set(letters)
    
    def test_large_input(self):
        """测试大量输入的情况"""
        letters = [chr(65 + i) for i in range(10)]  # A-J
        
        result = task_func(letters)
        
        # 验证所有字母都存在
        assert len(result) == 10
        for letter in letters:
            assert letter in result
        
        # 验证排序
        means = [statistics.mean(v) for v in result.values()]
        assert means == sorted(means, reverse=True)
    
    def test_randomness_control(self):
        """测试随机性控制"""
        letters = ['A', 'B']
        
        # 第一次调用
        with patch('random.randint') as mock_randint, \
             patch('random.randint') as mock_randint_len:
            
            mock_randint.side_effect = [10, 20, 30, 40]
            mock_randint_len.side_effect = [2, 2]
            
            result1 = task_func(letters)
        
        # 第二次调用（相同随机种子应产生相同结果）
        with patch('random.randint') as mock_randint, \
             patch('random.randint') as mock_randint_len:
            
            mock_randint.side_effect = [10, 20, 30, 40]
            mock_randint_len.side_effect = [2, 2]
            
            result2 = task_func(letters)
        
        # 验证两次结果相同
        assert result1 == result2
    
    def test_statistics_mean_usage(self):
        """测试statistics.mean的正确使用"""
        letters = ['A', 'B']
        
        with patch('random.randint') as mock_randint, \
             patch('random.randint') as mock_randint_len, \
             patch('statistics.mean') as mock_mean:
            
            mock_randint.side_effect = [10, 20, 30, 40]
            mock_randint_len.side_effect = [2, 2]
            
            # 模拟mean返回固定值
            mock_mean.side_effect = [15.0, 35.0]
            
            result = task_func(letters)
            
            # 验证mean被调用了正确次数
            assert mock_mean.call_count == 2
    
    def test_negative_case_empty_lists(self):
        """测试边界情况：空列表（理论上不会发生，因为random.randint(1,10)最小为1）"""
        letters = ['A']
        
        # 强制模拟长度为1的列表
        with patch('random.randint') as mock_randint, \
             patch('random.randint') as mock_randint_len:
            
            mock_randint.side_effect = [50]
            mock_randint_len.side_effect = [1]
            
            result = task_func(letters)
            
            # 验证单元素列表的均值计算
            assert statistics.mean(result['A']) == 50.0
    
    def test_return_type(self):
        """测试返回类型"""
        letters = ['A', 'B', 'C']
        result = task_func(letters)
        
        assert isinstance(result, dict)
        assert not isinstance(result, OrderedDict)  # 应该是普通dict
    
    def test_deterministic_with_fixed_seed(self):
        """测试使用固定随机种子时的确定性"""
        letters = ['X', 'Y', 'Z']
        
        # 使用固定随机种子
        random.seed(42)
        result1 = task_func(letters)
        
        random.seed(42)
        result2 = task_func(letters)
        
        # 验证两次结果相同
        assert result1 == result2
        
        # 恢复随机种子
        random.seed()