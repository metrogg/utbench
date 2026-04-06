from __future__ import annotations
from typing import TYPE_CHECKING, cast
import numpy as np
from pandas._libs import NaT, algos as libalgos, internals as libinternals, lib
from pandas._libs.missing import NA
from pandas.util._decorators import cache_readonly
from pandas.core.dtypes.cast import ensure_dtype_can_hold_na, find_common_type
from pandas.core.dtypes.common import is_1d_only_ea_dtype, needs_i8_conversion
from pandas.core.dtypes.concat import concat_compat
from pandas.core.dtypes.dtypes import ExtensionDtype
from pandas.core.dtypes.missing import is_valid_na_for_dtype
from pandas.core.construction import ensure_wrapped_if_datetimelike
from pandas.core.internals.blocks import ensure_block_shape, new_block_2d
from pandas.core.internals.managers import BlockManager, make_na_array
if TYPE_CHECKING:
    from collections.abc import Generator, Sequence
    from pandas._typing import ArrayLike, AxisInt, DtypeObj, Shape
    from pandas import Index
    from pandas.core.internals.blocks import Block, BlockPlacement

def concatenate_managers(mgrs_indexers, axes: list[Index], concat_axis: AxisInt, copy: bool) -> BlockManager:
    needs_copy = copy and concat_axis == 0
    if concat_axis == 0:
        mgrs = _maybe_reindex_columns_na_proxy(axes, mgrs_indexers, needs_copy)
        return mgrs[0].concat_horizontal(mgrs, axes)
    if len(mgrs_indexers) > 0 and mgrs_indexers[0][0].nblocks > 0:
        first_dtype = mgrs_indexers[0][0].blocks[0].dtype
        if first_dtype in [np.float64, np.float32]:
            if all((_is_homogeneous_mgr(mgr, first_dtype) for mgr, _ in mgrs_indexers)) and len(mgrs_indexers) > 1:
                shape = tuple((len(x) for x in axes))
                nb = _concat_homogeneous_fastpath(mgrs_indexers, shape, first_dtype)
                return BlockManager((nb,), axes)
    mgrs = _maybe_reindex_columns_na_proxy(axes, mgrs_indexers, needs_copy)
    if len(mgrs) == 1:
        mgr = mgrs[0]
        out = mgr.copy(deep=False)
        out.axes = axes
        return out
    blocks = []
    values: ArrayLike
    for placement, join_units in _get_combined_plan(mgrs):
        unit = join_units[0]
        blk = unit.block
        if _is_uniform_join_units(join_units):
            vals = [ju.block.values for ju in join_units]
            if not blk.is_extension:
                values = np.concatenate(vals, axis=1)
            elif is_1d_only_ea_dtype(blk.dtype):
                values = concat_compat(vals, axis=0, ea_compat_axis=True)
                values = ensure_block_shape(values, ndim=2)
            else:
                values = concat_compat(vals, axis=1)
            values = ensure_wrapped_if_datetimelike(values)
            fastpath = blk.values.dtype == values.dtype
        else:
            values = _concatenate_join_units(join_units, copy=copy)
            fastpath = False
        if fastpath:
            b = blk.make_block_same_class(values, placement=placement)
        else:
            b = new_block_2d(values, placement=placement)
        blocks.append(b)
    return BlockManager(tuple(blocks), axes)

def _maybe_reindex_columns_na_proxy(axes: list[Index], mgrs_indexers: list[tuple[BlockManager, dict[int, np.ndarray]]], needs_copy: bool) -> list[BlockManager]:
    new_mgrs = []
    for mgr, indexers in mgrs_indexers:
        for i, indexer in indexers.items():
            mgr = mgr.reindex_indexer(axes[i], indexer, axis=i, only_slice=True, allow_dups=True, use_na_proxy=True)
        if needs_copy and (not indexers):
            mgr = mgr.copy(deep=True)
        new_mgrs.append(mgr)
    return new_mgrs

def _is_homogeneous_mgr(mgr: BlockManager, first_dtype: DtypeObj) -> bool:
    if mgr.nblocks != 1:
        return False
    blk = mgr.blocks[0]
    if not (blk.mgr_locs.is_slice_like and blk.mgr_locs.as_slice.step == 1):
        return False
    return blk.dtype == first_dtype

def _concat_homogeneous_fastpath(mgrs_indexers, shape: Shape, first_dtype: np.dtype) -> Block:
    if all((not indexers for _, indexers in mgrs_indexers)):
        arrs = [mgr.blocks[0].values.T for mgr, _ in mgrs_indexers]
        arr = np.concatenate(arrs).T
        bp = libinternals.BlockPlacement(slice(shape[0]))
        nb = new_block_2d(arr, bp)
        return nb
    arr = np.empty(shape, dtype=first_dtype)
    if first_dtype == np.float64:
        take_func = libalgos.take_2d_axis0_float64_float64
    else:
        take_func = libalgos.take_2d_axis0_float32_float32
    start = 0
    for mgr, indexers in mgrs_indexers:
        mgr_len = mgr.shape[1]
        end = start + mgr_len
        if 0 in indexers:
            take_func(mgr.blocks[0].values, indexers[0], arr[:, start:end])
        else:
            arr[:, start:end] = mgr.blocks[0].values
        start += mgr_len
    bp = libinternals.BlockPlacement(slice(shape[0]))
    nb = new_block_2d(arr, bp)
    return nb

