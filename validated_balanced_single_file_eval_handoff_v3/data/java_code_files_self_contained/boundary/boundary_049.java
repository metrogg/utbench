import java.math.BigInteger;
import java.security.*;
import java.util.*;
import java.lang.*;

class Solution {
    /**
    Given a string "text", return its md5 hash equivalent string with length being 32.
    If "text" is an empty string, return Optional.empty().
    
    >>> stringToMd5("Hello world") == "3e25960a79dbc69b674cd4ec67a72c62"
     */
    public Optional<String> stringToMd5(String text) throws NoSuchAlgorithmException {
        if (text.isEmpty()) {
            return Optional.empty();
        }

        String md5 = new BigInteger(1, java.security.MessageDigest.getInstance("MD5").digest(text.getBytes())).toString(16);
        md5 = "0".repeat(32 - md5.length()) + md5;
        return Optional.of(md5);
    }
}
