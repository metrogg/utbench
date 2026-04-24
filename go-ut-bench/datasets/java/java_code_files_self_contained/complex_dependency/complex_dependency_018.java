import java.util.HashMap;
import java.util.Map;

class RegionTextureManager {
    private Map<String, TexturePack> regionPacks = new HashMap<>();
    private TexturePack defaultPack;
    private Map<String, Boolean> excludedPlayers = new HashMap<>();

    public RegionTextureManager(TexturePack defaultPack) {
        this.defaultPack = defaultPack;
    }

    /**
     * Handles player movement between regions with different texture packs
     * @param playerId The ID of the player moving
     * @param fromRegion The region the player is coming from (can be null)
     * @param toRegion The region the player is moving to (can be null)
     * @return A status message about the texture pack change
     */
    public String handlePlayerMovement(String playerId, String fromRegion, String toRegion) {
        if (isExcluded(playerId)) {
            return "Player " + playerId + " is excluded from automatic texture changes";
        }

        TexturePack fromPack = fromRegion != null ? regionPacks.get(fromRegion) : null;
        TexturePack toPack = toRegion != null ? regionPacks.get(toRegion) : null;

        if (toPack == fromPack) {
            return "No texture change needed - same pack in both regions";
        }

        if (toPack == null) {
            return exitRegion(playerId, fromPack, toPack);
        }

        if (fromPack == null) {
            return enterRegion(playerId, fromPack, toPack);
        }

        if (toPack.isCustom()) {
            return enterCustomRegion(playerId);
        }

        if (fromPack.isCustom()) {
            return exitCustomRegion(playerId);
        }

        return enterRegion(playerId, fromPack, toPack);
    }

    private boolean isExcluded(String playerId) {
        return excludedPlayers.getOrDefault(playerId, false);
    }

    private String enterRegion(String playerId, TexturePack fromPack, TexturePack toPack) {
        if (toPack.isCustom()) {
            return enterCustomRegion(playerId);
        }
        return "Applied texture pack " + toPack.getName() + " to player " + playerId;
    }

    private String exitRegion(String playerId, TexturePack fromPack, TexturePack toPack) {
        if (fromPack != null && fromPack.isCustom()) {
            return exitCustomRegion(playerId);
        }
        return "Applied default texture pack to player " + playerId;
    }

    private String enterCustomRegion(String playerId) {
        excludedPlayers.put(playerId, true);
        return "NOTE: Player " + playerId + " entered custom region and is excluded from automatic changes";
    }

    private String exitCustomRegion(String playerId) {
        return "Player " + playerId + " left custom region - use command to re-enable auto changes";
    }

    public void addRegionPack(String regionId, TexturePack pack) {
        regionPacks.put(regionId, pack);
    }

    public void setPlayerExcluded(String playerId, boolean excluded) {
        excludedPlayers.put(playerId, excluded);
    }
}

class TexturePack {
    private String name;
    private boolean isCustom;

    public TexturePack(String name, boolean isCustom) {
        this.name = name;
        this.isCustom = isCustom;
    }

    public String getName() {
        return name;
    }

    public boolean isCustom() {
        return isCustom;
    }
}
