import os
import warnings
from functools import partial
from math import ceil
from uuid import uuid4
import numpy as np
import pyarrow as pa
from multiprocess import get_context
try:
    from multiprocess.shared_memory import SharedMemory
except ImportError:
    SharedMemory = None
from .. import config

def minimal_tf_collate_fn(features):
    if isinstance(features, dict):
        return features
    elif config.TF_AVAILABLE:
        import tensorflow as tf
    else:
        raise ImportError('Called a Tensorflow-specific function but Tensorflow is not installed.')
    first = features[0]
    batch = {}
    for k, v in first.items():
        if isinstance(v, np.ndarray):
            batch[k] = np.stack([f[k] for f in features])
        elif isinstance(v, tf.Tensor):
            batch[k] = tf.stack([f[k] for f in features])
        else:
            batch[k] = np.array([f[k] for f in features])
    return batch

def minimal_tf_collate_fn_with_renaming(features):
    batch = minimal_tf_collate_fn(features)
    if 'label' in batch:
        batch['labels'] = batch['label']
        del batch['label']
    return batch

def is_numeric_pa_type(pa_type):
    if pa.types.is_list(pa_type):
        return is_numeric_pa_type(pa_type.value_type)
    return pa.types.is_integer(pa_type) or pa.types.is_floating(pa_type) or pa.types.is_decimal(pa_type)

def np_get_batch(indices, dataset, cols_to_retain, collate_fn, collate_fn_args, columns_to_np_types, return_dict=False):
    if not isinstance(indices, np.ndarray):
        indices = indices.numpy()
    is_batched = True
    if isinstance(indices, np.integer):
        batch = dataset[indices.item()]
        is_batched = False
    elif np.all(np.diff(indices) == 1):
        batch = dataset[indices[0]:indices[-1] + 1]
    elif isinstance(indices, np.ndarray):
        batch = dataset[indices]
    else:
        raise RuntimeError(f'Unexpected type for indices: {type(indices)}')
    if cols_to_retain is not None:
        batch = {key: value for key, value in batch.items() if key in cols_to_retain or key in ('label', 'label_ids', 'labels')}
    if is_batched:
        actual_size = len(list(batch.values())[0])
        batch = [{key: value[i] for key, value in batch.items()} for i in range(actual_size)]
    batch = collate_fn(batch, **collate_fn_args)
    if return_dict:
        out_batch = {}
        for col, cast_dtype in columns_to_np_types.items():
            array = np.array(batch[col])
            array = array.astype(cast_dtype)
            out_batch[col] = array
    else:
        out_batch = []
        for col, cast_dtype in columns_to_np_types.items():
            array = np.array(batch[col])
            array = array.astype(cast_dtype)
            out_batch.append(array)
    return out_batch

def dataset_to_tf(dataset, cols_to_retain, collate_fn, collate_fn_args, columns_to_np_types, output_signature, shuffle, batch_size, drop_remainder):
    if config.TF_AVAILABLE:
        import tensorflow as tf
    else:
        raise ImportError('Called a Tensorflow-specific function but Tensorflow is not installed.')
    if hasattr(tf, 'random_index_shuffle'):
        random_index_shuffle = tf.random_index_shuffle
    elif hasattr(tf.random.experimental, 'index_shuffle'):
        random_index_shuffle = tf.random.experimental.index_shuffle
    else:
        if len(dataset) > 10000000:
            warnings.warn('to_tf_dataset() can be memory-inefficient on versions of TensorFlow older than 2.9. If you are iterating over a dataset with a very large number of samples, consider upgrading to TF >= 2.9.')
        random_index_shuffle = None
    getter_fn = partial(np_get_batch, dataset=dataset, cols_to_retain=cols_to_retain, collate_fn=collate_fn, collate_fn_args=collate_fn_args, columns_to_np_types=columns_to_np_types, return_dict=False)
    tout = [tf.dtypes.as_dtype(dtype) for dtype in columns_to_np_types.values()]

    @tf.function(input_signature=[tf.TensorSpec(None, tf.int64)])
    def fetch_function(indices):
        output = tf.py_function(getter_fn, inp=[indices], Tout=tout)
        return {key: output[i] for i, key in enumerate(columns_to_np_types.keys())}
    tf_dataset = tf.data.Dataset.range(len(dataset))
    if shuffle and random_index_shuffle is not None:
        base_seed = tf.fill((3,), value=tf.cast(-1, dtype=tf.int64))

        def scan_random_index(state, index):
            if tf.reduce_all(state == -1):
                state = tf.random.uniform(shape=(3,), maxval=2 ** 62, dtype=tf.int64)
            shuffled_index = random_index_shuffle(index=index, seed=state, max_index=len(dataset) - 1)
            return (state, shuffled_index)
        tf_dataset = tf_dataset.scan(base_seed, scan_random_index)
    elif shuffle:
        tf_dataset = tf_dataset.shuffle(tf_dataset.cardinality())
    if batch_size is not None:
        tf_dataset = tf_dataset.batch(batch_size, drop_remainder=drop_remainder)
    tf_dataset = tf_dataset.map(fetch_function)
    if batch_size is not None:

        def ensure_shapes(input_dict):
            return {key: tf.ensure_shape(val, output_signature[key].shape) for key, val in input_dict.items()}
    else:

        def ensure_shapes(input_dict):
            return {key: tf.ensure_shape(val, output_signature[key].shape[1:]) for key, val in input_dict.items()}
    return tf_dataset.map(ensure_shapes)

