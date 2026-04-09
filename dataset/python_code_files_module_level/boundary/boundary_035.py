from __future__ import annotations
import os.path
import types
import zipimport
from collections.abc import Iterable
from typing import TYPE_CHECKING
from coverage import env
from coverage.exceptions import CoverageException, NoSource
from coverage.files import canonical_filename, relative_filename, zip_location
from coverage.misc import isolate_module, join_regex
from coverage.parser import PythonParser
from coverage.phystokens import source_encoding, source_token_lines
from coverage.plugin import CodeRegion, FileReporter
from coverage.regions import code_regions
from coverage.types import TArc, TLineNo, TMorf, TSourceTokenLines
if TYPE_CHECKING:
    from coverage import Coverage
os = isolate_module(os)
open = open

def read_python_source(filename: str) -> bytes:
    with open(filename, 'rb') as f:
        source = f.read()
    return source.replace(b'\r\n', b'\n').replace(b'\r', b'\n')

def get_python_source(filename: str) -> str:
    base, ext = os.path.splitext(filename)
    if ext == '.py' and env.WINDOWS:
        exts = ['.py', '.pyw']
    else:
        exts = [ext]
    source_bytes: bytes | None
    for ext in exts:
        try_filename = base + ext
        if os.path.exists(try_filename):
            source_bytes = read_python_source(try_filename)
            break
        source_bytes = get_zip_bytes(try_filename)
        if source_bytes is not None:
            break
    else:
        raise NoSource(f"No source for code: '{filename}'.", slug='no-source')
    source_bytes = source_bytes.replace(b'\x0c', b' ')
    source = source_bytes.decode(source_encoding(source_bytes), 'replace')
    if source and source[-1] != '\n':
        source += '\n'
    return source

def get_zip_bytes(filename: str) -> bytes | None:
    zipfile_inner = zip_location(filename)
    if zipfile_inner is not None:
        zipfile, inner = zipfile_inner
        try:
            zi = zipimport.zipimporter(zipfile)
        except zipimport.ZipImportError:
            return None
        try:
            data = zi.get_data(inner)
        except OSError:
            return None
        return data
    return None

def source_for_file(filename: str) -> str:
    if filename.endswith('.py'):
        return filename
    elif filename.endswith(('.pyc', '.pyo')):
        py_filename = filename[:-1]
        if os.path.exists(py_filename):
            return py_filename
        if env.WINDOWS:
            pyw_filename = py_filename + 'w'
            if os.path.exists(pyw_filename):
                return pyw_filename
        return py_filename
    return filename

def source_for_morf(morf: TMorf) -> str:
    if hasattr(morf, '__file__') and morf.__file__:
        filename = morf.__file__
    elif isinstance(morf, types.ModuleType):
        raise CoverageException(f'Module {morf} has no file')
    else:
        filename = morf
    filename = source_for_file(filename)
    return filename

class PythonFileReporter(FileReporter):

    def __init__(self, morf: TMorf, coverage: Coverage | None=None) -> None:
        self.coverage = coverage
        filename = source_for_morf(morf)
        fname = filename
        canonicalize = True
        if self.coverage is not None:
            if self.coverage.config.relative_files:
                canonicalize = False
        if canonicalize:
            fname = canonical_filename(filename)
        super().__init__(fname)
        if hasattr(morf, '__name__'):
            name = morf.__name__.replace('.', os.sep)
            if os.path.basename(filename).startswith('__init__.'):
                name += os.sep + '__init__'
            name += '.py'
        else:
            name = relative_filename(filename)
        self.relname = name
        self._source: str | None = None
        self._parser: PythonParser | None = None
        self._excluded = None

    def __repr__(self) -> str:
        return f'<PythonFileReporter {self.filename!r}>'

    def relative_filename(self) -> str:
        return self.relname

    @property
    def parser(self) -> PythonParser:
        assert self.coverage is not None
        if self._parser is None:
            self._parser = PythonParser(filename=self.filename, exclude=self.coverage._exclude_regex('exclude'))
            self._parser.parse_source()
        return self._parser

    def lines(self) -> set[TLineNo]:
        return self.parser.statements

    def multiline_map(self) -> dict[TLineNo, TLineNo]:
        return self.parser.multiline_map

    def excluded_lines(self) -> set[TLineNo]:
        return self.parser.excluded

    def translate_lines(self, lines: Iterable[TLineNo]) -> set[TLineNo]:
        return self.parser.translate_lines(lines)

    def translate_arcs(self, arcs: Iterable[TArc]) -> set[TArc]:
        return self.parser.translate_arcs(arcs)

    def no_branch_lines(self) -> set[TLineNo]:
        assert self.coverage is not None
        no_branch = self.parser.lines_matching(join_regex(self.coverage.config.partial_list + self.coverage.config.partial_always_list))
        return no_branch

    def arcs(self) -> set[TArc]:
        return self.parser.arcs()

    def exit_counts(self) -> dict[TLineNo, int]:
        return self.parser.exit_counts()

    def missing_arc_description(self, start: TLineNo, end: TLineNo, executed_arcs: Iterable[TArc] | None=None) -> str:
        return self.parser.missing_arc_description(start, end)

    def arc_description(self, start: TLineNo, end: TLineNo) -> str:
        return self.parser.arc_description(start, end)

    def source(self) -> str:
        if self._source is None:
            self._source = get_python_source(self.filename)
        return self._source

    def should_be_python(self) -> bool:
        _, ext = os.path.splitext(self.filename)
        if ext.startswith('.py'):
            return True
        if not ext:
            return True
        return False

    def source_token_lines(self) -> TSourceTokenLines:
        return source_token_lines(self.source())

    def code_regions(self) -> Iterable[CodeRegion]:
        return code_regions(self.source())

    def code_region_kinds(self) -> Iterable[tuple[str, str]]:
        return [('function', 'functions'), ('class', 'classes')]
