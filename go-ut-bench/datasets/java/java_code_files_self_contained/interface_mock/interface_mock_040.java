import java.util.ArrayList;
import java.util.List;

/**
 * Represents a Reddit feed item with various properties.
 */
class RedditFeedItem {
    private String title;
    private String url;
    private String subreddit;
    private String domain;
    private String id;
    private String author;
    private String score;
    private String thumbnail;
    private String permalink;
    private String created;
    private String numComments;

    public RedditFeedItem(String title, String url, String subreddit, String domain, String id,
                         String author, String score, String thumbnail, String permalink,
                         String created, String numComments) {
        this.title = title;
        this.url = url;
        this.subreddit = subreddit;
        this.domain = domain;
        this.id = id;
        this.author = author;
        this.score = score;
        this.thumbnail = thumbnail;
        this.permalink = permalink;
        this.created = created;
        this.numComments = numComments;
    }

    // Getters
    public String getTitle() { return title; }
    public String getUrl() { return url; }
    public String getSubreddit() { return subreddit; }
    public String getDomain() { return domain; }
    public String getId() { return id; }
    public String getAuthor() { return author; }
    public String getScore() { return score; }
    public String getThumbnail() { return thumbnail; }
    public String getPermalink() { return permalink; }
    public String getCreated() { return created; }
    public String getNumComments() { return numComments; }
}

/**
 * Processes and analyzes Reddit feed data.
 */
class RedditFeedAnalyzer {
    
    /**
     * Filters feed items by minimum score and returns a formatted summary.
     * 
     * @param feedItems List of Reddit feed items
     * @param minScore Minimum score threshold for filtering
     * @return Formatted summary of filtered items
     */
    public String analyzeAndFormatFeed(List<RedditFeedItem> feedItems, int minScore) {
        if (feedItems == null || feedItems.isEmpty()) {
            return "No feed items to analyze";
        }

        List<RedditFeedItem> filteredItems = new ArrayList<>();
        int totalScore = 0;
        int totalComments = 0;

        // Filter items and calculate totals
        for (RedditFeedItem item : feedItems) {
            try {
                int itemScore = Integer.parseInt(item.getScore());
                if (itemScore >= minScore) {
                    filteredItems.add(item);
                    totalScore += itemScore;
                    totalComments += Integer.parseInt(item.getNumComments());
                }
            } catch (NumberFormatException e) {
                // Skip items with invalid score or comment numbers
                continue;
            }
        }

        if (filteredItems.isEmpty()) {
            return "No items meet the minimum score requirement of " + minScore;
        }

        // Calculate averages
        double avgScore = (double) totalScore / filteredItems.size();
        double avgComments = (double) totalComments / filteredItems.size();

        // Build summary
        StringBuilder summary = new StringBuilder();
        summary.append("Feed Analysis Summary:\n");
        summary.append("----------------------\n");
        summary.append("Total items: ").append(filteredItems.size()).append("\n");
        summary.append("Average score: ").append(String.format("%.2f", avgScore)).append("\n");
        summary.append("Average comments: ").append(String.format("%.2f", avgComments)).append("\n");
        summary.append("Top 3 items by score:\n");

        // Get top 3 items by score
        filteredItems.sort((a, b) -> {
            try {
                return Integer.compare(Integer.parseInt(b.getScore()), Integer.parseInt(a.getScore()));
            } catch (NumberFormatException e) {
                return 0;
            }
        });

        for (int i = 0; i < Math.min(3, filteredItems.size()); i++) {
            RedditFeedItem item = filteredItems.get(i);
            summary.append(i + 1).append(". ")
                  .append(item.getTitle()).append(" (Score: ")
                  .append(item.getScore()).append(", Comments: ")
                  .append(item.getNumComments()).append(")\n");
        }

        return summary.toString();
    }
}
