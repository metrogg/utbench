// Converted Java method
import java.util.ArrayList;
import java.util.List;

class AirplaneManager {
    private List<Airplane> airplanes = new ArrayList<>();

    /**
     * Adds a new airplane to the manager with validation.
     * @param model The airplane model name
     * @param travelClass The travel class configuration
     * @throws IllegalArgumentException if model is empty or travelClass is null
     */
    public void addAirplane(String model, TravelClass travelClass) {
        if (model == null || model.trim().isEmpty()) {
            throw new IllegalArgumentException("Airplane model cannot be empty");
        }
        if (travelClass == null) {
            throw new IllegalArgumentException("Travel class cannot be null");
        }
        airplanes.add(new Airplane(model, travelClass));
    }

    /**
     * Finds all airplanes that have at least the specified number of total seats
     * @param minSeats The minimum number of seats required
     * @return List of matching airplanes
     */
    public List<Airplane> findAirplanesByMinSeats(int minSeats) {
        List<Airplane> result = new ArrayList<>();
        for (Airplane airplane : airplanes) {
            if (airplane.getTravelClass().getTotalSeats() >= minSeats) {
                result.add(airplane);
            }
        }
        return result;
    }

    /**
     * Gets the airplane with the most total seats
     * @return The airplane with maximum seats or null if no airplanes exist
     */
    public Airplane getLargestAirplane() {
        if (airplanes.isEmpty()) {
            return null;
        }
        Airplane largest = airplanes.get(0);
        for (Airplane airplane : airplanes) {
            if (airplane.getTravelClass().getTotalSeats() > largest.getTravelClass().getTotalSeats()) {
                largest = airplane;
            }
        }
        return largest;
    }
}

class Airplane {
    private String model;
    private TravelClass travelClass;

    public Airplane(String model, TravelClass travelClass) {
        this.model = model;
        this.travelClass = travelClass;
    }

    public String getModel() {
        return model;
    }

    public TravelClass getTravelClass() {
        return travelClass;
    }

    @Override
    public String toString() {
        return "Airplane{" +
                "model='" + model + '\'' +
                ", totalSeats=" + travelClass.getTotalSeats() +
                '}';
    }
}

class TravelClass {
    private int economySeats;
    private int businessSeats;
    private int firstClassSeats;

    public TravelClass(int economySeats, int businessSeats, int firstClassSeats) {
        this.economySeats = economySeats;
        this.businessSeats = businessSeats;
        this.firstClassSeats = firstClassSeats;
    }

    public int getTotalSeats() {
        return economySeats + businessSeats + firstClassSeats;
    }
}
