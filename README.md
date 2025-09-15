# Practice_2

## API Endpoints

### Offset Paging
- **Endpoint**: `GET /movies?page=&size=`
- **Description**: Retrieve movies using offset-based pagination.
- **Parameters**:
  - `page`: Page number (default: 1).
  - `size`: Number of items per page (default: 10).

### Cursor Paging
- **Endpoint**: `GET /movies?cursor=`
- **Description**: Retrieve movies using cursor-based pagination.
- **Parameters**:
  - `cursor`: Cursor value for the next page (default: empty for the first page).
  - `size`: Number of items per page (default: 10).

## Performance Comparison
- Use `EXPLAIN` to analyze query plans.
- Log execution time for both offset and cursor paging.
- Compare performance for page 1 and page 1000.