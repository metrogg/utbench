import hashlib
import binascii

def task_func(salt, cursor):
    if not isinstance(salt, str):
        raise TypeError
    cursor.execute('SELECT id, password FROM users')
    users = cursor.fetchall()
    count_updated = 0
    for user in users:
        password = user[1].encode('utf-8')
        salted_password = password + salt.encode('utf-8')
        hash_obj = hashlib.sha256(salted_password)
        hashed_password = binascii.hexlify(hash_obj.digest()).decode('utf-8')
        cursor.execute(f"UPDATE users SET password = '{hashed_password}' WHERE id = {user[0]}")
        count_updated += 1
    return count_updated
