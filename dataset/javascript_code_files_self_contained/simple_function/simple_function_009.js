function assocIndexOf(array, key) {
      var length = array.length;
      while (length--) {
        if (eq(array[length][0], key)) {
          return length;
        }
      }
      return -1;
    }

function eq(value, other) {
      return value === other || (value !== value && other !== other);
    }