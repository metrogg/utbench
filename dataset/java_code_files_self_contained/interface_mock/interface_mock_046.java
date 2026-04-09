import java.util.Date;
import java.util.concurrent.atomic.AtomicInteger;

/**
 * Enhanced factory simulation with configurable workshops and stocks.
 * This version adds more complex operations and statistics tracking.
 */
class EnhancedFactory {
    private final Stock inputStock;
    private final Stock intermediateStock;
    private final Stock outputStock;
    private final Workshop[] workshops;
    private final int version;
    private final AtomicInteger totalProcessedItems = new AtomicInteger(0);
    private final AtomicInteger processingErrors = new AtomicInteger(0);

    /**
     * Creates an enhanced factory with configurable parameters
     * @param version Factory configuration version (1-4)
     * @param initialInput Initial items in input stock
     * @param intermediateCapacity Capacity of intermediate stock (if used)
     */
    public EnhancedFactory(int version, int initialInput, int intermediateCapacity) {
        this.version = version;
        this.inputStock = new Stock("input", initialInput);
        
        if (version >= 2) {
            this.intermediateStock = new Stock("intermediate", 0, intermediateCapacity);
        } else {
            this.intermediateStock = new Stock("intermediate", 0);
        }
        
        this.outputStock = new Stock("output", 0);
        
        switch (version) {
            case 1:
                workshops = new Workshop[]{
                    new Workshop(inputStock, outputStock, 5, "W1"),
                    new Workshop(inputStock, outputStock, 5, "W2")
                };
                break;
            case 2:
                workshops = new Workshop[]{
                    new Workshop(inputStock, intermediateStock, 10, "W1"),
                    new Workshop(intermediateStock, outputStock, 10, "W2")
                };
                break;
            case 3:
                workshops = new Workshop[]{
                    new Workshop(inputStock, intermediateStock, 10, "W1"),
                    new Workshop(intermediateStock, outputStock, 5, "W2"),
                    new Workshop(intermediateStock, outputStock, 5, "W3")
                };
                break;
            case 4:
                workshops = new Workshop[]{
                    new Workshop(inputStock, intermediateStock, 5, "W1"),
                    new Workshop(inputStock, intermediateStock, 5, "W2"),
                    new Workshop(intermediateStock, outputStock, 5, "W3"),
                    new Workshop(intermediateStock, outputStock, 5, "W4")
                };
                break;
            default:
                throw new IllegalArgumentException("Invalid factory version: " + version);
        }
    }

    /**
     * Runs the factory simulation and returns performance statistics
     * @return FactoryPerformance object containing execution metrics
     */
    public FactoryPerformance runSimulation() {
        long startTime = System.currentTimeMillis();
        
        // Start all workshops
        for (Workshop workshop : workshops) {
            workshop.setTimeStamp(startTime);
            workshop.start();
        }

        // Wait for all workshops to complete
        for (Workshop workshop : workshops) {
            try {
                workshop.join();
                totalProcessedItems.addAndGet(workshop.getProcessedCount());
                processingErrors.addAndGet(workshop.getErrorCount());
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
                processingErrors.incrementAndGet();
            }
        }

        long endTime = System.currentTimeMillis();
        long duration = endTime - startTime;

        return new FactoryPerformance(
            version,
            duration,
            totalProcessedItems.get(),
            processingErrors.get(),
            inputStock.getCurrentCount(),
            outputStock.getCurrentCount(),
            intermediateStock.getCurrentCount()
        );
    }
}

/**
 * Represents factory performance metrics
 */
class FactoryPerformance {
    public final int version;
    public final long durationMs;
    public final int totalProcessed;
    public final int errors;
    public final int remainingInput;
    public final int outputProduced;
    public final int intermediateRemaining;

    public FactoryPerformance(int version, long durationMs, int totalProcessed, 
                            int errors, int remainingInput, int outputProduced, 
                            int intermediateRemaining) {
        this.version = version;
        this.durationMs = durationMs;
        this.totalProcessed = totalProcessed;
        this.errors = errors;
        this.remainingInput = remainingInput;
        this.outputProduced = outputProduced;
        this.intermediateRemaining = intermediateRemaining;
    }
}

/**
 * Simplified Workshop implementation for the benchmark
 */
class Workshop extends Thread {
    private final Stock source;
    private final Stock destination;
    private final int itemsToProcess;
    private final String name;
    private long timeStamp;
    private int processedCount = 0;
    private int errorCount = 0;

    public Workshop(Stock source, Stock destination, int itemsToProcess, String name) {
        this.source = source;
        this.destination = destination;
        this.itemsToProcess = itemsToProcess;
        this.name = name;
    }

    public void setTimeStamp(long timeStamp) {
        this.timeStamp = timeStamp;
    }

    public int getProcessedCount() {
        return processedCount;
    }

    public int getErrorCount() {
        return errorCount;
    }

    @Override
    public void run() {
        long startTime = System.currentTimeMillis();
        
        for (int i = 0; i < itemsToProcess; i++) {
            try {
                if (source.takeItem()) {
                    destination.addItem();
                    processedCount++;
                }
            } catch (Exception e) {
                errorCount++;
            }
        }
        
        this.timeStamp = System.currentTimeMillis() - startTime;
    }
}

/**
 * Simplified Stock implementation for the benchmark
 */
class Stock {
    private final String name;
    private int count;
    private final int capacity;

    public Stock(String name, int initialCount) {
        this(name, initialCount, Integer.MAX_VALUE);
    }

    public Stock(String name, int initialCount, int capacity) {
        this.name = name;
        this.count = initialCount;
        this.capacity = capacity;
    }

    public synchronized boolean takeItem() {
        if (count > 0) {
            count--;
            return true;
        }
        return false;
    }

    public synchronized boolean addItem() {
        if (count < capacity) {
            count++;
            return true;
        }
        return false;
    }

    public int getCurrentCount() {
        return count;
    }
}
