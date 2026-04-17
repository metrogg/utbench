function arrayMap(array, iteratee) {
    var index = -1,
        length = array == null ? 0 : array.length,
        result = Array(length);

    while (++index < length) {
      result[index] = iteratee(array[index], index, array);
    }
    return result;
  }

function baseValues(object, props) {
    return arrayMap(props, function(key) {
      return object[key];
    });
  }