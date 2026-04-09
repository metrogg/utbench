import importlib.util
import os
import tempfile
from pathlib import PurePath
from typing import TYPE_CHECKING, NamedTuple, Optional, Union
import fsspec
import numpy as np
from .features import List
from .utils import logging
from .utils import tqdm as hf_tqdm
if TYPE_CHECKING:
    from .arrow_dataset import Dataset
    try:
        from elasticsearch import Elasticsearch
    except ImportError:
        pass
    try:
        import faiss
    except ImportError:
        pass
_has_elasticsearch = importlib.util.find_spec('elasticsearch') is not None
_has_faiss = importlib.util.find_spec('faiss') is not None
logger = logging.get_logger(__name__)

class MissingIndex(Exception):
    pass

class SearchResults(NamedTuple):
    scores: list[float]
    indices: list[int]

class BatchedSearchResults(NamedTuple):
    total_scores: list[list[float]]
    total_indices: list[list[int]]

class NearestExamplesResults(NamedTuple):
    scores: list[float]
    examples: dict

class BatchedNearestExamplesResults(NamedTuple):
    total_scores: list[list[float]]
    total_examples: list[dict]

class BaseIndex:

    def search(self, query, k: int=10, **kwargs) -> SearchResults:
        raise NotImplementedError

    def search_batch(self, queries, k: int=10, **kwargs) -> BatchedSearchResults:
        total_scores, total_indices = ([], [])
        for query in queries:
            scores, indices = self.search(query, k)
            total_scores.append(scores)
            total_indices.append(indices)
        return BatchedSearchResults(total_scores, total_indices)

    def save(self, file: Union[str, PurePath]):
        raise NotImplementedError

    @classmethod
    def load(cls, file: Union[str, PurePath]) -> 'BaseIndex':
        raise NotImplementedError

class ElasticSearchIndex(BaseIndex):

    def __init__(self, host: Optional[str]=None, port: Optional[int]=None, es_client: Optional['Elasticsearch']=None, es_index_name: Optional[str]=None, es_index_config: Optional[dict]=None):
        if not _has_elasticsearch:
            raise ImportError('You must install ElasticSearch to use ElasticSearchIndex. To do so you can run `pip install elasticsearch==7.7.1 for example`')
        if es_client is not None and (host is not None or port is not None):
            raise ValueError('Please specify either `es_client` or `(host, port)`, but not both.')
        host = host or 'localhost'
        port = port or 9200
        import elasticsearch.helpers
        from elasticsearch import Elasticsearch
        self.es_client = es_client if es_client is not None else Elasticsearch([{'host': host, 'port': str(port)}])
        self.es_index_name = es_index_name if es_index_name is not None else 'huggingface_datasets_' + os.path.basename(tempfile.NamedTemporaryFile().name)
        self.es_index_config = es_index_config if es_index_config is not None else {'settings': {'number_of_shards': 1, 'analysis': {'analyzer': {'stop_standard': {'type': 'standard', ' stopwords': '_english_'}}}}, 'mappings': {'properties': {'text': {'type': 'text', 'analyzer': 'standard', 'similarity': 'BM25'}}}}

    def add_documents(self, documents: Union[list[str], 'Dataset'], column: Optional[str]=None):
        index_name = self.es_index_name
        index_config = self.es_index_config
        self.es_client.indices.create(index=index_name, body=index_config)
        number_of_docs = len(documents)
        progress = hf_tqdm(unit='docs', total=number_of_docs)
        successes = 0

        def passage_generator():
            if column is not None:
                for i, example in enumerate(documents):
                    yield {'text': example[column], '_id': i}
            else:
                for i, example in enumerate(documents):
                    yield {'text': example, '_id': i}
        import elasticsearch as es
        for ok, action in es.helpers.streaming_bulk(client=self.es_client, index=index_name, actions=passage_generator()):
            progress.update(1)
            successes += ok
        if successes != len(documents):
            logger.warning(f'Some documents failed to be added to ElasticSearch. Failures: {len(documents) - successes}/{len(documents)}')
        logger.info(f'Indexed {successes:d} documents')

    def search(self, query: str, k=10, **kwargs) -> SearchResults:
        response = self.es_client.search(index=self.es_index_name, body={'query': {'multi_match': {'query': query, 'fields': ['text'], 'type': 'cross_fields'}}, 'size': k}, **kwargs)
        hits = response['hits']['hits']
        return SearchResults([hit['_score'] for hit in hits], [int(hit['_id']) for hit in hits])

    def search_batch(self, queries, k: int=10, max_workers=10, **kwargs) -> BatchedSearchResults:
        import concurrent.futures
        total_scores, total_indices = ([None] * len(queries), [None] * len(queries))
        with concurrent.futures.ThreadPoolExecutor(max_workers=max_workers) as executor:
            future_to_index = {executor.submit(self.search, query, k, **kwargs): i for i, query in enumerate(queries)}
            for future in concurrent.futures.as_completed(future_to_index):
                index = future_to_index[future]
                results: SearchResults = future.result()
                total_scores[index] = results.scores
                total_indices[index] = results.indices
        return BatchedSearchResults(total_indices=total_indices, total_scores=total_scores)

