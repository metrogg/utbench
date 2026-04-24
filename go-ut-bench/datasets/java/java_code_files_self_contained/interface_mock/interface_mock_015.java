// Converted Java method
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

class FinancialReportAnalyzer {
    
    /**
     * Analyzes financial data and generates a report summary.
     * This is a simplified version of the original financial report generation logic.
     * 
     * @param policies List of policy data (simplified from original database records)
     * @param startDate Start date of reporting period (format: "YYYY-MM-DD")
     * @param endDate End date of reporting period (format: "YYYY-MM-DD")
     * @return Map containing report summary with keys: "preRiskSummary", "curRiskSummary", 
     *         "totalPremiums", "policyCounts", "reportDetails"
     */
    public Map<String, Object> analyzeFinancialData(List<PolicyData> policies, String startDate, String endDate) {
        Map<String, Object> report = new HashMap<>();
        
        // Validate input dates
        if (startDate == null || endDate == null || startDate.isEmpty() || endDate.isEmpty()) {
            throw new IllegalArgumentException("Start date and end date cannot be empty");
        }
        
        // Initialize data structures
        Map<String, RiskSummary> preRiskMap = new HashMap<>();
        Map<String, RiskSummary> curRiskMap = new HashMap<>();
        double totalPremium = 0;
        int totalPolicies = 0;
        
        // Process each policy
        for (PolicyData policy : policies) {
            // Update previous risk summary
            preRiskMap.computeIfAbsent(policy.getPreviousRiskCode(), k -> new RiskSummary(policy.getPreviousRiskName()))
                     .addPolicy(policy.getPremium());
            
            // Update current risk summary
            curRiskMap.computeIfAbsent(policy.getCurrentRiskCode(), k -> new RiskSummary(policy.getCurrentRiskName()))
                     .addPolicy(policy.getPremium());
            
            totalPremium += policy.getPremium();
            totalPolicies++;
        }
        
        // Prepare report data
        List<Map<String, String>> preRiskSummary = new ArrayList<>();
        List<Map<String, String>> curRiskSummary = new ArrayList<>();
        
        // Convert risk maps to report format
        preRiskMap.forEach((code, summary) -> {
            Map<String, String> entry = new HashMap<>();
            entry.put("riskCode", code);
            entry.put("riskName", summary.getRiskName());
            entry.put("premium", String.format("%.2f", summary.getTotalPremium()));
            entry.put("policyCount", String.valueOf(summary.getPolicyCount()));
            preRiskSummary.add(entry);
        });
        
        curRiskMap.forEach((code, summary) -> {
            Map<String, String> entry = new HashMap<>();
            entry.put("riskCode", code);
            entry.put("riskName", summary.getRiskName());
            entry.put("premium", String.format("%.2f", summary.getTotalPremium()));
            entry.put("policyCount", String.valueOf(summary.getPolicyCount()));
            curRiskSummary.add(entry);
        });
        
        // Populate final report
        report.put("preRiskSummary", preRiskSummary);
        report.put("curRiskSummary", curRiskSummary);
        report.put("totalPremiums", String.format("%.2f", totalPremium));
        report.put("policyCounts", totalPolicies);
        report.put("reportPeriod", startDate + " to " + endDate);
        
        return report;
    }
    
    // Helper class for risk summary
    private static class RiskSummary {
        private String riskName;
        private double totalPremium;
        private int policyCount;
        
        public RiskSummary(String riskName) {
            this.riskName = riskName;
        }
        
        public void addPolicy(double premium) {
            this.totalPremium += premium;
            this.policyCount++;
        }
        
        public String getRiskName() { return riskName; }
        public double getTotalPremium() { return totalPremium; }
        public int getPolicyCount() { return policyCount; }
    }
    
    // Simplified policy data structure
    public static class PolicyData {
        private String previousRiskCode;
        private String previousRiskName;
        private String currentRiskCode;
        private String currentRiskName;
        private double premium;
        
        public PolicyData(String prevCode, String prevName, String currCode, String currName, double premium) {
            this.previousRiskCode = prevCode;
            this.previousRiskName = prevName;
            this.currentRiskCode = currCode;
            this.currentRiskName = currName;
            this.premium = premium;
        }
        
        // Getters
        public String getPreviousRiskCode() { return previousRiskCode; }
        public String getPreviousRiskName() { return previousRiskName; }
        public String getCurrentRiskCode() { return currentRiskCode; }
        public String getCurrentRiskName() { return currentRiskName; }
        public double getPremium() { return premium; }
    }
}
