import org.junit.Test;
import org.junit.Before;
import static org.junit.Assert.*;

public class NameParserTest {

    private NameParser parser;

    @Before
    public void setUp() {
        parser = new NameParser();
    }

    @Test
    public void testParseSingleFixedToken() {
        Collection<NameToken> result = parser.parse("foo");
        assertEquals(1, result.size());
        assertTrue(result.iterator().next() instanceof FixedToken);
    }

    @Test
    public void testParseSingleNumericRange() {
        Collection<NameToken> result = parser.parse("[1-10]");
        assertEquals(1, result.size());
        assertTrue(result.iterator().next() instanceof NumericRange);
    }

    @Test
    public void testParseMultipleFixedTokens() {
        Collection<NameToken> result = parser.parse("foo/bar");
        assertEquals(2, result.size());
        Object[] tokens = result.toArray();
        assertTrue(tokens[0] instanceof FixedToken);
        assertTrue(tokens[1] instanceof FixedToken);
    }

    @Test
    public void testParseMixedTokens() {
        Collection<NameToken> result = parser.parse("prefix/[1-5]");
        assertEquals(2, result.size());
        Object[] tokens = result.toArray();
        assertTrue(tokens[0] instanceof FixedToken);
        assertTrue(tokens[1] instanceof NumericRange);
    }

    @Test
    public void testParseWithEmptySegments() {
        // Input "a//b" splits into ["a", "", "b"]
        // The empty middle segment should be skipped
        Collection<NameToken> result = parser.parse("a//b");
        assertEquals(2, result.size());
        Object[] tokens = result.toArray();
        assertTrue(tokens[0] instanceof FixedToken);
        assertTrue(tokens[1] instanceof FixedToken);
    }

    @Test
    public void testParseWithLeadingSlash() {
        // Input "/foo" splits into ["", "foo"]
        // The empty leading segment should be skipped
        Collection<NameToken> result = parser.parse("/foo");
        assertEquals(1, result.size());
        assertTrue(result.iterator().next() instanceof FixedToken);
    }

    @Test
    public void testParseWithTrailingSlash() {
        // Input "foo/" splits into ["foo"]
        // Trailing empty strings are discarded by split
        Collection<NameToken> result = parser.parse("foo/");
        assertEquals(1, result.size());
        assertTrue(result.iterator().next() instanceof FixedToken);
    }

    @Test
    public void testParseWithWhitespace() {
        // Input "  foo  " should be trimmed to "foo"
        Collection<NameToken> result = parser.parse("  foo  ");
        assertEquals(1, result.size());
        assertTrue(result.iterator().next() instanceof FixedToken);
    }

    @Test
    public void testParseEmptyString() {
        Collection<NameToken> result = parser.parse("");
        assertEquals(0, result.size());
    }

    @Test
    public void testParseOnlySlashes() {
        // Input "///" splits into []
        Collection<NameToken> result = parser.parse("///");
        assertEquals(0, result.size());
    }

    @Test(expected = NullPointerException.class)
    public void testParseNullInput() {
        parser.parse(null);
    }
}