class FaissIndex(BaseIndex):

    def __init__(self, device: Optional[Union[int, list[int]]]=None, string_factory: Optional[str]=None, metric_type: Optional[int]=None, custom_index: Optional['faiss.Index']=None):
        if string_factory is not None and custom_index is not None:
            raise ValueError('Please specify either `string_factory` or `custom_index` but not both.')
        if device is not None and custom_index is not None:
            raise ValueError("Cannot pass both 'custom_index' and 'device'. Pass 'custom_index' already transferred to the target device instead.")
        self.device = device
        self.string_factory = string_factory
        self.metric_type = metric_type
        self.faiss_index = custom_index
        if not _has_faiss:
            raise ImportError('You must install Faiss to use FaissIndex. To do so you can run `conda install -c pytorch faiss-cpu` or `conda install -c pytorch faiss-gpu`. A community supported package is also available on pypi: `pip install faiss-cpu` or `pip install faiss-gpu`. Note that pip may not have the latest version of FAISS, and thus, some of the latest features and bug fixes may not be available.')

    def add_vectors(self, vectors: Union[np.array, 'Dataset'], column: Optional[str]=None, batch_size: int=1000, train_size: Optional[int]=None, faiss_verbose: Optional[bool]=None):
        import faiss
        if column and (not isinstance(vectors.features[column], List)):
            raise ValueError(f"Wrong feature type for column '{column}'. Expected 1d array, got {vectors.features[column]}")
        if self.faiss_index is None:
            size = len(vectors[0]) if column is None else len(vectors[0][column])
            if self.string_factory is not None:
                if self.metric_type is None:
                    index = faiss.index_factory(size, self.string_factory)
                else:
                    index = faiss.index_factory(size, self.string_factory, self.metric_type)
            elif self.metric_type is None:
                index = faiss.IndexFlat(size)
            else:
                index = faiss.IndexFlat(size, self.metric_type)
            self.faiss_index = self._faiss_index_to_device(index, self.device)
            logger.info(f'Created faiss index of type {type(self.faiss_index)}')
        if faiss_verbose is not None:
            self.faiss_index.verbose = faiss_verbose
            if hasattr(self.faiss_index, 'index') and self.faiss_index.index is not None:
                self.faiss_index.index.verbose = faiss_verbose
            if hasattr(self.faiss_index, 'quantizer') and self.faiss_index.quantizer is not None:
                self.faiss_index.quantizer.verbose = faiss_verbose
            if hasattr(self.faiss_index, 'clustering_index') and self.faiss_index.clustering_index is not None:
                self.faiss_index.clustering_index.verbose = faiss_verbose
        if train_size is not None:
            train_vecs = vectors[:train_size] if column is None else vectors[:train_size][column]
            logger.info(f'Training the index with the first {len(train_vecs)} vectors')
            self.faiss_index.train(train_vecs)
        else:
            logger.info('Ignored the training step of the faiss index as `train_size` is None.')
        logger.info(f'Adding {len(vectors)} vectors to the faiss index')
        for i in hf_tqdm(range(0, len(vectors), batch_size)):
            vecs = vectors[i:i + batch_size] if column is None else vectors[i:i + batch_size][column]
            self.faiss_index.add(vecs)

    @staticmethod
    def _faiss_index_to_device(index: 'faiss.Index', device: Optional[Union[int, list[int]]]=None) -> 'faiss.Index':
        if device is None:
            return index
        import faiss
        if isinstance(device, int):
            if device > -1:
                faiss_res = faiss.StandardGpuResources()
                index = faiss.index_cpu_to_gpu(faiss_res, device, index)
            else:
                index = faiss.index_cpu_to_all_gpus(index)
        elif isinstance(device, (list, tuple)):
            index = faiss.index_cpu_to_gpus_list(index, gpus=list(device))
        else:
            raise TypeError(f'The argument type: {type(device)} is not expected. ' + 'Please pass in either nothing, a positive int, a negative int, or a list of positive ints.')
        return index

    def search(self, query: np.array, k=10, **kwargs) -> SearchResults:
        if len(query.shape) != 1 and (len(query.shape) != 2 or query.shape[0] != 1):
            raise ValueError('Shape of query is incorrect, it has to be either a 1D array or 2D (1, N)')
        queries = query.reshape(1, -1)
        if not queries.flags.c_contiguous:
            queries = np.asarray(queries, order='C')
        scores, indices = self.faiss_index.search(queries, k, **kwargs)
        return SearchResults(scores[0], indices[0].astype(int))

    def search_batch(self, queries: np.array, k=10, **kwargs) -> BatchedSearchResults:
        if len(queries.shape) != 2:
            raise ValueError('Shape of query must be 2D')
        if not queries.flags.c_contiguous:
            queries = np.asarray(queries, order='C')
        scores, indices = self.faiss_index.search(queries, k, **kwargs)
        return BatchedSearchResults(scores, indices.astype(int))

    def save(self, file: Union[str, PurePath], storage_options: Optional[dict]=None):
        import faiss
        if self.device is not None and isinstance(self.device, (int, list, tuple)):
            index = faiss.index_gpu_to_cpu(self.faiss_index)
        else:
            index = self.faiss_index
        with fsspec.open(str(file), 'wb', **storage_options or {}) as f:
            faiss.write_index(index, faiss.BufferedIOWriter(faiss.PyCallbackIOWriter(f.write)))

    @classmethod
    def load(cls, file: Union[str, PurePath], device: Optional[Union[int, list[int]]]=None, storage_options: Optional[dict]=None) -> 'FaissIndex':
        import faiss
        faiss_index = cls(device=device)
        with fsspec.open(str(file), 'rb', **storage_options or {}) as f:
            index = faiss.read_index(faiss.BufferedIOReader(faiss.PyCallbackIOReader(f.read)))
        faiss_index.faiss_index = faiss_index._faiss_index_to_device(index, faiss_index.device)
        return faiss_index

