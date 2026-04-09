// Converted Java method
import java.util.regex.Pattern;

class UserValidator {

    /**
     * Validates user credentials and checks password strength.
     * 
     * @param user The user object to validate
     * @return A validation result message
     * @throws IllegalArgumentException if user is null
     */
    public String validateUser(User user) {
        if (user == null) {
            throw new IllegalArgumentException("User cannot be null");
        }

        StringBuilder validationMessage = new StringBuilder();

        // Validate login
        if (user.getLogin() == null || user.getLogin().trim().isEmpty()) {
            validationMessage.append("Login cannot be empty. ");
        } else if (user.getLogin().length() < 4 || user.getLogin().length() > 20) {
            validationMessage.append("Login must be between 4 and 20 characters. ");
        }

        // Validate password strength
        if (user.getPassword() == null || user.getPassword().trim().isEmpty()) {
            validationMessage.append("Password cannot be empty. ");
        } else {
            String password = user.getPassword();
            boolean hasLetter = Pattern.compile("[a-zA-Z]").matcher(password).find();
            boolean hasDigit = Pattern.compile("[0-9]").matcher(password).find();
            boolean hasSpecial = Pattern.compile("[^a-zA-Z0-9]").matcher(password).find();

            if (password.length() < 8) {
                validationMessage.append("Password must be at least 8 characters long. ");
            }
            if (!hasLetter) {
                validationMessage.append("Password must contain at least one letter. ");
            }
            if (!hasDigit) {
                validationMessage.append("Password must contain at least one digit. ");
            }
            if (!hasSpecial) {
                validationMessage.append("Password must contain at least one special character. ");
            }
        }

        // Validate role
        if (user.getRole() == null) {
            validationMessage.append("Role must be specified. ");
        }

        return validationMessage.length() == 0 ? "User is valid" : validationMessage.toString().trim();
    }
}

class User {
    private Integer userId;
    private String login;
    private String password;
    private Role role;

    public User(Integer userId, String login, String password, Role role) {
        this.userId = userId;
        this.login = login;
        this.password = password;
        this.role = role;
    }

    public Integer getUserId() {
        return userId;
    }

    public String getLogin() {
        return login;
    }

    public String getPassword() {
        return password;
    }

    public Role getRole() {
        return role;
    }
}

enum Role {
    ADMIN, CUSTOMER, GUEST
}
