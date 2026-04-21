// Converted Java method
import java.io.Serializable;
import java.util.HashMap;
import java.util.Map;

class AccidentNatureInfo implements Serializable {
    private String natureCode;
    private String description;
    private int severityLevel;
    private static final Map<String, String> NATURE_DESCRIPTIONS = new HashMap<>();
    
    static {
        NATURE_DESCRIPTIONS.put("FIRE", "Fire-related accident");
        NATURE_DESCRIPTIONS.put("CHEM", "Chemical exposure or spill");
        NATURE_DESCRIPTIONS.put("MECH", "Mechanical failure");
        NATURE_DESCRIPTIONS.put("ELEC", "Electrical incident");
        NATURE_DESCRIPTIONS.put("FALL", "Slip, trip or fall");
    }

    public AccidentNatureInfo() {
        this("UNKN", "Unknown nature", 1);
    }

    public AccidentNatureInfo(String natureCode, String description, int severityLevel) {
        if (natureCode == null || natureCode.trim().isEmpty()) {
            throw new IllegalArgumentException("Nature code cannot be null or empty");
        }
        if (severityLevel < 1 || severityLevel > 5) {
            throw new IllegalArgumentException("Severity level must be between 1 and 5");
        }
        
        this.natureCode = natureCode.toUpperCase();
        this.description = description;
        this.severityLevel = severityLevel;
    }

    public String getNatureCode() {
        return natureCode;
    }

    public String getDescription() {
        return NATURE_DESCRIPTIONS.getOrDefault(natureCode, description);
    }

    public int getSeverityLevel() {
        return severityLevel;
    }

    public void setSeverityLevel(int severityLevel) {
        if (severityLevel < 1 || severityLevel > 5) {
            throw new IllegalArgumentException("Severity level must be between 1 and 5");
        }
        this.severityLevel = severityLevel;
    }

    public boolean isHighRisk() {
        return severityLevel >= 4;
    }

    public static boolean isValidNatureCode(String code) {
        return NATURE_DESCRIPTIONS.containsKey(code.toUpperCase());
    }

    @Override
    public String toString() {
        return String.format("AccidentNatureInfo{code=%s, description='%s', severity=%d}",
                natureCode, getDescription(), severityLevel);
    }
}
