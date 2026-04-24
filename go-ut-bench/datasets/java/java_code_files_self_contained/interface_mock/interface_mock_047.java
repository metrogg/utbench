// Converted Java method
import java.util.List;
import java.util.Objects;

class EnhancedPositionDao {
    private final PositionMapper positionMapper;
    private final PositionValidator positionValidator;

    public EnhancedPositionDao(PositionMapper positionMapper, PositionValidator positionValidator) {
        this.positionMapper = Objects.requireNonNull(positionMapper, "PositionMapper cannot be null");
        this.positionValidator = Objects.requireNonNull(positionValidator, "PositionValidator cannot be null");
    }

    /**
     * Creates a new position with validation and returns the generated ID.
     * @param position Position to create
     * @return Created position with generated ID
     * @throws IllegalArgumentException if position is invalid
     */
    public PositionPo createPosition(PositionPo position) {
        positionValidator.validateForCreate(position);
        positionMapper.insert(position);
        return position;
    }

    /**
     * Updates an existing position with validation and returns the number of affected rows.
     * @param position Position to update
     * @return Number of affected rows
     * @throws IllegalArgumentException if position is invalid or doesn't exist
     */
    public int updatePosition(PositionPo position) {
        positionValidator.validateForUpdate(position);
        if (positionMapper.selectOne(position) == null) {
            throw new IllegalArgumentException("Position not found");
        }
        return positionMapper.updateByPrimaryKey(position);
    }

    /**
     * Deletes a position and returns the number of affected rows.
     * @param position Position to delete
     * @return Number of affected rows
     * @throws IllegalArgumentException if position doesn't exist
     */
    public int deletePosition(PositionPo position) {
        if (positionMapper.selectOne(position) == null) {
            throw new IllegalArgumentException("Position not found");
        }
        return positionMapper.delete(position);
    }

    /**
     * Finds positions matching criteria with pagination.
     * @param criteria Search criteria
     * @param pageNumber Page number (1-based)
     * @param pageSize Page size
     * @return List of matching positions
     */
    public List<PositionPo> findPositions(PositionPo criteria, int pageNumber, int pageSize) {
        if (pageNumber < 1 || pageSize < 1) {
            throw new IllegalArgumentException("Page number and size must be positive");
        }
        int offset = (pageNumber - 1) * pageSize;
        return positionMapper.selectWithPagination(criteria, offset, pageSize);
    }

    /**
     * Counts positions matching criteria.
     * @param criteria Search criteria
     * @return Count of matching positions
     */
    public int countPositions(PositionPo criteria) {
        return positionMapper.selectCount(criteria);
    }
}

// Supporting interface and class (needed for the implementation)
interface PositionMapper {
    int insert(PositionPo position);
    int updateByPrimaryKey(PositionPo position);
    int delete(PositionPo position);
    PositionPo selectOne(PositionPo position);
    List<PositionPo> selectWithPagination(PositionPo criteria, int offset, int pageSize);
    int selectCount(PositionPo criteria);
}

class PositionValidator {
    public void validateForCreate(PositionPo position) {
        if (position == null) {
            throw new IllegalArgumentException("Position cannot be null");
        }
        // Add more validation logic here
    }

    public void validateForUpdate(PositionPo position) {
        if (position == null || position.getId() == null) {
            throw new IllegalArgumentException("Position and ID cannot be null");
        }
        // Add more validation logic here
    }
}

class PositionPo {
    private Integer id;
    // Other position properties

    public Integer getId() {
        return id;
    }

    public void setId(Integer id) {
        this.id = id;
    }
    // Other getters and setters
}