class IndexableMixin:

    def __init__(self):
        self._indexes: dict[str, BaseIndex] = {}

    def __len__(self):
        raise NotImplementedError

    def __getitem__(self, key):
        raise NotImplementedError

    def is_index_initialized(self, index_name: str) -> bool:
        return index_name in self._indexes

    def _check_index_is_initialized(self, index_name: str):
        if not self.is_index_initialized(index_name):
            raise MissingIndex(f"Index with index_name '{index_name}' not initialized yet. Please make sure that you call `add_faiss_index` or `add_elasticsearch_index` first.")

    def list_indexes(self) -> list[str]:
        return list(self._indexes)

    def get_index(self, index_name: str) -> BaseIndex:
        self._check_index_is_initialized(index_name)
        return self._indexes[index_name]

    def add_faiss_index(self, column: str, index_name: Optional[str]=None, device: Optional[Union[int, list[int]]]=None, string_factory: Optional[str]=None, metric_type: Optional[int]=None, custom_index: Optional['faiss.Index']=None, batch_size: int=1000, train_size: Optional[int]=None, faiss_verbose: bool=False):
        index_name = index_name if index_name is not None else column
        faiss_index = FaissIndex(device=device, string_factory=string_factory, metric_type=metric_type, custom_index=custom_index)
        faiss_index.add_vectors(self, column=column, batch_size=batch_size, train_size=train_size, faiss_verbose=faiss_verbose)
        self._indexes[index_name] = faiss_index

    def add_faiss_index_from_external_arrays(self, external_arrays: np.array, index_name: str, device: Optional[Union[int, list[int]]]=None, string_factory: Optional[str]=None, metric_type: Optional[int]=None, custom_index: Optional['faiss.Index']=None, batch_size: int=1000, train_size: Optional[int]=None, faiss_verbose: bool=False):
        faiss_index = FaissIndex(device=device, string_factory=string_factory, metric_type=metric_type, custom_index=custom_index)
        faiss_index.add_vectors(external_arrays, column=None, batch_size=batch_size, train_size=train_size, faiss_verbose=faiss_verbose)
        self._indexes[index_name] = faiss_index

    def save_faiss_index(self, index_name: str, file: Union[str, PurePath], storage_options: Optional[dict]=None):
        index = self.get_index(index_name)
        if not isinstance(index, FaissIndex):
            raise ValueError(f"Index '{index_name}' is not a FaissIndex but a '{type(index)}'")
        index.save(file, storage_options=storage_options)
        logger.info(f'Saved FaissIndex {index_name} at {file}')

    def load_faiss_index(self, index_name: str, file: Union[str, PurePath], device: Optional[Union[int, list[int]]]=None, storage_options: Optional[dict]=None):
        index = FaissIndex.load(file, device=device, storage_options=storage_options)
        if index.faiss_index.ntotal != len(self):
            raise ValueError(f"Index size should match Dataset size, but Index '{index_name}' at {file} has {index.faiss_index.ntotal} elements while the dataset has {len(self)} examples.")
        self._indexes[index_name] = index
        logger.info(f'Loaded FaissIndex {index_name} from {file}')

    def add_elasticsearch_index(self, column: str, index_name: Optional[str]=None, host: Optional[str]=None, port: Optional[int]=None, es_client: Optional['Elasticsearch']=None, es_index_name: Optional[str]=None, es_index_config: Optional[dict]=None):
        index_name = index_name if index_name is not None else column
        es_index = ElasticSearchIndex(host=host, port=port, es_client=es_client, es_index_name=es_index_name, es_index_config=es_index_config)
        es_index.add_documents(self, column=column)
        self._indexes[index_name] = es_index

    def load_elasticsearch_index(self, index_name: str, es_index_name: str, host: Optional[str]=None, port: Optional[int]=None, es_client: Optional['Elasticsearch']=None, es_index_config: Optional[dict]=None):
        self._indexes[index_name] = ElasticSearchIndex(host=host, port=port, es_client=es_client, es_index_name=es_index_name, es_index_config=es_index_config)

    def drop_index(self, index_name: str):
        del self._indexes[index_name]

    def search(self, index_name: str, query: Union[str, np.array], k: int=10, **kwargs) -> SearchResults:
        self._check_index_is_initialized(index_name)
        return self._indexes[index_name].search(query, k, **kwargs)

    def search_batch(self, index_name: str, queries: Union[list[str], np.array], k: int=10, **kwargs) -> BatchedSearchResults:
        self._check_index_is_initialized(index_name)
        return self._indexes[index_name].search_batch(queries, k, **kwargs)

    def get_nearest_examples(self, index_name: str, query: Union[str, np.array], k: int=10, **kwargs) -> NearestExamplesResults:
        self._check_index_is_initialized(index_name)
        scores, indices = self.search(index_name, query, k, **kwargs)
        top_indices = [i for i in indices if i >= 0]
        return NearestExamplesResults(scores[:len(top_indices)], self[top_indices])

    def get_nearest_examples_batch(self, index_name: str, queries: Union[list[str], np.array], k: int=10, **kwargs) -> BatchedNearestExamplesResults:
        self._check_index_is_initialized(index_name)
        total_scores, total_indices = self.search_batch(index_name, queries, k, **kwargs)
        total_scores = [scores_i[:len([i for i in indices_i if i >= 0])] for scores_i, indices_i in zip(total_scores, total_indices)]
        total_samples = [self[[i for i in indices if i >= 0]] for indices in total_indices]
        return BatchedNearestExamplesResults(total_scores, total_samples)
