// Converted Java method
import java.util.*;
import java.util.stream.Collectors;

class ClientManager {
    private Map<Long, ClientTO> clients = new HashMap<>();
    private long nextId = 1;

    /**
     * Creates a new client with validation checks.
     * 
     * @param client the client to create (without ID)
     * @return the created client with assigned ID
     * @throws IllegalArgumentException if client is null, has an ID, or validation fails
     */
    public ClientTO createClient(ClientTO client) {
        if (client == null) {
            throw new IllegalArgumentException("Client cannot be null");
        }
        if (client.getId() != null) {
            throw new IllegalArgumentException("New client cannot have an existing ID");
        }
        
        validateClientAttributes(client);
        
        ClientTO newClient = new ClientTO(
            nextId++, 
            client.getFirstName(), 
            client.getLastName(), 
            client.getEmail(), 
            client.getPhone()
        );
        clients.put(newClient.getId(), newClient);
        return newClient;
    }

    /**
     * Finds a client by ID with optional authorization check.
     * 
     * @param id the client ID to search for
     * @param requireAdmin whether to enforce admin privileges
     * @return the found client or null if not found
     * @throws IllegalArgumentException if id is null
     * @throws SecurityException if requireAdmin is true and user doesn't have admin privileges
     */
    public ClientTO findClient(Long id, boolean requireAdmin) {
        if (id == null) {
            throw new IllegalArgumentException("ID cannot be null");
        }
        if (requireAdmin && !hasAdminPrivileges()) {
            throw new SecurityException("Admin privileges required");
        }
        return clients.get(id);
    }

    /**
     * Updates an existing client with validation checks.
     * 
     * @param client the client to update
     * @return the updated client
     * @throws IllegalArgumentException if client is null or doesn't exist
     */
    public ClientTO updateClient(ClientTO client) {
        if (client == null) {
            throw new IllegalArgumentException("Client cannot be null");
        }
        if (!clients.containsKey(client.getId())) {
            throw new IllegalArgumentException("Client does not exist");
        }
        
        validateClientAttributes(client);
        
        clients.put(client.getId(), client);
        return client;
    }

    /**
     * Deletes a client with optional authorization check.
     * 
     * @param id the ID of the client to delete
     * @param requireAdmin whether to enforce admin privileges
     * @return true if client was deleted, false if not found
     * @throws IllegalArgumentException if id is null
     * @throws SecurityException if requireAdmin is true and user doesn't have admin privileges
     */
    public boolean deleteClient(Long id, boolean requireAdmin) {
        if (id == null) {
            throw new IllegalArgumentException("ID cannot be null");
        }
        if (requireAdmin && !hasAdminPrivileges()) {
            throw new SecurityException("Admin privileges required");
        }
        return clients.remove(id) != null;
    }

    /**
     * Finds all clients with optional filtering.
     * 
     * @param filter an optional filter string to match against client names or emails
     * @param requireAdmin whether to enforce admin privileges
     * @return list of matching clients or all clients if no filter
     * @throws SecurityException if requireAdmin is true and user doesn't have admin privileges
     */
    public List<ClientTO> findAllClients(String filter, boolean requireAdmin) {
        if (requireAdmin && !hasAdminPrivileges()) {
            throw new SecurityException("Admin privileges required");
        }
        
        if (filter == null || filter.isEmpty()) {
            return new ArrayList<>(clients.values());
        }
        
        String lowerFilter = filter.toLowerCase();
        return clients.values().stream()
            .filter(c -> c.getFirstName().toLowerCase().contains(lowerFilter) ||
                        c.getLastName().toLowerCase().contains(lowerFilter) ||
                        c.getEmail().toLowerCase().contains(lowerFilter))
            .collect(Collectors.toList());
    }

    /**
     * Finds clients by name (first or last name).
     * 
     * @param name the name to search for
     * @return list of matching clients
     * @throws IllegalArgumentException if name is null or empty
     */
    public List<ClientTO> findClientsByName(String name) {
        if (name == null || name.isEmpty()) {
            throw new IllegalArgumentException("Name cannot be null or empty");
        }
        
        String lowerName = name.toLowerCase();
        return clients.values().stream()
            .filter(c -> c.getFirstName().toLowerCase().contains(lowerName) ||
                        c.getLastName().toLowerCase().contains(lowerName))
            .collect(Collectors.toList());
    }

    private void validateClientAttributes(ClientTO client) {
        if (client.getFirstName() == null || client.getFirstName().isEmpty()) {
            throw new IllegalArgumentException("First name is required");
        }
        if (client.getLastName() == null || client.getLastName().isEmpty()) {
            throw new IllegalArgumentException("Last name is required");
        }
        if (client.getEmail() == null || !client.getEmail().contains("@")) {
            throw new IllegalArgumentException("Valid email is required");
        }
    }

    private boolean hasAdminPrivileges() {
        // In a real implementation, this would check the security context
        return false; // Simplified for this example
    }
}

class ClientTO {
    private Long id;
    private String firstName;
    private String lastName;
    private String email;
    private String phone;

    public ClientTO(Long id, String firstName, String lastName, String email, String phone) {
        this.id = id;
        this.firstName = firstName;
        this.lastName = lastName;
        this.email = email;
        this.phone = phone;
    }

    // Getters and setters omitted for brevity
    public Long getId() { return id; }
    public String getFirstName() { return firstName; }
    public String getLastName() { return lastName; }
    public String getEmail() { return email; }
    public String getPhone() { return phone; }
}
