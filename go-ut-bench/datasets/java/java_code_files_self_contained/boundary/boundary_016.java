// Converted Java method
import java.util.HashMap;
import java.util.Map;

class Venue {
    private int numRows;
    private int numCols;
    private Map<String, Boolean> seats;

    /**
     * Creates a new Venue with the specified number of rows and columns.
     * @param numRows Number of rows (must be positive)
     * @param numCols Number of columns (must be positive)
     * @throws IllegalArgumentException if rows or columns are not positive
     */
    public Venue(int numRows, int numCols) {
        if (numRows <= 0 || numCols <= 0) {
            throw new IllegalArgumentException("Rows and columns must be positive");
        }
        this.numRows = numRows;
        this.numCols = numCols;
        this.seats = new HashMap<>();
        initializeSeats();
    }

    private void initializeSeats() {
        for (int row = 1; row <= numRows; row++) {
            for (int col = 1; col <= numCols; col++) {
                seats.put(generateSeatKey(row, col), false);
            }
        }
    }

    private String generateSeatKey(int row, int col) {
        return row + "-" + col;
    }

    /**
     * Gets the total number of seats in the venue
     * @return Total number of seats
     */
    public int getNumberOfSeats() {
        return numRows * numCols;
    }

    /**
     * Reserves a specific seat in the venue
     * @param row Row number (1-based)
     * @param col Column number (1-based)
     * @return true if seat was successfully reserved, false if already reserved
     * @throws IllegalArgumentException if seat coordinates are invalid
     */
    public boolean reserveSeat(int row, int col) {
        String seatKey = generateSeatKey(row, col);
        if (!seats.containsKey(seatKey)) {
            throw new IllegalArgumentException("Invalid seat coordinates");
        }
        if (seats.get(seatKey)) {
            return false;
        }
        seats.put(seatKey, true);
        return true;
    }

    /**
     * Gets the number of available seats
     * @return Count of available seats
     */
    public int getAvailableSeats() {
        return (int) seats.values().stream().filter(occupied -> !occupied).count();
    }

    /**
     * Gets the best available seat (closest to front-center)
     * @return int array with [row, col], or null if no seats available
     */
    public int[] getBestAvailableSeat() {
        for (int row = 1; row <= numRows; row++) {
            int midCol = (numCols + 1) / 2;
            // Check middle column first
            if (!seats.get(generateSeatKey(row, midCol))) {
                return new int[]{row, midCol};
            }
            // Check alternating columns around middle
            for (int offset = 1; offset <= numCols / 2; offset++) {
                if (midCol - offset >= 1 && !seats.get(generateSeatKey(row, midCol - offset))) {
                    return new int[]{row, midCol - offset};
                }
                if (midCol + offset <= numCols && !seats.get(generateSeatKey(row, midCol + offset))) {
                    return new int[]{row, midCol + offset};
                }
            }
        }
        return null;
    }
}
