import time
import threading

def task_func(delay_time: float=1.0, num_threads: int=5):
    results = []

    def delay():
        time.sleep(delay_time)
        results.append(f'Delay in thread {threading.current_thread().name} completed')
    for i in range(num_threads):
        t = threading.Thread(target=delay, name=str(i))
        t.start()
        t.join()
    return results
