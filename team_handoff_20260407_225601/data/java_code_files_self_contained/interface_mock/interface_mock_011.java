// Converted Java method
import java.util.Set;
import java.util.HashSet;
import java.util.HashMap;
import java.util.Map;

class EnhancedTariffManager {
    private Map<Integer, Tariff> tariffs;
    private Map<Integer, Set<Contract>> tariffContracts;

    public EnhancedTariffManager() {
        this.tariffs = new HashMap<>();
        this.tariffContracts = new HashMap<>();
    }

    /**
     * Adds a new tariff to the system
     * @param tariff The tariff to add
     * @throws IllegalArgumentException if tariff is null or already exists
     */
    public void addTariff(Tariff tariff) {
        if (tariff == null) {
            throw new IllegalArgumentException("Tariff cannot be null");
        }
        if (tariffs.containsKey(tariff.getId())) {
            throw new IllegalArgumentException("Tariff with ID " + tariff.getId() + " already exists");
        }
        tariffs.put(tariff.getId(), tariff);
        tariffContracts.put(tariff.getId(), new HashSet<>());
    }

    /**
     * Gets all contracts subscribed to a specific tariff
     * @param tariffId The ID of the tariff
     * @return Set of contracts using this tariff
     * @throws IllegalArgumentException if tariff doesn't exist
     */
    public Set<Contract> getContractsByTariff(Integer tariffId) {
        if (!tariffs.containsKey(tariffId)) {
            throw new IllegalArgumentException("Tariff with ID " + tariffId + " doesn't exist");
        }
        return new HashSet<>(tariffContracts.get(tariffId));
    }

    /**
     * Subscribes a contract to a tariff
     * @param contract The contract to subscribe
     * @param tariffId The tariff ID to subscribe to
     * @throws IllegalArgumentException if contract or tariff doesn't exist
     */
    public void subscribeContractToTariff(Contract contract, Integer tariffId) {
        if (contract == null) {
            throw new IllegalArgumentException("Contract cannot be null");
        }
        if (!tariffs.containsKey(tariffId)) {
            throw new IllegalArgumentException("Tariff with ID " + tariffId + " doesn't exist");
        }
        tariffContracts.get(tariffId).add(contract);
    }

    /**
     * Gets the most popular tariff (with most contracts)
     * @return The most popular tariff or null if no tariffs exist
     */
    public Tariff getMostPopularTariff() {
        if (tariffs.isEmpty()) {
            return null;
        }
        
        int maxContracts = -1;
        Tariff popularTariff = null;
        
        for (Map.Entry<Integer, Set<Contract>> entry : tariffContracts.entrySet()) {
            if (entry.getValue().size() > maxContracts) {
                maxContracts = entry.getValue().size();
                popularTariff = tariffs.get(entry.getKey());
            }
        }
        
        return popularTariff;
    }
}

class Tariff {
    private Integer id;
    private String name;
    private double monthlyFee;
    
    public Tariff(Integer id, String name, double monthlyFee) {
        this.id = id;
        this.name = name;
        this.monthlyFee = monthlyFee;
    }
    
    public Integer getId() { return id; }
    public String getName() { return name; }
    public double getMonthlyFee() { return monthlyFee; }
}

class Contract {
    private Integer id;
    private String clientName;
    
    public Contract(Integer id, String clientName) {
        this.id = id;
        this.clientName = clientName;
    }
    
    public Integer getId() { return id; }
    public String getClientName() { return clientName; }
    
    @Override
    public String toString() {
        return "Contract[id=" + id + ", client=" + clientName + "]";
    }
}
