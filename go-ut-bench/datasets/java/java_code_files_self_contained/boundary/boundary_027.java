// Converted Java method
import java.math.BigDecimal;
import java.math.RoundingMode;

class TaxPriceCalculator {
    
    /**
     * Calculates the price and tax amounts based on the total amount and tax rate.
     * Handles both tax-included and tax-excluded scenarios.
     * 
     * @param totalAmount The total amount (could be tax-included or tax-excluded)
     * @param taxRate The tax rate as a percentage (e.g., 13 for 13%)
     * @param isTaxIncluded Flag indicating if totalAmount includes tax
     * @return An array where [0] is price amount, [1] is tax amount
     * @throws IllegalArgumentException if inputs are invalid
     */
    public static BigDecimal[] calculatePriceAndTax(BigDecimal totalAmount, BigDecimal taxRate, boolean isTaxIncluded) {
        if (totalAmount == null || taxRate == null) {
            throw new IllegalArgumentException("Amount and tax rate cannot be null");
        }
        if (totalAmount.compareTo(BigDecimal.ZERO) < 0) {
            throw new IllegalArgumentException("Amount cannot be negative");
        }
        if (taxRate.compareTo(BigDecimal.ZERO) < 0) {
            throw new IllegalArgumentException("Tax rate cannot be negative");
        }

        BigDecimal[] result = new BigDecimal[2];
        BigDecimal rate = taxRate.divide(new BigDecimal("100"), 10, RoundingMode.HALF_UP);

        if (isTaxIncluded) {
            // Calculate price from tax-included amount
            result[0] = totalAmount.divide(BigDecimal.ONE.add(rate), 10, RoundingMode.HALF_UP);
            result[1] = totalAmount.subtract(result[0]);
        } else {
            // Calculate tax from tax-excluded amount
            result[0] = totalAmount;
            result[1] = totalAmount.multiply(rate);
        }

        // Round to 2 decimal places for currency
        result[0] = result[0].setScale(2, RoundingMode.HALF_UP);
        result[1] = result[1].setScale(2, RoundingMode.HALF_UP);

        return result;
    }
}
