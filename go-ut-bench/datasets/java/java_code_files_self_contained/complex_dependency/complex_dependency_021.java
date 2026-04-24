// Converted Java method
import java.util.ArrayList;
import java.util.List;

class TicTacToeGame {
    private char[][] board;
    private char currentPlayer;
    private boolean gameOver;
    private String winner;

    /**
     * Initializes a new Tic-Tac-Toe game with a board of specified size.
     * @param size The size of the board (size x size)
     */
    public TicTacToeGame(int size) {
        if (size < 3) {
            throw new IllegalArgumentException("Board size must be at least 3x3");
        }
        this.board = new char[size][size];
        initializeBoard();
        this.currentPlayer = 'X';
        this.gameOver = false;
        this.winner = null;
    }

    private void initializeBoard() {
        for (int i = 0; i < board.length; i++) {
            for (int j = 0; j < board[i].length; j++) {
                board[i][j] = '-';
            }
        }
    }

    /**
     * Makes a move on the board at the specified coordinates.
     * @param row The row index (0-based)
     * @param col The column index (0-based)
     * @return true if the move was successful, false if the spot was already taken or game is over
     */
    public boolean makeMove(int row, int col) {
        if (gameOver || row < 0 || row >= board.length || col < 0 || col >= board.length) {
            return false;
        }

        if (board[row][col] != '-') {
            return false;
        }

        board[row][col] = currentPlayer;
        
        if (checkWin(row, col)) {
            gameOver = true;
            winner = String.valueOf(currentPlayer);
            return true;
        }

        if (checkDraw()) {
            gameOver = true;
            return true;
        }

        currentPlayer = (currentPlayer == 'X') ? 'O' : 'X';
        return true;
    }

    private boolean checkWin(int row, int col) {
        return checkRow(row) || checkColumn(col) || checkDiagonals(row, col);
    }

    private boolean checkRow(int row) {
        for (int i = 0; i < board.length; i++) {
            if (board[row][i] != currentPlayer) {
                return false;
            }
        }
        return true;
    }

    private boolean checkColumn(int col) {
        for (int i = 0; i < board.length; i++) {
            if (board[i][col] != currentPlayer) {
                return false;
            }
        }
        return true;
    }

    private boolean checkDiagonals(int row, int col) {
        if (row != col && row + col != board.length - 1) {
            return false;
        }

        boolean diagonal1 = true;
        boolean diagonal2 = true;
        
        for (int i = 0; i < board.length; i++) {
            if (board[i][i] != currentPlayer) {
                diagonal1 = false;
            }
            if (board[i][board.length - 1 - i] != currentPlayer) {
                diagonal2 = false;
            }
        }
        
        return diagonal1 || diagonal2;
    }

    private boolean checkDraw() {
        for (int i = 0; i < board.length; i++) {
            for (int j = 0; j < board[i].length; j++) {
                if (board[i][j] == '-') {
                    return false;
                }
            }
        }
        return true;
    }

    /**
     * Gets the current state of the board.
     * @return A list of strings representing each row of the board
     */
    public List<String> getBoardState() {
        List<String> state = new ArrayList<>();
        for (char[] row : board) {
            state.add(new String(row));
        }
        return state;
    }

    /**
     * Gets the current game status.
     * @return A string representing the game status ("X wins", "O wins", "Draw", or "In progress")
     */
    public String getGameStatus() {
        if (winner != null) {
            return winner + " wins";
        }
        if (gameOver) {
            return "Draw";
        }
        return "In progress";
    }

    /**
     * Gets the current player's turn.
     * @return The symbol of the current player ('X' or 'O')
     */
    public char getCurrentPlayer() {
        return currentPlayer;
    }

    /**
     * Checks if the game is over.
     * @return true if the game is over, false otherwise
     */
    public boolean isGameOver() {
        return gameOver;
    }
}
