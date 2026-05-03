# Example: Testing a function with I/O and error handling

import pytest
from unittest import mock
from module_under_test import read_config


class TestReadConfig:
    """Tests for read_config(path, default=None)."""

    def test_returns_parsed_json(self):
        """Happy path: valid JSON file returns parsed dict."""
        mock_open = mock.mock_open(read_data='{"key": "value"}')
        with mock.patch("builtins.open", mock_open):
            result = read_config("/fake/config.json")
        assert result == {"key": "value"}
        mock_open.assert_called_once_with("/fake/config.json", "r")

    def test_returns_default_when_file_not_found(self):
        """FileNotFoundError returns the default value."""
        with mock.patch("builtins.open", side_effect=FileNotFoundError):
            result = read_config("/missing.json", default={"fallback": True})
        assert result == {"fallback": True}

    def test_raises_on_invalid_json(self):
        """Malformed JSON raises ValueError with descriptive message."""
        mock_open = mock.mock_open(read_data="not json")
        with mock.patch("builtins.open", mock_open):
            with pytest.raises(ValueError, match="Invalid JSON"):
                read_config("/bad/config.json")

    def test_raises_when_no_default_and_file_missing(self):
        """FileNotFoundError with no default raises the original error."""
        with mock.patch("builtins.open", side_effect=FileNotFoundError):
            with pytest.raises(FileNotFoundError):
                read_config("/missing.json")

    @pytest.mark.parametrize("content,expected", [
        ('{}', {}),
        ('{"a": 1, "b": 2}', {"a": 1, "b": 2}),
        ('[]', []),
        ('null', None),
    ])
    def test_various_json_types(self, content, expected):
        """Different valid JSON types are parsed correctly."""
        mock_open = mock.mock_open(read_data=content)
        with mock.patch("builtins.open", mock_open):
            assert read_config("/f.json") == expected
