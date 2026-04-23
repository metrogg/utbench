import java.util.LinkedList;
import java.util.HashMap;
import java.util.HashSet;
import java.util.Set;

class DFA {
    private LinkedList<State> states;
    private State initialState;
    private Set<State> finalStates;
    private Set<Character> alphabet;

    public DFA() {
        this.states = new LinkedList<>();
        this.finalStates = new HashSet<>();
        this.alphabet = new HashSet<>();
    }

    /**
     * Adds a state to the DFA
     * @param state The state to add
     * @param isFinal Whether the state is a final/accepting state
     */
    public void addState(State state, boolean isFinal) {
        if (states.isEmpty()) {
            initialState = state;
        }
        states.add(state);
        if (isFinal) {
            finalStates.add(state);
        }
    }

    /**
     * Adds a transition between states
     * @param from The source state
     * @param to The destination state
     * @param symbol The transition symbol
     */
    public void addTransition(State from, State to, char symbol) {
        from.addTransition(symbol, to);
        alphabet.add(symbol);
    }

    /**
     * Checks if a string is accepted by the DFA
     * @param input The string to check
     * @return true if the string is accepted, false otherwise
     */
    public boolean accepts(String input) {
        State current = initialState;
        for (char c : input.toCharArray()) {
            current = current.getNextState(c);
            if (current == null) {
                return false;
            }
        }
        return finalStates.contains(current);
    }

    /**
     * Minimizes the DFA using Hopcroft's algorithm
     */
    public void minimize() {
        // Implementation of Hopcroft's algorithm would go here
        // This is a placeholder for the actual minimization logic
    }

    @Override
    public String toString() {
        StringBuilder sb = new StringBuilder();
        sb.append("DFA with ").append(states.size()).append(" states\n");
        sb.append("Alphabet: ").append(alphabet).append("\n");
        sb.append("Initial state: ").append(initialState.getName()).append("\n");
        sb.append("Final states: ");
        for (State s : finalStates) {
            sb.append(s.getName()).append(" ");
        }
        sb.append("\nTransitions:\n");
        for (State s : states) {
            sb.append(s.toString()).append("\n");
        }
        return sb.toString();
    }
}

class State {
    private String name;
    private HashMap<Character, State> transitions;

    public State(String name) {
        this.name = name;
        this.transitions = new HashMap<>();
    }

    public String getName() {
        return name;
    }

    public void addTransition(char symbol, State nextState) {
        transitions.put(symbol, nextState);
    }

    public State getNextState(char symbol) {
        return transitions.get(symbol);
    }

    @Override
    public String toString() {
        StringBuilder sb = new StringBuilder();
        sb.append(name).append(": ");
        for (HashMap.Entry<Character, State> entry : transitions.entrySet()) {
            sb.append("--").append(entry.getKey()).append("-> ").append(entry.getValue().getName()).append(" ");
        }
        return sb.toString();
    }
}
