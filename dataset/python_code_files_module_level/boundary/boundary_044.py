from dataclasses import dataclass
from typing import Literal, Optional, Union
import pyarrow as pa
import pyarrow.dataset as ds
import pyarrow.parquet as pq
import datasets
from datasets.builder import Key
from datasets.table import table_cast
logger = datasets.utils.logging.get_logger(__name__)

@dataclass
class ParquetConfig(datasets.BuilderConfig):
    batch_size: Optional[int] = None
    columns: Optional[list[str]] = None
    features: Optional[datasets.Features] = None
    filters: Optional[Union[ds.Expression, list[tuple], list[list[tuple]]]] = None
    fragment_scan_options: Optional[ds.ParquetFragmentScanOptions] = None
    on_bad_files: Literal['error', 'warn', 'skip'] = 'error'

    def __post_init__(self):
        super().__post_init__()

class Parquet(datasets.ArrowBasedBuilder):
    BUILDER_CONFIG_CLASS = ParquetConfig

    def _info(self):
        if self.config.columns is not None and self.config.features is not None and (set(self.config.columns) != set(self.config.features)):
            raise ValueError('The columns and features argument must contain the same columns, but got ', f'{self.config.columns} and {self.config.features}')
        return datasets.DatasetInfo(features=self.config.features)

    def _split_generators(self, dl_manager):
        if not self.config.data_files:
            raise ValueError(f'At least one data file must be specified, but got data_files={self.config.data_files}')
        dl_manager.download_config.extract_on_the_fly = True
        data_files = dl_manager.download(self.config.data_files)
        splits = []
        for split_name, files in data_files.items():
            if self.info.features is None:
                for file in files:
                    try:
                        with open(file, 'rb') as f:
                            self.info.features = datasets.Features.from_arrow_schema(pq.read_schema(f))
                            break
                    except pa.ArrowInvalid as e:
                        if self.config.on_bad_files == 'error':
                            logger.error(f"Failed to read schema from '{file}' with error {type(e).__name__}: {e}")
                            raise
                        elif self.config.on_bad_files == 'warn':
                            logger.warning(f"Skipping bad schema from '{file}'. {type(e).__name__}: {e}`")
                        else:
                            logger.debug(f"Skipping bad schema from '{file}'. {type(e).__name__}: {e}`")
            if self.info.features is None:
                raise ValueError(f'At least one valid data file must be specified, all the data_files are invalid: {self.config.data_files}')
            splits.append(datasets.SplitGenerator(name=split_name, gen_kwargs={'files': files, 'row_groups_list': [None] * len(files)}))
        if self.config.columns is not None and set(self.config.columns) != set(self.info.features):
            self.info.features = datasets.Features({col: feat for col, feat in self.info.features.items() if col in self.config.columns})
        return splits

    def _cast_table(self, pa_table: pa.Table) -> pa.Table:
        if self.info.features is not None:
            pa_table = table_cast(pa_table, self.info.features.arrow_schema)
        return pa_table

    def _generate_shards(self, files, row_groups_list):
        if not row_groups_list:
            yield from files
        else:
            for file, row_groups in zip(files, row_groups_list):
                yield {'fragment_data_file': file, 'fragment_row_groups': row_groups}

    def _generate_more_gen_kwargs(self, files, row_groups_list):
        if not row_groups_list:
            parquet_file_format = ds.ParquetFileFormat(default_fragment_scan_options=self.config.fragment_scan_options)
            for file in files:
                with open(file, 'rb') as f:
                    parquet_fragment = parquet_file_format.make_fragment(f)
                    yield {'files': [file] * parquet_fragment.num_row_groups, 'row_groups_list': [(row_group_id,) for row_group_id in range(parquet_fragment.num_row_groups)]}
        else:
            for file, row_groups in zip(files, row_groups_list):
                yield {'files': [file], 'row_groups_list': [row_groups]}

    def _generate_tables(self, files, row_groups_list):
        if self.config.features is not None and self.config.columns is not None:
            if sorted((field.name for field in self.info.features.arrow_schema)) != sorted(self.config.columns):
                raise ValueError(f"Tried to load parquet data with columns '{self.config.columns}' with mismatching features '{self.info.features}'")
        filter_expr = pq.filters_to_expression(self.config.filters) if isinstance(self.config.filters, list) else self.config.filters
        parquet_file_format = ds.ParquetFileFormat(default_fragment_scan_options=self.config.fragment_scan_options)
        for file_idx, (file, row_groups) in enumerate(zip(files, row_groups_list)):
            try:
                with open(file, 'rb') as f:
                    parquet_fragment = parquet_file_format.make_fragment(f)
                    if row_groups is not None:
                        parquet_fragment.subset(row_group_ids=row_groups)
                    if parquet_fragment.row_groups:
                        batch_size = self.config.batch_size or parquet_fragment.row_groups[0].num_rows
                        for batch_idx, record_batch in enumerate(parquet_fragment.to_batches(batch_size=batch_size, columns=self.config.columns, filter=filter_expr, batch_readahead=0, fragment_readahead=0)):
                            pa_table = pa.Table.from_batches([record_batch])
                            yield (Key(file_idx, batch_idx), self._cast_table(pa_table))
            except (pa.ArrowInvalid, ValueError) as e:
                if self.config.on_bad_files == 'error':
                    logger.error(f"Failed to read file '{file}' with error {type(e).__name__}: {e}")
                    raise
                elif self.config.on_bad_files == 'warn':
                    logger.warning(f"Skipping bad file '{file}'. {type(e).__name__}: {e}`")
                else:
                    logger.debug(f"Skipping bad file '{file}'. {type(e).__name__}: {e}`")
