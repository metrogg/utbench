import java.util.ArrayList;
import java.util.List;
import java.util.Objects;

class FrequencyPlanner {
    /**
     * Represents a frequency assignment with bandwidth in kHz.
     */
    public static class FrequencyAssignment {
        private final int bandwidthInkHz;

        public FrequencyAssignment(int bandwidthInkHz) {
            if (bandwidthInkHz <= 0) {
                throw new IllegalArgumentException("Bandwidth must be positive");
            }
            this.bandwidthInkHz = bandwidthInkHz;
        }

        public int getBandwidthInkHz() {
            return bandwidthInkHz;
        }

        @Override
        public boolean equals(Object o) {
            if (this == o) return true;
            if (o == null || getClass() != o.getClass()) return false;
            FrequencyAssignment that = (FrequencyAssignment) o;
            return bandwidthInkHz == that.bandwidthInkHz;
        }

        @Override
        public int hashCode() {
            return Objects.hash(bandwidthInkHz);
        }
    }

    /**
     * Represents a channel with name and available bandwidth in kHz.
     */
    public static class Channel {
        private final String name;
        private final int availableBandwidthInkHz;

        public Channel(String name, int availableBandwidthInkHz) {
            if (availableBandwidthInkHz <= 0) {
                throw new IllegalArgumentException("Bandwidth must be positive");
            }
            this.name = Objects.requireNonNull(name, "Channel name cannot be null");
            this.availableBandwidthInkHz = availableBandwidthInkHz;
        }

        public String getName() {
            return name;
        }

        public int getAvailableBandwidthInkHz() {
            return availableBandwidthInkHz;
        }

        @Override
        public boolean equals(Object o) {
            if (this == o) return true;
            if (o == null || getClass() != o.getClass()) return false;
            Channel channel = (Channel) o;
            return availableBandwidthInkHz == channel.availableBandwidthInkHz &&
                    name.equals(channel.name);
        }

        @Override
        public int hashCode() {
            return Objects.hash(name, availableBandwidthInkHz);
        }
    }

    /**
     * Represents a plan containing frequency assignments and channels.
     */
    public static class Plan {
        private final List<FrequencyAssignment> assignments;
        private final List<Channel> channels;

        public Plan(List<FrequencyAssignment> assignments, List<Channel> channels) {
            this.assignments = new ArrayList<>(Objects.requireNonNull(assignments));
            this.channels = new ArrayList<>(Objects.requireNonNull(channels));
        }

        public List<FrequencyAssignment> getAssignments() {
            return new ArrayList<>(assignments);
        }

        public List<Channel> getChannels() {
            return new ArrayList<>(channels);
        }

        /**
         * Calculates the total bandwidth required by all frequency assignments.
         * @return Total bandwidth in kHz
         */
        public int getTotalRequiredBandwidth() {
            return assignments.stream().mapToInt(FrequencyAssignment::getBandwidthInkHz).sum();
        }

        /**
         * Calculates the total available bandwidth from all channels.
         * @return Total available bandwidth in kHz
         */
        public int getTotalAvailableBandwidth() {
            return channels.stream().mapToInt(Channel::getAvailableBandwidthInkHz).sum();
        }

        /**
         * Checks if the plan is feasible (total available bandwidth >= total required bandwidth).
         * @return true if feasible, false otherwise
         */
        public boolean isFeasible() {
            return getTotalAvailableBandwidth() >= getTotalRequiredBandwidth();
        }

        /**
         * Finds the channel with the maximum available bandwidth.
         * @return Channel with maximum bandwidth or null if no channels exist
         */
        public Channel findChannelWithMaxBandwidth() {
            return channels.stream()
                    .max((c1, c2) -> Integer.compare(c1.getAvailableBandwidthInkHz(), 
                                                    c2.getAvailableBandwidthInkHz()))
                    .orElse(null);
        }

        @Override
        public boolean equals(Object o) {
            if (this == o) return true;
            if (o == null || getClass() != o.getClass()) return false;
            Plan plan = (Plan) o;
            return assignments.equals(plan.assignments) &&
                    channels.equals(plan.channels);
        }

        @Override
        public int hashCode() {
            return Objects.hash(assignments, channels);
        }
    }

    /**
     * Creates a new plan from hardcoded data (simulating CSV input).
     * @param assignmentBandwidths Array of bandwidths for frequency assignments
     * @param channelData Array of channel data (name and bandwidth pairs)
     * @return New Plan object
     */
    public static Plan createPlan(int[] assignmentBandwidths, String[][] channelData) {
        List<FrequencyAssignment> assignments = new ArrayList<>();
        List<Channel> channels = new ArrayList<>();

        for (int bandwidth : assignmentBandwidths) {
            if (bandwidth > 0) {
                assignments.add(new FrequencyAssignment(bandwidth));
            }
        }

        for (String[] channel : channelData) {
            if (channel.length == 2) {
                try {
                    int bandwidth = Integer.parseInt(channel[1]);
                    if (bandwidth > 0) {
                        channels.add(new Channel(channel[0], bandwidth));
                    }
                } catch (NumberFormatException e) {
                    // Skip invalid entries
                }
            }
        }

        return new Plan(assignments, channels);
    }
}
