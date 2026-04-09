// Converted Java method
import java.util.HashMap;
import java.util.Map;
import java.math.BigDecimal;
import java.math.RoundingMode;

class CurrencyConverter {
    private static final Map<String, String> CURRENCY_NAMES = new HashMap<>();
    private static final Map<String, String> ACCOUNT_COL_NAMES = new HashMap<>();
    private static final Map<String, String> DEBTS_COL_NAMES = new HashMap<>();
    private static final Map<String, Integer> DECIMAL_PLACES = new HashMap<>();
    private static final Map<String, BigDecimal> EXCHANGE_RATES = new HashMap<>();

    static {
        // Initialize currency data
        CURRENCY_NAMES.put("CNY", "Chinese Yuan");
        CURRENCY_NAMES.put("TWD", "New Taiwan Dollar");
        CURRENCY_NAMES.put("USD", "US Dollar");
        CURRENCY_NAMES.put("JPY", "Japanese Yen");
        CURRENCY_NAMES.put("EUR", "Euro");
        CURRENCY_NAMES.put("GBP", "British Pound");

        ACCOUNT_COL_NAMES.put("CNY", "sum_rmb");
        ACCOUNT_COL_NAMES.put("TWD", "sum_ntd");
        ACCOUNT_COL_NAMES.put("USD", "sum_usd");
        ACCOUNT_COL_NAMES.put("JPY", "sum_yen");
        ACCOUNT_COL_NAMES.put("EUR", "sum_eur");
        ACCOUNT_COL_NAMES.put("GBP", "sum_gbp");

        DEBTS_COL_NAMES.put("CNY", "debts_rmb");
        DEBTS_COL_NAMES.put("TWD", "debts_ntd");
        DEBTS_COL_NAMES.put("USD", "debts_usd");
        DEBTS_COL_NAMES.put("JPY", "debts_yen");
        DEBTS_COL_NAMES.put("EUR", "debts_eur");
        DEBTS_COL_NAMES.put("GBP", "debts_gbp");

        DECIMAL_PLACES.put("CNY", 2);
        DECIMAL_PLACES.put("TWD", 0);
        DECIMAL_PLACES.put("USD", 2);
        DECIMAL_PLACES.put("JPY", 0);
        DECIMAL_PLACES.put("EUR", 2);
        DECIMAL_PLACES.put("GBP", 2);

        // Initialize exchange rates (relative to USD)
        EXCHANGE_RATES.put("USD", BigDecimal.ONE);
        EXCHANGE_RATES.put("CNY", new BigDecimal("6.5"));
        EXCHANGE_RATES.put("TWD", new BigDecimal("28.0"));
        EXCHANGE_RATES.put("JPY", new BigDecimal("110.0"));
        EXCHANGE_RATES.put("EUR", new BigDecimal("0.85"));
        EXCHANGE_RATES.put("GBP", new BigDecimal("0.75"));
    }

    /**
     * Converts an amount from one currency to another
     * @param amount The amount to convert
     * @param fromCurrency Source currency code
     * @param toCurrency Target currency code
     * @return Converted amount with proper decimal places
     * @throws IllegalArgumentException if currency codes are invalid
     */
    public static BigDecimal convertCurrency(BigDecimal amount, String fromCurrency, String toCurrency) {
        if (!EXCHANGE_RATES.containsKey(fromCurrency) || !EXCHANGE_RATES.containsKey(toCurrency)) {
            throw new IllegalArgumentException("Invalid currency code");
        }

        BigDecimal fromRate = EXCHANGE_RATES.get(fromCurrency);
        BigDecimal toRate = EXCHANGE_RATES.get(toCurrency);
        
        // Convert to USD first, then to target currency
        BigDecimal amountInUSD = amount.divide(fromRate, 10, RoundingMode.HALF_UP);
        BigDecimal convertedAmount = amountInUSD.multiply(toRate);
        
        // Round to appropriate decimal places
        int decimalPlaces = DECIMAL_PLACES.getOrDefault(toCurrency, 2);
        return convertedAmount.setScale(decimalPlaces, RoundingMode.HALF_UP);
    }

    /**
     * Gets the display name for a currency
     * @param currencyCode The currency code
     * @return Display name of the currency
     */
    public static String getCurrencyName(String currencyCode) {
        return CURRENCY_NAMES.getOrDefault(currencyCode, "Unknown Currency");
    }

    /**
     * Gets the account column name for a currency
     * @param currencyCode The currency code
     * @return Column name for account balance
     */
    public static String getAccountColumnName(String currencyCode) {
        return ACCOUNT_COL_NAMES.getOrDefault(currencyCode, "");
    }

    /**
     * Gets the debts column name for a currency
     * @param currencyCode The currency code
     * @return Column name for debts
     */
    public static String getDebtsColumnName(String currencyCode) {
        return DEBTS_COL_NAMES.getOrDefault(currencyCode, "");
    }

    /**
     * Gets the number of decimal places for a currency
     * @param currencyCode The currency code
     * @return Number of decimal places
     */
    public static int getDecimalPlaces(String currencyCode) {
        return DECIMAL_PLACES.getOrDefault(currencyCode, 2);
    }
}
