import java.util.List;
import java.util.ArrayList;
import java.util.LinkedList;

class CelestialSimulator {
    
    /**
     * Represents a celestial body with position, velocity, acceleration, and mass
     */
    public static class Celestial {
        public final String name;
        public final double mass;
        public final Vector position;
        public final Vector velocity;
        public final Vector acceleration;
        
        public Celestial(String name, double mass, Vector position, Vector velocity, Vector acceleration) {
            this.name = name;
            this.mass = mass;
            this.position = position;
            this.velocity = velocity;
            this.acceleration = acceleration;
        }
        
        /**
         * Updates the position and velocity of this celestial body based on its acceleration
         * @param timeStep The time step for the simulation
         */
        public void update(double timeStep) {
            // Update velocity: v = v0 + a*t
            Vector newVelocity = Vector.add(velocity, Vector.scale(acceleration, timeStep));
            
            // Update position: x = x0 + v*t
            Vector newPosition = Vector.add(position, Vector.scale(newVelocity, timeStep));
            
            // Update the celestial body's state
            velocity.x = newVelocity.x;
            velocity.y = newVelocity.y;
            position.x = newPosition.x;
            position.y = newPosition.y;
        }
        
        /**
         * Calculates the gravitational force exerted by another celestial body
         * @param other The other celestial body
         * @return The force vector
         */
        public Vector calculateGravityForce(Celestial other) {
            final double G = 6.67430e-11; // Gravitational constant
            
            // Calculate distance between bodies
            double dx = other.position.x - position.x;
            double dy = other.position.y - position.y;
            double distance = Math.sqrt(dx*dx + dy*dy);
            
            // Avoid division by zero for extremely close bodies
            if (distance < 1e-10) {
                return Vector.make(0, 0);
            }
            
            // Calculate force magnitude: F = G*m1*m2/r^2
            double forceMagnitude = G * mass * other.mass / (distance * distance);
            
            // Calculate force direction (unit vector)
            double fx = forceMagnitude * dx / distance;
            double fy = forceMagnitude * dy / distance;
            
            return Vector.make(fx, fy);
        }
    }
    
    /**
     * Represents a 2D vector with x and y components
     */
    public static class Vector {
        public double x;
        public double y;
        
        public Vector(double x, double y) {
            this.x = x;
            this.y = y;
        }
        
        public static Vector make(double x, double y) {
            return new Vector(x, y);
        }
        
        public static Vector add(Vector a, Vector b) {
            return new Vector(a.x + b.x, a.y + b.y);
        }
        
        public static Vector subtract(Vector a, Vector b) {
            return new Vector(a.x - b.x, a.y - b.y);
        }
        
        public static Vector scale(Vector v, double scalar) {
            return new Vector(v.x * scalar, v.y * scalar);
        }
        
        public static Vector rotate90(Vector v) {
            return new Vector(-v.y, v.x);
        }
        
        public static Vector unit(Vector v) {
            double magnitude = Math.sqrt(v.x*v.x + v.y*v.y);
            if (magnitude < 1e-10) return make(0, 0);
            return new Vector(v.x/magnitude, v.y/magnitude);
        }
    }
    
    /**
     * Simulates the movement of celestial bodies over time using Newtonian physics
     * @param bodies List of celestial bodies
     * @param timeStep Time step for the simulation
     * @param steps Number of simulation steps to perform
     * @return List of states at each step
     */
    public static List<List<Celestial>> simulate(List<Celestial> bodies, double timeStep, int steps) {
        List<List<Celestial>> simulationResults = new ArrayList<>();
        
        for (int step = 0; step < steps; step++) {
            // Calculate forces and update accelerations for all bodies
            for (Celestial body : bodies) {
                Vector totalForce = Vector.make(0, 0);
                
                // Sum forces from all other bodies
                for (Celestial other : bodies) {
                    if (body != other) {
                        Vector force = body.calculateGravityForce(other);
                        totalForce = Vector.add(totalForce, force);
                    }
                }
                
                // Update acceleration: a = F/m
                body.acceleration.x = totalForce.x / body.mass;
                body.acceleration.y = totalForce.y / body.mass;
            }
            
            // Update positions and velocities
            for (Celestial body : bodies) {
                body.update(timeStep);
            }
            
            // Save current state
            simulationResults.add(copyState(bodies));
        }
        
        return simulationResults;
    }
    
    /**
     * Creates a deep copy of the current state of all celestial bodies
     */
    private static List<Celestial> copyState(List<Celestial> bodies) {
        List<Celestial> copy = new ArrayList<>();
        for (Celestial body : bodies) {
            copy.add(new Celestial(
                body.name,
                body.mass,
                new Vector(body.position.x, body.position.y),
                new Vector(body.velocity.x, body.velocity.y),
                new Vector(body.acceleration.x, body.acceleration.y)
            ));
        }
        return copy;
    }
    
    /**
     * Creates a simple solar system with one sun and several planets
     */
    public static List<Celestial> createSolarSystem() {
        List<Celestial> system = new LinkedList<>();
        
        // Create sun at center with large mass
        Celestial sun = new Celestial(
            "sun", 
            1.989e30,  // Sun's mass in kg
            new Vector(0, 0), 
            new Vector(0, 0), 
            new Vector(0, 0)
        );
        system.add(sun);
        
        // Create some planets
        Celestial earth = new Celestial(
            "earth",
            5.972e24,
            new Vector(1.496e11, 0),  // ~1 AU from sun
            new Vector(0, 29.78e3),   // Orbital velocity in m/s
            new Vector(0, 0)
        );
        system.add(earth);
        
        Celestial mars = new Celestial(
            "mars",
            6.39e23,
            new Vector(2.279e11, 0),
            new Vector(0, 24.07e3),
            new Vector(0, 0)
        );
        system.add(mars);
        
        return system;
    }
}
