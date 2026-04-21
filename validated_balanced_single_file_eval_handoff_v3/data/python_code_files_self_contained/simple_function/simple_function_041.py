import binascii
import io
import gzip

def task_func(compressed_hex):
    try:
        compressed_bytes = binascii.unhexlify(compressed_hex)
        decompressed_bytes = gzip.GzipFile(fileobj=io.BytesIO(compressed_bytes)).read()
        decoded_string = decompressed_bytes.decode('utf-8')
        return decoded_string
    except gzip.BadGzipFile as e:
        return 'Error during decompression: ' + str(e)
