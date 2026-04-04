import org.junit.Test;
import org.junit.Before;
import static org.junit.Assert.*;
import java.util.Collection;
import java.util.ArrayList;

class NameParserTest {

    private NameParser parser;

    @Before
    public void setUp() {
        parser = new NameParser();
    }

    @Test
    public void testParse_EmptyString() {
        Collection<NameToken> result = parser.parse("");
        assertNotNull(result);
        assertEquals(0, result.size());
    }

    @Test
    public void testParse_SingleFixedToken() {
        Collection<NameToken> result = parser.parse("test");
        assertNotNull(result);
        assertEquals(1, result.size());
        NameToken token = result.iterator().next();
        assertTrue(token instanceof FixedToken);
        assertEquals("test", ((FixedToken) token).getValue());
    }

    @Test
    public void testParse_SingleNumericRange() {
        Collection<NameToken> result = parser.parse("[1-10]");
        assertNotNull(result);
        assertEquals(1, result.size());
        NameToken token = result.iterator().next();
        assertTrue(token instanceof NumericRange);
    }

    @Test
    public void testParse_MultipleFixedTokens() {
        Collection<NameToken> result = parser.parse("folder/file/name");
        assertNotNull(result);
        assertEquals(3, result.size());
        
        ArrayList<NameToken> list = new ArrayList<>(result);
        assertTrue(list.get(0) instanceof FixedToken);
        assertEquals("folder", ((FixedToken) list.get(0)).getValue());
        assertTrue(list.get(1) instanceof FixedToken);
        assertEquals("file", ((FixedToken) list.get(1)).getValue());
        assertTrue(list.get(2) instanceof FixedToken);
        assertEquals("name", ((FixedToken) list.get(2)).getValue());
    }

    @Test
    public void testParse_MixedTokens() {
        Collection<NameToken> result = parser.parse("prefix/[1-5]/suffix");
        assertNotNull(result);
        assertEquals(3, result.size());
        
        ArrayList<NameToken> list = new ArrayList<>(result);
        assertTrue(list.get(0) instanceof FixedToken);
        assertEquals("prefix", ((FixedToken) list.get(0)).getValue());
        assertTrue(list.get(1) instanceof NumericRange);
        assertTrue(list.get(2) instanceof FixedToken);
        assertEquals("suffix", ((FixedToken) list.get(2)).getValue());
    }

    @Test
    public void testParse_WithSpaces() {
        Collection<NameToken> result = parser.parse("  part1  /  part2  ");
        assertNotNull(result);
        assertEquals(2, result.size());
        
        ArrayList<NameToken> list = new ArrayList<>(result);
        assertEquals("part1", ((FixedToken) list.get(0)).getValue());
        assertEquals("part2", ((FixedToken) list.get(1)).getValue());
    }

    @Test
    public void testParse_EmptyPartsIgnored() {
        Collection<NameToken> result = parser.parse("a//b/");
        assertNotNull(result);
        assertEquals(2, result.size());
        
        ArrayList<NameToken> list = new ArrayList<>(result);
        assertEquals("a", ((FixedToken) list.get(0)).getValue());
        assertEquals("b", ((FixedToken) list.get(1)).getValue());
    }

    @Test
    public void testParse_OnlyEmptyParts() {
        Collection<NameToken> result = parser.parse("///");
        assertNotNull(result);
        assertEquals(0, result.size());
    }

    @Test
    public void testParse_LeadingTrailingSlashes() {
        Collection<NameToken> result = parser.parse("/start/end/");
        assertNotNull(result);
        assertEquals(2, result.size());
        
        ArrayList<NameToken> list = new ArrayList<>(result);
        assertEquals("start", ((FixedToken) list.get(0)).getValue());
        assertEquals("end", ((FixedToken) list.get(1)).getValue());
    }

    @Test
    public void testParse_MultipleNumericRanges() {
        Collection<NameToken> result = parser.parse("[1-10]/[20-30]");
        assertNotNull(result);
        assertEquals(2, result.size());
        
        ArrayList<NameToken> list = new ArrayList<>(result);
        assertTrue(list.get(0) instanceof NumericRange);
        assertTrue(list.get(1) instanceof NumericRange);
    }

    @Test
    public void testParse_NullInput() {
        try {
            parser.parse(null);
            fail("Expected NullPointerException");
        } catch (NullPointerException e) {
            // Expected exception
        }
    }

    @Test
    public void testParse_SingleCharacterToken() {
        Collection<NameToken> result = parser.parse("a");
        assertNotNull(result);
        assertEquals(1, result.size());
        assertEquals("a", ((FixedToken) result.iterator().next()).getValue());
    }

    @Test
    public void testParse_WhitespaceOnlyParts() {
        Collection<NameToken> result = parser.parse("   /  / ");
        assertNotNull(result);
        assertEquals(0, result.size());
    }

    @Test
    public void testParse_ComplexMixedCase() {
        Collection<NameToken> result = parser.parse("  dir1  /[100-200]/  file  /[1-5]/end  ");
        assertNotNull(result);
        assertEquals(5, result.size());
        
        ArrayList<NameToken> list = new ArrayList<>(result);
        assertEquals("dir1", ((FixedToken) list.get(0)).getValue());
        assertTrue(list.get(1) instanceof NumericRange);
        assertEquals("file", ((FixedToken) list.get(2)).getValue());
        assertTrue(list.get(3) instanceof NumericRange);
        assertEquals("end", ((FixedToken) list.get(4)).getValue());
    }
}