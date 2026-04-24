import java.util.*;
import java.lang.*;

class Solution {
    /**
    You have to write a function which validates a given date string and
    returns true if the date is valid otherwise false.
    The date is valid if all of the following rules are satisfied:
    1. The date string is not empty.
    2. The number of days is not less than 1 or higher than 31 days for months 1,3,5,7,8,10,12. And the number of days is not less than 1 or higher than 30 days for months 4,6,9,11. And, the number of days is not less than 1 or higher than 29 for the month 2.
    3. The months should not be less than 1 or higher than 12.
    4. The date should be in the format: mm-dd-yyyy

    for example:
    validDate("03-11-2000") => true
    validDate("15-01-2012") => false
    validDate("04-0-2040") => false
    validDate("06-04-2020") => true
    validDate("06/04/2020") => false
     */
    public boolean validDate(String date) {
        try {
            date = date.strip();
            String[] dates = date.split("-" );
            String m = dates[0];
            while (!m.isEmpty() && m.charAt(0) == '0') {
                m = m.substring(1);
            }
            String d = dates[1];
            while (!d.isEmpty() && d.charAt(0) == '0') {
                d = d.substring(1);
            }
            String y = dates[2];
            while (!y.isEmpty() && y.charAt(0) == '0') {
                y = y.substring(1);
            }
            int month = Integer.parseInt(m), day = Integer.parseInt(d), year = Integer.parseInt(y);
            if (month < 1 || month > 12) {
                return false;
            }
            if (Arrays.asList(1, 3, 5, 7, 8, 10, 12).contains(month) && (day < 1 || day > 31)) {
                return false;
            }
            if (Arrays.asList(4, 6, 9, 11).contains(month) && (day < 1 || day > 30)) {
                return false;
            }
            if (month == 2 && (day < 1 || day > 29)) {
                return false;
            }
            return true;
        } catch (Exception e) {
            return false;
        }
    }
}
