// Converted Java method
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

class ScenarioAnalyzer {
    /**
     * Analyzes a scenario tree structure and validates synchronous call pairs.
     * This method assigns indices to all nodes in the tree and identifies valid
     * synchronous request-response pairs.
     *
     * @param rootNode The root node of the scenario tree
     * @return A map containing:
     *         - "nodeIndices": Map of node names to their assigned indices
     *         - "syncPairs": Map of request indices to list of response indices
     *         - "warnings": List of warning messages for missing responses
     * @throws IllegalArgumentException if rootNode is null
     */
    public Map<String, Object> analyzeScenario(TaskNode rootNode) {
        if (rootNode == null) {
            throw new IllegalArgumentException("Root node cannot be null");
        }

        Map<String, Object> results = new HashMap<>();
        Map<String, Integer> nodeIndices = new HashMap<>();
        Map<Integer, List<Integer>> syncPairs = new HashMap<>();
        List<String> warnings = new ArrayList<>();

        // Synchronous Request Server Pair: To -> From
        Map<String, String> syncReqServerPair = new HashMap<>();
        Map<String, List<String>> serverResponseMap = new HashMap<>();

        List<TaskNode> currentLevel = new ArrayList<>();
        int idx = 0;
        rootNode.index = idx++;
        nodeIndices.put(rootNode.name, rootNode.index);
        currentLevel.add(rootNode);

        while (!currentLevel.isEmpty()) {
            List<TaskNode> nextLevel = new ArrayList<>();

            for (TaskNode node : currentLevel) {
                if (node.isSync) {
                    syncReqServerPair.put(node.serverName, node.parent.serverName);
                }

                for (TaskNode child : node.children) {
                    child.index = idx++;
                    nodeIndices.put(child.name, child.index);
                    nextLevel.add(child);

                    if (syncReqServerPair.containsKey(node.getServerName()) && 
                        syncReqServerPair.get(node.getServerName()).equals(child.serverName)) {
                        
                        if (serverResponseMap.containsKey(node.serverName)) {
                            serverResponseMap.get(node.serverName).add(child.serverName);
                            syncPairs.get(node.index).add(child.index);
                        } else {
                            List<String> serverList = new ArrayList<>();
                            serverList.add(child.serverName);
                            serverResponseMap.put(node.serverName, serverList);

                            List<Integer> indexList = new ArrayList<>();
                            indexList.add(child.index);
                            syncPairs.put(node.index, indexList);
                        }
                    }
                }
            }
            currentLevel = nextLevel;
        }

        // Check for missing responses
        for (String toServerName : syncReqServerPair.keySet()) {
            if (!serverResponseMap.containsKey(toServerName)) {
                String warning = "Warning: There is no corresponding response to server \"" + 
                               toServerName + "\" from server \"" + 
                               syncReqServerPair.get(toServerName) + "\"";
                warnings.add(warning);
            }
        }

        results.put("nodeIndices", nodeIndices);
        results.put("syncPairs", syncPairs);
        results.put("warnings", warnings);
        return results;
    }
}

class TaskNode {
    String name;
    String serverName;
    TaskNode parent;
    List<TaskNode> children = new ArrayList<>();
    int index;
    boolean isSync;
    boolean isRoot;

    public TaskNode(String name, String serverName, boolean isSync, boolean isRoot) {
        this.name = name;
        this.serverName = serverName;
        this.isSync = isSync;
        this.isRoot = isRoot;
    }

    public String getServerName() {
        return serverName;
    }

    public void addChild(TaskNode child) {
        child.parent = this;
        children.add(child);
    }
}
