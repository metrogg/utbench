// Converted Java method
import java.util.Calendar;
import java.util.HashMap;
import java.util.Map;

class DateUtils {
    private static final String[] DAYS_OF_WEEK = {
        "SUNDAY", "MONDAY", "TUESDAY", "WEDNESDAY", "THURSDAY", "FRIDAY", "SATURDAY"
    };
    
    /**
     * Performs comprehensive date analysis including:
     * - Day of the week
     * - Whether the year is a leap year
     * - Number of days in the month
     * - Zodiac sign based on the date
     * 
     * @param day Day of month (1-31)
     * @param month Month of year (1-12)
     * @param year Year (1000-9999)
     * @return Map containing all date analysis results
     * @throws IllegalArgumentException if date is invalid
     */
    public static Map<String, String> analyzeDate(int day, int month, int year) {
        if (!isValidDate(day, month, year)) {
            throw new IllegalArgumentException("Invalid date");
        }
        
        Map<String, String> result = new HashMap<>();
        
        // Day of week calculation
        Calendar c = Calendar.getInstance();
        c.set(year, month - 1, day);
        int dayOfWeek = c.get(Calendar.DAY_OF_WEEK) - 1; // Calendar.SUNDAY = 1
        result.put("dayOfWeek", DAYS_OF_WEEK[dayOfWeek]);
        
        // Leap year check
        boolean isLeap = isLeapYear(year);
        result.put("isLeapYear", isLeap ? "Yes" : "No");
        
        // Days in month
        int daysInMonth = getDaysInMonth(month, year);
        result.put("daysInMonth", String.valueOf(daysInMonth));
        
        // Zodiac sign
        String zodiac = getZodiacSign(day, month);
        result.put("zodiacSign", zodiac);
        
        return result;
    }
    
    private static boolean isValidDate(int day, int month, int year) {
        if (year < 1000 || year > 9999) return false;
        if (month < 1 || month > 12) return false;
        
        int maxDays = getDaysInMonth(month, year);
        return day >= 1 && day <= maxDays;
    }
    
    private static boolean isLeapYear(int year) {
        if (year % 4 != 0) return false;
        if (year % 100 != 0) return true;
        return year % 400 == 0;
    }
    
    private static int getDaysInMonth(int month, int year) {
        switch (month) {
            case 2:
                return isLeapYear(year) ? 29 : 28;
            case 4: case 6: case 9: case 11:
                return 30;
            default:
                return 31;
        }
    }
    
    private static String getZodiacSign(int day, int month) {
        if ((month == 3 && day >= 21) || (month == 4 && day <= 19)) return "Aries";
        if ((month == 4 && day >= 20) || (month == 5 && day <= 20)) return "Taurus";
        if ((month == 5 && day >= 21) || (month == 6 && day <= 20)) return "Gemini";
        if ((month == 6 && day >= 21) || (month == 7 && day <= 22)) return "Cancer";
        if ((month == 7 && day >= 23) || (month == 8 && day <= 22)) return "Leo";
        if ((month == 8 && day >= 23) || (month == 9 && day <= 22)) return "Virgo";
        if ((month == 9 && day >= 23) || (month == 10 && day <= 22)) return "Libra";
        if ((month == 10 && day >= 23) || (month == 11 && day <= 21)) return "Scorpio";
        if ((month == 11 && day >= 22) || (month == 12 && day <= 21)) return "Sagittarius";
        if ((month == 12 && day >= 22) || (month == 1 && day <= 19)) return "Capricorn";
        if ((month == 1 && day >= 20) || (month == 2 && day <= 18)) return "Aquarius";
        return "Pisces";
    }
}
