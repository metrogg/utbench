import org.junit.Test;
import static org.junit.Assert.*;
import java.util.ArrayList;
import java.util.Collection;
import java.util.List;

class NameParserTest {

    private final NameParser parser = new NameParser();

    @Test(expected = NullPointerException.class)
    public void parse_NullInput_ThrowsNullPointerException() {
        parser.parse(null);
    }

    @Test
    public void parse_EmptyString_ReturnsEmptyCollection() {
        Collection<NameToken> result = parser.parse("");
        assertTrue(result.isEmpty());
    }

    @Test
    public void parse_OnlyWhitespaceAndSlashes_ReturnsEmptyCollection() {
        Collection<NameToken> result = parser.parse("  /  // \t\n /  ");
        assertTrue(result.isEmpty());
    }

    @Test
    public void parse_SingleFixedToken_ReturnsCorrectFixedToken() {
        Collection<NameToken> result = parser.parse("sampleName");
        assertEquals(1, result.size());
        assertTrue(result.iterator().next() instanceof FixedToken);
    }

    @Test
    public void parse_SingleNumericRangeToken_ReturnsCorrectNumericRange() {
        Collection<NameToken> result = parser.parse("[1-50]");
        assertEquals(1, result.size());
        assertTrue(result.iterator().next() instanceof NumericRange);
    }

    @Test
    public void parse_MultipleFixedTokens_ReturnsAllFixedTokens() {
        Collection<NameToken> result = parser.parse("part1/part2/part3");
        assertEquals(3, result.size());
        for (NameToken token : result) {
            assertTrue(token instanceof FixedToken);
        }
    }

    @Test
    public void parse_MixedFixedAndRangeTokens_ReturnsCorrectTokenTypes() {
        Collection<NameToken> result = parser.parse("prefix/[10-20]/middle/[30-40]/suffix");
        assertEquals(5, result.size());
        List<NameToken> tokenList = new ArrayList<>(result);
        
        assertTrue(tokenList.get(0) instanceof FixedToken);
        assertTrue(tokenList.get(1) instanceof NumericRange);
        assertTrue(tokenList.get(2) instanceof FixedToken);
        assertTrue(tokenList.get(3) instanceof NumericRange);
        assertTrue(tokenList.get(4) instanceof FixedToken);
    }

    @Test
    public void parse_PartsWithWhitespace_TrimsWhitespaceAndSkipsEmpty() {
        Collection<NameToken> result = parser.parse("  first  /  [5-15]  //  third  ");
        assertEquals(3, result.size());
        List<NameToken> tokenList = new ArrayList<>(result);
        
        assertTrue(tokenList.get(0) instanceof FixedToken);
        assertTrue(tokenList.get(1) instanceof NumericRange);
        assertTrue(tokenList.get(2) instanceof FixedToken);
    }

    @Test
    public void parse_LeadingTrailingAndConsecutiveSlashes_SkipsEmptyParts() {
        Collection<NameToken> result = parser.parse("/a//b/c//");
        assertEquals(3, result.size());
    }

    @Test
    public void parse_SingleBracketToken_ReturnsNumericRange() {
        Collection<NameToken> result = parser.parse("[");
        assertEquals(1, result.size());
        assertTrue(result.iterator().next() instanceof NumericRange);
    }
}