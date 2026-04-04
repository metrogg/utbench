import pytest
from below_zero import below_zero

class TestBelowZero:
    def test_empty_list(self):
        """Test with empty operations list - should return False"""
        assert below_zero([]) == False
    
    def test_all_positive(self):
        """Test with all positive operations - should return False"""
        assert below_zero([1, 2, 3]) == False
        assert below_zero([100, 50, 25]) == False
    
    def test_becomes_negative(self):
        """Test when balance goes negative - should return True"""
        assert below_zero([1, 2, -4, 5]) == True
        assert below_zero([-1]) == True
    
    def test_single_positive(self):
        """Test single positive operation"""
        assert below_zero([5]) == False
    
    def test_single_negative(self):
        """Test single negative operation - balance starts at 0, so negative makes it below zero"""
        assert below_zero([-1]) == True
    
    def test_all_negative(self):
        """Test all negative operations - should return True"""
        assert below_zero([-1, -2, -3]) == True
    
    def test_balance_never_below_zero(self):
        """Test operations that never make balance go below zero"""
        assert below_zero([1, -1, 2, -2, 3]) == False
    
    def test_balance_exactly_zero(self):
        """Test operations that result in balance of exactly zero"""
        assert below_zero([1, -1]) == False
    
    def test_goes_below_then_above(self):
        """Test when balance goes below zero then back above"""
        assert below_zero([1, 2, -5, 10]) == True