import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.HashMap;
import java.util.stream.Collectors;

/**
 * Simulates user information management with enhanced features
 */
class UserInfoManager {
    
    /**
     * Processes phone contacts in batches with enhanced error handling and statistics
     * 
     * @param contacts List of phone contacts to process
     * @return Map containing processed users and processing statistics
     * @throws IllegalArgumentException if contacts list is null
     */
    public Map<String, Object> processContacts(List<PhoneContact> contacts) {
        if (contacts == null) {
            throw new IllegalArgumentException("Contacts list cannot be null");
        }

        Map<String, Object> result = new HashMap<>();
        List<UserInfo> processedUsers = new ArrayList<>();
        int batchSize = 200;
        int totalBatches = (int) Math.ceil((double) contacts.size() / batchSize);
        int failedProcesses = 0;
        int duplicateContacts = 0;

        // Track seen phone numbers to identify duplicates
        Map<String, Boolean> seenNumbers = new HashMap<>();

        for (int i = 0; i < totalBatches; i++) {
            int fromIndex = i * batchSize;
            int toIndex = Math.min(fromIndex + batchSize, contacts.size());
            List<PhoneContact> batch = contacts.subList(fromIndex, toIndex);

            try {
                List<UserInfo> batchResults = processBatch(batch);
                
                // Filter out duplicates within batch
                for (UserInfo user : batchResults) {
                    if (user != null) {
                        if (seenNumbers.containsKey(user.getPhoneNumber())) {
                            duplicateContacts++;
                        } else {
                            seenNumbers.put(user.getPhoneNumber(), true);
                            processedUsers.add(user);
                        }
                    } else {
                        failedProcesses++;
                    }
                }
            } catch (Exception e) {
                failedProcesses += batch.size();
            }
        }

        result.put("users", processedUsers);
        result.put("totalContacts", contacts.size());
        result.put("processedContacts", processedUsers.size());
        result.put("duplicateContacts", duplicateContacts);
        result.put("failedProcesses", failedProcesses);
        result.put("batchesProcessed", totalBatches);

        return result;
    }

    /**
     * Processes a single batch of contacts
     */
    private List<UserInfo> processBatch(List<PhoneContact> batch) {
        // Simulate API call and processing
        return batch.stream()
                .map(contact -> {
                    try {
                        // Simulate potential processing failures
                        if (contact.getPhoneNumber().isEmpty()) {
                            return null;
                        }
                        return new UserInfo(
                            contact.getName(),
                            contact.getPhoneNumber(),
                            "user_" + contact.getPhoneNumber() + "@example.com"
                        );
                    } catch (Exception e) {
                        return null;
                    }
                })
                .collect(Collectors.toList());
    }

    // Helper classes
    public static class PhoneContact {
        private String name;
        private String phoneNumber;

        public PhoneContact(String name, String phoneNumber) {
            this.name = name;
            this.phoneNumber = phoneNumber;
        }

        public String getName() { return name; }
        public String getPhoneNumber() { return phoneNumber; }
    }

    public static class UserInfo {
        private String name;
        private String phoneNumber;
        private String email;

        public UserInfo(String name, String phoneNumber, String email) {
            this.name = name;
            this.phoneNumber = phoneNumber;
            this.email = email;
        }

        public String getName() { return name; }
        public String getPhoneNumber() { return phoneNumber; }
        public String getEmail() { return email; }
    }
}