class SharedMemoryContext:

    def __init__(self):
        self.created_shms = []
        self.opened_shms = []

    def get_shm(self, name, size, create):
        shm = SharedMemory(size=int(size), name=name, create=create)
        if create:
            self.created_shms.append(shm)
        else:
            self.opened_shms.append(shm)
        return shm

    def get_array(self, name, shape, dtype, create):
        shm = self.get_shm(name=name, size=np.prod(shape) * np.dtype(dtype).itemsize, create=create)
        return np.ndarray(shape, dtype=dtype, buffer=shm.buf)

    def __enter__(self):
        return self

    def __exit__(self, exc_type, exc_value, traceback):
        for shm in self.created_shms:
            shm.close()
            shm.unlink()
        for shm in self.opened_shms:
            shm.close()

class NumpyMultiprocessingGenerator:

    def __init__(self, dataset, cols_to_retain, collate_fn, collate_fn_args, columns_to_np_types, output_signature, shuffle, batch_size, drop_remainder, num_workers):
        self.dataset = dataset
        self.cols_to_retain = cols_to_retain
        self.collate_fn = collate_fn
        self.collate_fn_args = collate_fn_args
        self.string_columns = [col for col, dtype in columns_to_np_types.items() if dtype is np.str_]
        self.columns_to_np_types = {col: dtype if col not in self.string_columns else np.dtype('U1') for col, dtype in columns_to_np_types.items()}
        self.output_signature = output_signature
        self.shuffle = shuffle
        self.batch_size = batch_size
        self.drop_remainder = drop_remainder
        self.num_workers = num_workers
        self.columns_to_ranks = {col: int(spec.shape.rank) if col not in self.string_columns else int(spec.shape.rank) + 1 for col, spec in output_signature.items()}

    def __iter__(self):
        num_workers = min(self.num_workers, int(ceil(len(self.dataset) / self.batch_size)))
        per_worker_batches, final_batch, final_batch_worker = self.distribute_batches(self.dataset, self.batch_size, self.drop_remainder, num_workers, self.shuffle)
        ctx = get_context('spawn')
        names = []
        shape_arrays = []
        workers = []
        array_ready_events = [ctx.Event() for _ in range(num_workers)]
        array_loaded_events = [ctx.Event() for _ in range(num_workers)]
        base_args = {'dataset': self.dataset, 'cols_to_retain': self.cols_to_retain, 'collate_fn': self.collate_fn, 'collate_fn_args': self.collate_fn_args, 'columns_to_np_types': self.columns_to_np_types, 'columns_to_ranks': self.columns_to_ranks, 'string_columns': self.string_columns}
        with SharedMemoryContext() as shm_ctx:
            for i in range(num_workers):
                worker_random_id = str(uuid4())
                worker_name = f'dw_{i}_{worker_random_id}'[:10]
                names.append(worker_name)
                worker_shape_arrays = {col: shm_ctx.get_array(f'{worker_name}_{col}_shape', shape=(rank,), dtype=np.int64, create=True) for col, rank in self.columns_to_ranks.items()}
                shape_arrays.append(worker_shape_arrays)
                worker_indices = per_worker_batches[i]
                if i == final_batch_worker and final_batch is not None:
                    final_batch_arg = final_batch
                else:
                    final_batch_arg = None
                worker_kwargs = {'worker_name': worker_name, 'indices': worker_indices, 'extra_batch': final_batch_arg, 'array_ready_event': array_ready_events[i], 'array_loaded_event': array_loaded_events[i], **base_args}
                worker = ctx.Process(target=self.worker_loop, kwargs=worker_kwargs, daemon=True)
                worker.start()
                workers.append(worker)
            end_signal_received = False
            while not end_signal_received:
                for i in range(num_workers):
                    if not array_ready_events[i].wait(timeout=60):
                        raise TimeoutError('Data loading worker timed out!')
                    array_ready_events[i].clear()
                    array_shapes = shape_arrays[i]
                    if any((np.any(shape < 0) for shape in array_shapes.values())):
                        end_signal_received = True
                        break
                    with SharedMemoryContext() as batch_shm_ctx:
                        arrays = {col: batch_shm_ctx.get_array(f'{names[i]}_{col}', shape=shape, dtype=self.columns_to_np_types[col], create=False) for col, shape in array_shapes.items()}
                        arrays = {col: np.copy(arr) for col, arr in arrays.items()}
                        for string_col in self.string_columns:
                            arrays[string_col] = arrays[string_col].view(f'U{arrays[string_col].shape[-1]}').squeeze(-1)
                    yield arrays
                    array_loaded_events[i].set()
            for worker in workers:
                worker.join()

    def __call__(self):
        return self

    @staticmethod
    def worker_loop(dataset, cols_to_retain, collate_fn, collate_fn_args, columns_to_np_types, columns_to_ranks, string_columns, indices, extra_batch, worker_name, array_ready_event, array_loaded_event):
        os.environ['TF_CPP_MIN_LOG_LEVEL'] = '3'
        if config.TF_AVAILABLE:
            import tensorflow as tf
        else:
            raise ImportError('Called a Tensorflow-specific function but Tensorflow is not installed.')
        tf.config.set_visible_devices([], 'GPU')

        def send_batch_to_parent(indices):
            batch = np_get_batch(indices=indices, dataset=dataset, cols_to_retain=cols_to_retain, collate_fn=collate_fn, collate_fn_args=collate_fn_args, columns_to_np_types=columns_to_np_types, return_dict=True)
            out_arrays = {}
            with SharedMemoryContext() as batch_shm_ctx:
                for col, cast_dtype in columns_to_np_types.items():
                    array = batch[col]
                    if col in string_columns:
                        array = array.view('U1').reshape(array.shape + (-1,))
                    shape_arrays[col][:] = array.shape
                    out_arrays[col] = batch_shm_ctx.get_array(f'{worker_name}_{col}', shape=array.shape, dtype=cast_dtype, create=True)
                    out_arrays[col][:] = array
                array_ready_event.set()
                array_loaded_event.wait()
                array_loaded_event.clear()
        with SharedMemoryContext() as shm_ctx:
            shape_arrays = {col: shm_ctx.get_array(f'{worker_name}_{col}_shape', shape=(rank,), dtype=np.int64, create=False) for col, rank in columns_to_ranks.items()}
            for batch in indices:
                send_batch_to_parent(batch)
            if extra_batch is not None:
                send_batch_to_parent(extra_batch)
            for col, array in shape_arrays.items():
                array[:] = -1
            array_ready_event.set()

    @staticmethod
    def distribute_batches(dataset, batch_size, drop_remainder, num_workers, shuffle):
        indices = np.arange(len(dataset))
        if shuffle:
            np.random.shuffle(indices)
        num_samples = len(indices)
        incomplete_batch_cutoff = num_samples - num_samples % batch_size
        indices, last_incomplete_batch = np.split(indices, [incomplete_batch_cutoff])
        if drop_remainder or len(last_incomplete_batch) == 0:
            last_incomplete_batch = None
        indices = indices.reshape(-1, batch_size)
        num_batches = len(indices)
        final_batches_cutoff = num_batches - num_batches % num_workers
        indices, final_batches = np.split(indices, [final_batches_cutoff])
        indices = indices.reshape(-1, num_workers, batch_size)
        per_worker_indices = np.split(indices, indices.shape[1], axis=1)
        per_worker_indices = [np.squeeze(worker_indices, 1) for worker_indices in per_worker_indices]
        for i in range(len(final_batches)):
            per_worker_indices[i] = np.concatenate([per_worker_indices[i], final_batches[i].reshape(1, -1)], axis=0)
        if last_incomplete_batch is not None:
            incomplete_batch_worker_idx = len(final_batches)
        else:
            incomplete_batch_worker_idx = None
        return (per_worker_indices, last_incomplete_batch, incomplete_batch_worker_idx)

def multiprocess_dataset_to_tf(dataset, cols_to_retain, collate_fn, collate_fn_args, columns_to_np_types, output_signature, shuffle, batch_size, drop_remainder, num_workers):
    if config.TF_AVAILABLE:
        import tensorflow as tf
    else:
        raise ImportError('Called a Tensorflow-specific function but Tensorflow is not installed.')
    data_generator = NumpyMultiprocessingGenerator(dataset=dataset, cols_to_retain=cols_to_retain, collate_fn=collate_fn, collate_fn_args=collate_fn_args, columns_to_np_types=columns_to_np_types, output_signature=output_signature, shuffle=shuffle, batch_size=batch_size, drop_remainder=drop_remainder, num_workers=num_workers)
    tf_dataset = tf.data.Dataset.from_generator(data_generator, output_signature=output_signature)
    if drop_remainder:
        dataset_length = int(len(dataset) // batch_size)
    else:
        dataset_length = int(ceil(len(dataset) / batch_size))
    return tf_dataset.apply(tf.data.experimental.assert_cardinality(dataset_length))
