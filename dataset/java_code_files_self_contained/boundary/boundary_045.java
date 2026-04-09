// Converted Java method
import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.HashMap;
import java.util.Map;

class TimeDifferenceCalculator {
    private static final String DEFAULT_DATETIME_PATTERN = "yyyy-MM-dd HH:mm:ss";
    
    /**
     * Calculates the time difference between two dates in a human-readable format.
     * Returns a map containing multiple representations of the time difference.
     * 
     * @param earlierTime The earlier time string in "yyyy-MM-dd HH:mm:ss" format
     * @param laterTime The later time string in "yyyy-MM-dd HH:mm:ss" format
     * @return Map containing:
     *         - "humanReadable": Human-readable string (e.g., "2 days ago")
     *         - "seconds": Difference in seconds
     *         - "minutes": Difference in minutes
     *         - "hours": Difference in hours
     *         - "days": Difference in days
     *         - "weeks": Difference in weeks
     * @throws ParseException if the time strings are in incorrect format
     */
    public static Map<String, Object> calculateDetailedTimeDifference(String earlierTime, String laterTime) 
            throws ParseException {
        SimpleDateFormat format = new SimpleDateFormat(DEFAULT_DATETIME_PATTERN);
        Date date1 = format.parse(earlierTime);
        Date date2 = format.parse(laterTime);
        
        long diffInMillis = date2.getTime() - date1.getTime();
        if (diffInMillis < 0) {
            diffInMillis = -diffInMillis;
        }
        
        long seconds = diffInMillis / 1000;
        long minutes = seconds / 60;
        long hours = minutes / 60;
        long days = hours / 24;
        long weeks = days / 7;
        
        Map<String, Object> result = new HashMap<>();
        result.put("seconds", seconds);
        result.put("minutes", minutes);
        result.put("hours", hours);
        result.put("days", days);
        result.put("weeks", weeks);
        
        // Human-readable format
        String humanReadable;
        if (weeks > 4) {
            humanReadable = "more than a month ago";
        } else if (weeks > 0) {
            humanReadable = weeks + (weeks == 1 ? " week ago" : " weeks ago");
        } else if (days > 0) {
            humanReadable = days + (days == 1 ? " day ago" : " days ago");
        } else if (hours > 0) {
            humanReadable = hours + (hours == 1 ? " hour ago" : " hours ago");
        } else if (minutes > 0) {
            humanReadable = minutes + (minutes == 1 ? " minute ago" : " minutes ago");
        } else {
            humanReadable = "just now";
        }
        
        result.put("humanReadable", humanReadable);
        return result;
    }
}
