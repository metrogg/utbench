import java.util.ArrayList;
import java.util.List;

/**
 * Represents a company with hierarchical employee structure.
 * Uses Singleton pattern to ensure only one company instance exists.
 */
class Company {
    private static Company instance;
    private Manager ceo;
    private List<Employee> allEmployees;
    private String companyName;
    private int employeeCounter;

    private Company() {
        this.allEmployees = new ArrayList<>();
        this.employeeCounter = 1;
    }

    /**
     * Gets the singleton instance of the Company
     * @return Company instance
     */
    public static Company getInstance() {
        if (instance == null) {
            instance = new Company();
        }
        return instance;
    }

    /**
     * Sets the company name
     * @param name Name of the company
     */
    public void setCompanyName(String name) {
        this.companyName = name;
    }

    /**
     * Hires a CEO for the company
     * @param m Manager to be hired as CEO
     * @throws IllegalStateException if CEO is already set
     */
    public void hireCEO(Manager m) {
        if (ceo != null) {
            throw new IllegalStateException("CEO already exists: " + ceo.getName());
        }
        ceo = m;
        m.setId(generateEmployeeId());
        allEmployees.add(m);
    }

    /**
     * Hires an employee under a manager
     * @param manager The hiring manager
     * @param employee The employee to be hired
     * @throws IllegalArgumentException if manager is not part of company
     */
    public void hireEmployee(Manager manager, Employee employee) {
        if (!allEmployees.contains(manager)) {
            throw new IllegalArgumentException("Manager not found in company");
        }
        employee.setId(generateEmployeeId());
        manager.addSubordinate(employee);
        allEmployees.add(employee);
    }

    /**
     * Generates a unique employee ID
     * @return Generated employee ID
     */
    private String generateEmployeeId() {
        return "EMP-" + (employeeCounter++);
    }

    /**
     * Gets the CEO of the company
     * @return CEO Manager object
     */
    public Manager getCEO() {
        return ceo;
    }

    /**
     * Gets all employees in the company
     * @return List of all employees
     */
    public List<Employee> getAllEmployees() {
        return new ArrayList<>(allEmployees);
    }

    /**
     * Gets the total number of employees
     * @return Employee count
     */
    public int getEmployeeCount() {
        return allEmployees.size();
    }

    /**
     * Gets the company name
     * @return Company name
     */
    public String getCompanyName() {
        return companyName;
    }
}

/**
 * Base class for all employees
 */
class Employee {
    private String id;
    private String name;
    private double salary;

    public Employee(String name, double salary) {
        this.name = name;
        this.salary = salary;
    }

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getName() {
        return name;
    }

    public double getSalary() {
        return salary;
    }

    public void setSalary(double salary) {
        this.salary = salary;
    }

    @Override
    public String toString() {
        return "Employee{" + "id='" + id + '\'' + ", name='" + name + '\'' + '}';
    }
}

/**
 * Manager class that can have subordinates
 */
class Manager extends Employee {
    private List<Employee> subordinates;

    public Manager(String name, double salary) {
        super(name, salary);
        this.subordinates = new ArrayList<>();
    }

    public void addSubordinate(Employee employee) {
        subordinates.add(employee);
    }

    public List<Employee> getSubordinates() {
        return new ArrayList<>(subordinates);
    }

    public int getTeamSize() {
        return subordinates.size();
    }
}
