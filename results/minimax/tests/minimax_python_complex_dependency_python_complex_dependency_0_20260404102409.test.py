import pytest
from unittest.mock import patch, MagicMock


def task_func(LETTERS):
    import random
    import statistics
    
    random_dict = {k: [random.randint(0, 100) for _ in range(random.randint(1, 10))] for k in LETTERS}
    sorted_dict = dict(sorted(random_dict.items(), key=lambda item: statistics.mean(item[1]), reverse=True))
    return sorted_dict


class TestTaskFunc:
    def test_valid_input_letters(self):
        """Test with valid list of letters."""
        letters = ['a', 'b', 'c']
        mock_random_values = [
            50, 50, 50, 50, 50, 50, 50, 50, 50, 50,  # for 'a': 5 values -> mean=50
            30, 70, 30, 70, 30, 70, 30, 70, 30, 70,  # for 'b': 5 values -> mean=50
            90, 10, 90, 10, 90, 10, 90, 10, 90, 10   # for 'c': 5 values -> mean=50
        ]
        
        with patch('task_func.random.randint') as mock_randint:
            with patch('task_func.random.choice') as mock_choice:
                mock_randint.side_effect = mock_random_values
                
                result = task_func(letters)
                
                assert isinstance(result, dict)
                assert set(result.keys()) == set(letters)

    def test_sorted_by_mean_descending(self):
        """Test that dictionary is sorted by mean values in descending order."""
        letters = ['a', 'b', 'c']
        
        # Mock randint to return specific values for list lengths and elements
        def mock_randint_side_effect(min_val, max_val):
            if min_val == 0:
                # These are the list element values
                if hasattr(mock_randint_side_effect, 'call_count'):
                    pass
            return 1
        
        # We need to control both list lengths and element values
        call_count = [0]
        def randint_mock(min_val, max_val):
            call_count[0] += 1
            if call_count[0] <= 3:
                return 5  # Each letter gets 5 elements in list
            elif call_count[0] <= 6:
                return 30  # 'a' values: 30, 30, 30, 30, 30 -> mean=30
            elif call_count[0] <= 11:
                return 60  # 'b' values: 60, 60, 60, 60, 60 -> mean=60
            else:
                return 90  # 'c' values: 90, 90, 90, 90, 90 -> mean=90
        
        with patch('task_func.random.randint', side_effect=randint_mock):
            result = task_func(letters)
        
        result_means = []
        import statistics
        for key in result:
            result_means.append(statistics.mean(result[key]))
        
        # Check descending order
        assert result_means == sorted(result_means, reverse=True)

    def test_empty_input(self):
        """Test with empty list of letters."""
        result = task_func([])
        assert result == {}
        assert isinstance(result, dict)

    def test_single_letter(self):
        """Test with a single letter input."""
        call_count = [0]
        def randint_mock(min_val, max_val):
            call_count[0] += 1
            if call_count[0] <= 1:
                return 5  # List length
            return 42  # List elements
        
        with patch('task_func.random.randint', side_effect=randint_mock):
            result = task_func(['x'])
        
        assert isinstance(result, dict)
        assert 'x' in result
        assert len(result['x']) == 5
        assert all(v == 42 for v in result['x'])

    def test_invalid_input_type(self):
        """Test with invalid input type (string instead of list)."""
        with pytest.raises(TypeError):
            task_func("abc")

    def test_invalid_input_none(self):
        """Test with None input."""
        with pytest.raises(TypeError):
            task_func(None)

    def test_duplicate_letters(self):
        """Test with duplicate letters in input."""
        letters = ['a', 'a', 'b']
        call_count = [0]
        def randint_mock(min_val, max_val):
            call_count[0] += 1
            if call_count[0] <= 3:
                return 5
            elif call_count[0] <= 8:
                return 50
            else:
                return 75
        
        with patch('task_func.random.randint', side_effect=randint_mock):
            result = task_func(letters)
        
        # Dictionary cannot have duplicate keys, so only one 'a' will exist
        assert 'b' in result

    def test_list_lengths_vary(self):
        """Test that lists can have varying lengths."""
        call_count = [0]
        lengths = {'a': 3, 'b': 7, 'c': 1}
        current_letter = [None]
        current_call_in_letter = [0]
        
        def randint_mock(min_val, max_val):
            call_count[0] += 1
            # First 3 calls are for list lengths (one per letter)
            if call_count[0] <= 3:
                return lengths[['a', 'b', 'c'][call_count[0] - 1]]
            return 50
        
        with patch('task_func.random.randint', side_effect=randint_mock):
            result = task_func(['a', 'b', 'c'])
        
        assert len(result['a']) == 3
        assert len(result['b']) == 7
        assert len(result['c']) == 1

    def test_mean_calculation(self):
        """Test that mean is correctly calculated."""
        call_count = [0]
        def randint_mock(min_val, max_val):
            call_count[0] += 1
            if call_count[0] <= 1:
                return 4  # 'a' has 4 elements
            elif call_count[0] <= 5:
                return 10  # 'a' values: [10, 10, 10, 10] -> mean=10
            elif call_count[0] <= 2:
                return 2  # 'b' has 2 elements
            else:
                return 90  # 'b' values: [90, 90] -> mean=90
        
        with patch('task_func.random.randint', side_effect=randint_mock):
            result = task_func(['a', 'b'])
        
        import statistics
        assert statistics.mean(result['a']) == 10
        assert statistics.mean(result['b']) == 90

    def test_first_key_has_highest_mean(self):
        """Test that the first key in result has the highest mean."""
        letters = ['x', 'y', 'z']
        
        call_count = [0]
        def randint_mock(min_val, max_val):
            call_count[0] += 1
            if call_count[0] <= 3:
                return 3
            if call_count[0] <= 6:
                return 20  # 'x' mean = 20
            elif call_count[0] <= 9:
                return 50  # 'y' mean = 50
            else:
                return 80  # 'z' mean = 80
        
        with patch('task_func.random.randint', side_effect=randint_mock):
            result = task_func(letters)
        
        keys = list(result.keys())
        import statistics
        first_mean = statistics.mean(result[keys[0]])
        
        for key in keys[1:]:
            assert first_mean >= statistics.mean(result[key])

    def test_output_dict_type(self):
        """Test that output is a dict type."""
        with patch('task_func.random.randint', return_value=5):
            result = task_func(['a', 'b'])
        
        assert type(result) is dict

    def test_all_values_in_range(self):
        """Test that all values in lists are within valid range (0-100)."""
        letters = ['a', 'b', 'c', 'd', 'e']
        
        with patch('task_func.random.randint', return_value=5):
            result = task_func(letters)
        
        for key, values in result.items():
            for value in values:
                assert 0 <= value <= 100, f"Value {value} out of range for key {key}"

    def test_mixed_list_lengths(self):
        """Test with various list lengths to ensure sorting works correctly."""
        letters = ['a', 'b', 'c', 'd']
        expected_lengths = [1, 3, 5, 8]
        
        call_count = [0]
        def randint_mock(min_val, max_val):
            call_count[0] += 1
            if call_count[0] <= 4:
                return expected_lengths[call_count[0] - 1]
            return 50
        
        with patch('task_func.random.randint', side_effect=randint_mock):
            result = task_func(letters)
        
        actual_lengths = [len(result[k]) for k in result]
        assert actual_lengths == expected_lengths