def _get_combined_plan(mgrs: list[BlockManager]) -> Generator[tuple[BlockPlacement, list[JoinUnit]]]:
    max_len = mgrs[0].shape[0]
    blknos_list = [mgr.blknos for mgr in mgrs]
    pairs = libinternals.get_concat_blkno_indexers(blknos_list)
    for blknos, bp in pairs:
        units_for_bp = []
        for k, mgr in enumerate(mgrs):
            blkno = blknos[k]
            nb = _get_block_for_concat_plan(mgr, bp, blkno, max_len=max_len)
            unit = JoinUnit(nb)
            units_for_bp.append(unit)
        yield (bp, units_for_bp)

def _get_block_for_concat_plan(mgr: BlockManager, bp: BlockPlacement, blkno: int, *, max_len: int) -> Block:
    blk = mgr.blocks[blkno]
    if len(bp) == len(blk.mgr_locs) and (blk.mgr_locs.is_slice_like and blk.mgr_locs.as_slice.step == 1):
        nb = blk
    else:
        ax0_blk_indexer = mgr.blklocs[bp.indexer]
        slc = lib.maybe_indices_to_slice(ax0_blk_indexer, max_len)
        if isinstance(slc, slice):
            nb = blk.slice_block_columns(slc)
        else:
            nb = blk.take_block_columns(slc)
    return nb

class JoinUnit:

    def __init__(self, block: Block) -> None:
        self.block = block

    def __repr__(self) -> str:
        return f'{type(self).__name__}({self.block!r})'

    def _is_valid_na_for(self, dtype: DtypeObj) -> bool:
        if not self.is_na:
            return False
        blk = self.block
        if blk.dtype.kind == 'V':
            return True
        if blk.dtype == object:
            values = blk.values
            return all((is_valid_na_for_dtype(x, dtype) for x in values.ravel(order='K')))
        na_value = blk.fill_value
        if na_value is NaT and blk.dtype != dtype:
            return False
        if na_value is NA and needs_i8_conversion(dtype):
            return False
        return is_valid_na_for_dtype(na_value, dtype)

    @cache_readonly
    def is_na(self) -> bool:
        blk = self.block
        if blk.dtype.kind == 'V':
            return True
        return False

    def get_reindexed_values(self, empty_dtype: DtypeObj, upcasted_na) -> ArrayLike:
        values: ArrayLike
        if upcasted_na is None and self.block.dtype.kind != 'V':
            return self.block.values
        else:
            fill_value = upcasted_na
            if self._is_valid_na_for(empty_dtype):
                blk_dtype = self.block.dtype
                if blk_dtype == np.dtype('object'):
                    values = cast(np.ndarray, self.block.values)
                    if values.size and values[0, 0] is None:
                        fill_value = None
                return make_na_array(empty_dtype, self.block.shape, fill_value)
            return self.block.values

def _concatenate_join_units(join_units: list[JoinUnit], copy: bool) -> ArrayLike:
    empty_dtype = _get_empty_dtype(join_units)
    has_none_blocks = any((unit.block.dtype.kind == 'V' for unit in join_units))
    upcasted_na = _dtype_to_na_value(empty_dtype, has_none_blocks)
    to_concat = [ju.get_reindexed_values(empty_dtype=empty_dtype, upcasted_na=upcasted_na) for ju in join_units]
    if any((is_1d_only_ea_dtype(t.dtype) for t in to_concat)):
        to_concat = [t if is_1d_only_ea_dtype(t.dtype) else t[0, :] for t in to_concat]
        concat_values = concat_compat(to_concat, axis=0, ea_compat_axis=True)
        concat_values = ensure_block_shape(concat_values, 2)
    else:
        concat_values = concat_compat(to_concat, axis=1)
    return concat_values

def _dtype_to_na_value(dtype: DtypeObj, has_none_blocks: bool):
    if isinstance(dtype, ExtensionDtype):
        return dtype.na_value
    elif dtype.kind in 'mM':
        return dtype.type('NaT')
    elif dtype.kind in 'fc':
        return dtype.type('NaN')
    elif dtype.kind == 'b':
        return None
    elif dtype.kind in 'iu':
        if not has_none_blocks:
            return None
        return np.nan
    elif dtype.kind == 'O':
        return np.nan
    raise NotImplementedError

def _get_empty_dtype(join_units: Sequence[JoinUnit]) -> DtypeObj:
    if lib.dtypes_all_equal([ju.block.dtype for ju in join_units]):
        empty_dtype = join_units[0].block.dtype
        return empty_dtype
    has_none_blocks = any((unit.block.dtype.kind == 'V' for unit in join_units))
    dtypes = [unit.block.dtype for unit in join_units if not unit.is_na]
    dtype = find_common_type(dtypes)
    if has_none_blocks:
        dtype = ensure_dtype_can_hold_na(dtype)
    return dtype

def _is_uniform_join_units(join_units: list[JoinUnit]) -> bool:
    first = join_units[0].block
    if first.dtype.kind == 'V':
        return False
    return all((type(ju.block) is type(first) for ju in join_units)) and all((ju.block.dtype == first.dtype or ju.block.dtype.kind in 'iub' for ju in join_units)) and all((not ju.is_na or ju.block.is_extension for ju in join_units))
